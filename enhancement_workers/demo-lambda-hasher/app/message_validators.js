const MESSAGE_TYPE_TEST = "test"
const MESSAGE_TYPE_EVIDENCE_CREATED = "evidence_created";

function isTestRequest(data) {
  return data.type === MESSAGE_TYPE_TEST;
}

function isEvidenceCreatedEvent(data) {
  return (
    data.type === MESSAGE_TYPE_EVIDENCE_CREATED &&
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

function isValidRequest(data) {
  const validInputFns = [isTestRequest, isEvidenceCreatedEvent];

  return validInputFns.findIndex((fn) => fn(data)) !== -1;
}

module.exports = {
  isTestRequest,
  isEvidenceCreatedEvent,
  isValidRequest,
  MESSAGE_TYPE_TEST,
  MESSAGE_TYPE_EVIDENCE_CREATED,
};
