const { createHmac, createHash } = require("crypto");
const http = require('http');

class AShirtService {
  constructor(config) {
    this.apiUrl = config.apiUrl;
    this.apiPort = config.apiPort;
    this.accessKey = config.accessKey;
    this.secretKey = Buffer.from(config.secretKeyB64, "base64");
  }

  async getEvidence(operationSlug, evidenceUuid) {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}`
    });
  }

  async getEvidenceContent(operationSlug, evidenceUuid, type) {
    return this.makeRequest({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}/${type}`,
      responseType: 'arraybuffer'
    });
  }

  async createOperation(body) {
    return this.makeRequest({
      method: 'POST',
      path: `/api/operations`,
      body: JSON.stringify(body)
    });
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
        let data = Buffer.from([])
        res.on("data", (chunk) =>
          data = Buffer.concat([data, chunk])
        );
        res.on('close', ()=>{
          resolve({
            data,
            statusCode: res.statusCode,
            headers: res.headers,
          })
        })
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
          "Content-Type": "application/json",
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

module.exports = {
  AShirtService: new AShirtService({
    apiUrl: process.env.ASHIRT_BACKEND_URL ?? "",
    apiPort: parseInt(process.env.ASHIRT_BACKEND_PORT, 10),
    accessKey: process.env.ASHIRT_ACCESS_KEY ?? "",
    secretKeyB64: process.env.ASHIRT_SECRET_KEY ?? "",
  }),
};
