import { default as axios, AxiosRequestConfig } from 'axios'
import { createHmac, createHash } from 'crypto'
import { formatRFC7231 } from 'date-fns'

export type RequestConfig = {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE'
  path: string
  body?: string
  responseType?: 'arraybuffer' | 'document' | 'json' | 'text' | 'stream'
}

export class AShirtService {
  private secretKey: Buffer
  constructor(
    private apiUrl: string,
    private accessKey: string,
    secretKeyB64: string
  ) {
    this.secretKey = Buffer.from(secretKeyB64, "base64")
  }

  async getEvidence(operationSlug: string, evidenceUuid: string) {
    return this.makeRequest<EvidenceOutput>({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}`
    })
  }

  async getEvidenceContent(operationSlug: string, evidenceUuid: string, type: 'media' | 'preview' = 'media'): Promise<ResponseWrapper<Buffer>> {
    return this.makeRequest<Buffer>({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}/${type}`,
      responseType: 'arraybuffer'
    })
  }

  async makeRequest<T>(
    config: RequestConfig,
    guard?: (o: unknown) => o is T,
  ): Promise<ResponseWrapper<T>> {
    const reqConfig = this.buildRequestConfig(config)
    const resp = await axios(reqConfig)
    const respData = resp.data

    if (guard && !guard(respData)) {
        Promise.reject(new Error("Response is not in the right format."))
    }

    const reqResult:ResponseWrapper<T> = {
      contentType: resp.headers['content-type'],
      responseCode: resp.status,
      data: respData
    }

    return reqResult
  }

  private buildRequestConfig(config: RequestConfig) {
    const sendBody = config.body ?? ''
    const now = AShirtService.nowInRFC1123()
    const auth = this.generateAuthorizationHeaderValue({
      method: config.method,
      body: sendBody,
      date: now,
      path: config.path,
    })

    const req: AxiosRequestConfig = {
      method: config.method,
      url: `${this.apiUrl}${config.path}`,
      headers: {
        "Content-Type": "application/json",
        "Date": now,
        "Authorization": auth,
      }
    }
    if (sendBody != '') {
      req.data = sendBody
    }
    if( config.responseType) {
      req.responseType = config.responseType
    }
    return req
  }

  private static nowInRFC1123() {
    return formatRFC7231(new Date())
  }

  private generateAuthorizationHeaderValue(data: {
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' // more methods with a similar naming style are possible
    path: string
    date: string // in RFC1123 format
    body: string
  }) {
    const stringBuff = Buffer.from(
      data.method + "\n" +
      data.path + "\n" +
      data.date + "\n"
    )
    // note that this isn't encoded -- the result is a series of raw bytes.
    const bodyDigest = createHash('sha256').update(data.body).digest()

    const message = Buffer.concat([stringBuff, bodyDigest])
    const hmacMessage = createHmac('sha256', this.secretKey)
      .update(message)
      .digest('base64')

    return `${this.accessKey}:${hmacMessage}`
  }
}

type EvidenceOutput = {
  uuid: string
  description: string
  contentType: string
  occurredAt: Date
}

type ResponseWrapper<T> = {
  responseCode: number,
  contentType: string
  data: T
}
