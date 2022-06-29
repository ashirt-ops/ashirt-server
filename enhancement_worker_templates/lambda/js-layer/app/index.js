const {
  badRequest,
  errorProcessing,
  notImplemented,
  rejectEvidence,
  testPassed,
  processSuccess,
} = require("./responses");

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
  if (event.type === "evidence_created") {
    return await handleEvidenceCreated(event);
  }

  // we should never actually get here -- validation above should trap anything that won't work
  return notImplemented();
};

function isValidateInput(data) {
  if (data.type === "test") {
    return true;
  }
  return (
    data.type === "evidence_created" &&
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

async function handleEvidenceCreated(requestData) {
  // TODO handle your custom logic here!
  // filter out unprocessable evidence
  if (requestData.contentType != "image") {
    return rejectEvidence();
  }

  try {
    // fetch evidence content
    // for example...
    // const content = AShirtService.getEvidenceContent(
    //   requestData.operationSlug,
    //   requestData.evidenceUuid,
    //   "media"
    // );
    doProcessing();
    return processSuccess("It all works!"); // TODO replace your content here!
  } catch (err) {
    return errorProcessing(err);
  }
}

function doProcessing() {
  // TODO
}
