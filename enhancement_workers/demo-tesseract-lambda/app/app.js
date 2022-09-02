const {
  badRequest,
  errorProcessing,
  notImplemented,
  rejectEvidence,
  testPassed,
  processSuccess,
} = require("./responses");
const { AShirtService } = require("./ashirt_service");
const tesseract = require("node-tesseract-ocr");

exports.handler = handleLambdaEvent;

const tesseractConfig = {
  lang: "eng",
  oem: 1,
  psm: 12,
};

async function handleLambdaEvent(event) {
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
}

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
  if (requestData.contentType != "image") {
    return rejectEvidence();
  }

  try {
    const resp = await AShirtService.getEvidenceContent(
      requestData.operationSlug,
      requestData.evidenceUuid,
      "media"
    );

    const result = await tesseract.recognize(resp.data, tesseractConfig);
    return processSuccess(result);
  } catch (err) {
    return errorProcessing(err);
  }
}
