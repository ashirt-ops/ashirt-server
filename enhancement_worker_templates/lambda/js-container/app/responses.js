module.exports = {
  testPassed: () => ({
    statusCode: 200,
    body: JSON.stringify({
      status: "ok",
    }),
  }),

  errorProcessing: (message) => ({
    statusCode: 500,
    body: message
      ? JSON.stringify({
          action: "error",
          content: message,
        })
      : undefined,
  }),
  rejectEvidence: () => ({
    statusCode: 406,
  }),
  badRequest: (message) => ({
    statusCode: 400,
    body: JSON.stringify({
      action: "error",
      content: message,
    }),
  }),
  notImplemented: () => ({
    statusCode: 501,
  }),
  processSuccess: (content) => ({
    statusCode: 200,
    body: JSON.stringify({
      action: "processed",
      content,
    }),
  }),
};
