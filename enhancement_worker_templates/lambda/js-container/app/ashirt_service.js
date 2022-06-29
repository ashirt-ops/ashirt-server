const { createHmac, createHash, randomFillSync } = require("crypto");
const http = require("http");

class AShirtService {
  constructor(config) {
    this.apiUrl = config.apiUrl;
    this.apiPort = config.apiPort;
    this.accessKey = config.accessKey;
    this.secretKey = Buffer.from(config.secretKeyB64, "base64");
  }

  async getOperations() {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations`
    })
  }

  async checkConnection() {
    return this.makeRequest({
      method: 'GET',
      path: `/api/checkconnection`
    })
  }

  async createOperation(body) {
    return this.makeRequest({
      method: 'POST',
      path: `/api/operations`,
      body: JSON.stringify(body)
    })
  }

  async getEvidence(operationSlug, evidenceUuid) {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}`
    })
  }

  async getEvidenceContent(operationSlug, evidenceUuid, type = 'media') {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}/${type}`,
    })
  }

  async createEvidence(operationSlug, body) {
    const { file } = body
    const fields = {
      notes: body.notes,
      contentType: body.contentType,
      occurred_at: body.occurred_at,
      tagIds: JSON.stringify(body.tagIds)
    }

    const { boundary: boundary, data } = encodeForm(fields, { file })
    return this.makeRequest({
      method: 'POST',
      path: `/api/operations/${operationSlug}/evidence`,
      body: data,
      multipartFormBoundary: boundary,
    })
  }

  async updateEvidence(operationSlug, evidenceUuid, body) {
    const { file } = body
    const fields = {
      notes: body.notes,
      contentType: body.contentType,
      occurred_at: body.occurred_at,
      tagsToAdd: body.tagsToAdd ? JSON.stringify(body.tagsToAdd) : undefined,
      tagsToRemove: body.tagsToRemove ? JSON.stringify(body.tagsToRemove) : undefined,
    }

    const { boundary: boundary, data } = encodeForm(fields, { file })

    return this.makeRequest({
      method: 'PUT',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}`,
      body: data,
      multipartFormBoundary: boundary,
    })
  }

  async upsertEvidenceMetadata(operationSlug, evidenceUuid, body) {
    return this.makeRequest({
      method: 'PUT',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}/metadata`,
      body: JSON.stringify(body)
    })
  }

  async getOperationTags(operationSlug) {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations/${operationSlug}/tags`
    })
  }

  async createOperationTag(operationSlug, body) {
    return this.makeRequest({
      method: 'POST',
      path: `/api/operations/${operationSlug}/tags`,
      body: JSON.stringify(body),
    })
  }

  async makeRequest(config) {
    const reqConfig = this.buildRequestConfig(config);

    const resp = await this._reqPromise(reqConfig);

    return {
      contentType: resp.headers["content-type"],
      statusCode: resp.statusCode,
      data: resp.data,
    };
  }

  _reqPromise(config) {
    return new Promise((resolve, reject) => {
      const req = http.request(config.httpsOptions, (res) => {
        let data = Buffer.from([]);
        res.on("data", (chunk) => (data = Buffer.concat([data, chunk])));
        res.on("close", () => {
          resolve({
            data,
            statusCode: res.statusCode,
            headers: res.headers,
          });
        });
      });
      req.on("error", reject);
      if (config.data) {
        req.write(config.data);
      }
      req.end();
    });
  }

  buildRequestConfig(config) {
    const sendBody = config.body ?? "";
    const now = AShirtService.nowInRFC1123();
    const auth = this.generateAuthorizationHeaderValue({
      method: config.method,
      body: sendBody,
      date: now,
      path: config.path,
    });

    const rtn = {
      httpsOptions: {
        host: this.apiUrl,
        port: this.apiPort,
        path: config.path,
        method: config.method,
        headers: {
          "Content-Type": config.multipartFormBoundary
            ? `multipart/form-data; boundary=${config.multipartFormBoundary}`
            : "application/json",
          Date: now,
          Authorization: auth,
        },
      },
    };

    if (sendBody != "") {
      rtn.data = sendBody;
    }

    return rtn;
  }

  static nowInRFC1123() {
    return new Date().toUTCString();
  }

  generateAuthorizationHeaderValue(data) {
    const stringBuff = Buffer.from(
      data.method + "\n" + data.path + "\n" + data.date + "\n"
    );
    // note that this isn't encoded -- the result is a series of raw bytes.
    const bodyDigest = createHash("sha256").update(data.body).digest();

    const message = Buffer.concat([stringBuff, bodyDigest]);
    const hmacMessage = createHmac("sha256", this.secretKey)
      .update(message)
      .digest("base64");

    return `${this.accessKey}:${hmacMessage}`;
  }
}

function randomChars(length) {
  const buff = Buffer.alloc(length);
  return randomFillSync(buff).toString("base64url");
}

function encodeForm(fields, files) {
  const boundary = "----AShirtFormData-" + randomChars(30);
  const newline = "\r\n";
  const boundaryStart = "--" + boundary + newline;
  const lastBoundary = "--" + boundary + "--" + newline;

  let fieldBuffer = Buffer.from("");
  Object.entries(fields).forEach(([key, value]) => {
    if (value === undefined) {
      return;
    }
    const text =
      boundaryStart +
      `Content-Disposition: form-data; name="${key}"` +
      newline +
      newline +
      value +
      newline;
    fieldBuffer = Buffer.concat([fieldBuffer, Buffer.from(text)]);
  });

  let fileBuffer = Buffer.from("");
  Object.entries(files).forEach(([key, fd]) => {
    if (fd === undefined) {
      return;
    }
    const textPart =
      `${boundaryStart}` +
      `Content-Disposition: form-data; name="${key}"; filename="${fd.filename}"` +
      `${newline}Content-Type: ${fd.mimetype}` +
      `${newline}${newline}`;

    fileBuffer = Buffer.concat([
      fileBuffer,
      Buffer.from(textPart),
      fd.content,
      Buffer.from(newline),
    ]);
  });

  return {
    boundary: boundary,
    data: Buffer.concat([fieldBuffer, fileBuffer, Buffer.from(lastBoundary)]),
  };
}


module.exports = {
  AShirtService: new AShirtService({
    apiUrl: process.env.ASHIRT_BACKEND_URL ?? "",
    apiPort: parseInt(process.env.ASHIRT_BACKEND_PORT, 10),
    accessKey: process.env.ASHIRT_ACCESS_KEY ?? "",
    secretKeyB64: process.env.ASHIRT_SECRET_KEY ?? "",
  }),
};
