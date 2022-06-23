import { default as axios, AxiosRequestConfig } from 'axios'
import { createHmac, createHash } from 'crypto'
import {
  CheckConnectionOutput,
  CreateEvidenceInput,
  CreateOperationInput,
  CreateTagInput,
  EvidenceOutput,
  ListOperationsOutput,
  ListOperationTagsOutput,
  OperationOutputItem,
  ReadEvidenceOutput,
  ResponseWrapper,
  UpdateEvidenceInput,
  UpsertMetadataInput,
} from './types'

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

  async getOperations() {
    return this.makeRequest<ListOperationsOutput>({
      method: 'GET',
      path: `/api/operations`
    })
  }

  async checkConnection() {
    return this.makeRequest<CheckConnectionOutput>({
      method: 'GET',
      path: `/api/checkconnection`
    })
  }

  async createOperation(body: CreateOperationInput) {
    return this.makeRequest<OperationOutputItem>({
      method: 'POST',
      path: `/api/operations`,
      body: JSON.stringify(body)
    })
  }

  async getEvidence(operationSlug: string, evidenceUuid: string) {
    return this.makeRequest<ReadEvidenceOutput>({
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

  // TODO
  async createEvidence(operationSlug: string, body: CreateEvidenceInput) {
    return this.makeRequest<EvidenceOutput>({
      method: 'POST',
      path: `/api/operations/${operationSlug}/evidence`,
      body: "" // TODO
    })
  }

  // TODO
  async updateEvidence(operationSlug: string, evidenceUuid: string, body: UpdateEvidenceInput) {
    return this.makeRequest<void>({
      method: 'PUT',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}`,
      body: "" // TODO
    })
  }

  async upsertEvidenceMetadata(operationSlug: string, evidenceUuid: string, body: UpsertMetadataInput) {
    return this.makeRequest<void>({
      method: 'PUT',
      path: `/api/operations/${operationSlug}/evidence/${evidenceUuid}/metadata`,
      body: JSON.stringify(body)
    })
  }

  async getOperationTags(operationSlug: string) {
    return this.makeRequest<ListOperationTagsOutput>({
      method: 'GET',
      path: `/api/operations/${operationSlug}/evidence/tags`
    })
  }

  async createOperationTag(operationSlug: string, body: CreateTagInput) {
    return this.makeRequest<ReadEvidenceOutput>({
      method: 'POST',
      path: `/api/operations/${operationSlug}/evidence/tags`,
      body: JSON.stringify(body),
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

    const reqResult: ResponseWrapper<T> = {
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
    if (config.responseType) {
      req.responseType = config.responseType
    }
    return req
  }

  private static nowInRFC1123() {
    return new Date().toUTCString()
  }

  generateAuthorizationHeaderValue(data: {
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
