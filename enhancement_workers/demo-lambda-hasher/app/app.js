const {
  badRequest,
  errorProcessing,
  notImplemented,
  testPassed,
  processSuccess,
} = require("./responses");

const {
  isValidRequest,
  MESSAGE_TYPE_TEST,
  MESSAGE_TYPE_EVIDENCE_CREATED,
} = require("./message_validators");

const { createHash } = require("crypto");

const { AShirtService } = require("./ashirt_service");

exports.handler = async (event) => {
  // Parse and validate the data
  if (!isValidRequest(event)) {
    return badRequest("Body is not in the expected format");
  }

  // handle the request type
  switch (event.type) {
    case MESSAGE_TYPE_TEST:
      return testPassed();
    case MESSAGE_TYPE_EVIDENCE_CREATED:
      return await handleProcess(event);
  }

  // we should never actually get here -- validation above should trap anything that won't work
  return notImplemented();
};

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
      const buffer = Buffer.from(content.data);
      const hashResult = createHash("sha256")
        .update(buffer)
        .digest()
        .toString("hex");

      return processSuccess(hashResult);
    }
    return errorProcessing("Unable to retrieve content");
  } catch (err) {
    return errorProcessing(err);
  }
}
