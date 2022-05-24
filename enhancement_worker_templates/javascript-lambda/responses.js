module.exports = {
  errorProcessing: (message) => ({
    statusCode: 500,
    body: message
      ? JSON.stringify({
          message,
        })
      : undefined,
  }),
  rejectEvidence: () => ({
    statusCode: 406,
  }),
  badRequest: (message) => ({
    statusCode: 400,
    body: JSON.stringify({ message }),
  }),
  notImplemented: () => ({
    statusCode: 501,
  }),
  testPassed: () => ({
    statusCode: 200,
    body: "ok",
  }),
  processSuccess: (body) => ({
    statusCode: 200,
    body: JSON.stringify(body),
  })
};
