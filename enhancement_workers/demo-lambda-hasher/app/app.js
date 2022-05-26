const {
  badRequest,
  errorProcessing,
  notImplemented,
  testPassed,
  processSuccess,
} = require("./responses");

const {createHash} = require("crypto")

const { AShirtService } = require("./ashirt_service");

exports.handler = async (event) => {
  
  // Parse and validate the data
  if (!isValidateInput(event)) {
    return badRequest("Body is not in the expected format");
  }

  // handle the request type
  if (event.type === "test") {
    return testPassed();
  }
  if (event.type === "process") {
    return await handleProcess(event);
  }

  // we should never actually get here -- validation above should trap anything that won't work
  return notImplemented();
};

function isValidateInput(data) {
  if (data.type === "test") {
    return true;
  }
  return (
    data.type === "process" &&
    typeof data.evidenceUuid === "string" &&
    typeof data.operationSlug === "string" &&
    [
      // Note: this covers all of current types as of June 2022
      "http-request-cycle",
      "terminal-recording",
      "codeblock",
      "event",
      "image",
      "none",
    ].includes(data.contentType)
  );
}

async function handleProcess(requestData) {
  // accept all forms of content

  try {
    const content = await AShirtService.getEvidenceContent(
      requestData.operationSlug,
      requestData.evidenceUuid,
      "media"
    );
    
    if (content.statusCode == 200) {
      // hash the file contents
      console.log("???", content.data)
      const buffer = Buffer.from(content.data);
      const hashResult = createHash("sha256")
        .update(buffer)
        .digest()
        .toString("hex");

      return processSuccess(hashResult);
    }
    return errorProcessing("Unable to retrieve content")
  } catch (err) {
    return errorProcessing(err);
  }
}
