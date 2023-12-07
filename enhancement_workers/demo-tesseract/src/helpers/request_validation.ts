
export type SupportedMessage =
  | TestMessage
  | EvidenceCreatedMessage

export type TestMessage = {
  type: "test"
}

export type EvidenceCreatedMessage = {
  type: "evidence_created",
  evidenceUuid: string,
  operationSlug: string,
  contentType: typeof SupportedContentTypes[number]
  globalVariables: Record<string, unknown>[]
  operationVariables: Record<string, unknown>[]
}

/**
 * isValidProcessRequest is a type guard that reviews the given body to make sure it matches
 * a supported request type.
 * 
 * @param body Any json-parsed content. This is not guaranteed to be correct for random bodies --
 * especially esoteric javascript objects.
 * @returns true if the value given has a supported request type. Otherwise, this returns false
 */
export const isSupportedMessage = (body: unknown): body is SupportedMessage => {
  if (
    isRecordWithTypeField(body)
    && isString(body.type)
    && SupportedRequestTypes.includes((body.type as typeof SupportedRequestTypes[number]))
  ) {
    const basicRequest = body as BasicRequest
    return (
      isTestRequest(basicRequest)
      || isProcessRequest(basicRequest)
    )
  }
  return false
}

/**
 * isJsonWithType is a type guard that verifies that the given value is an object with a "type"
 * field.
 * @param body Anything.
 * @returns true if the value given is an object with a type field. Otherwise, this returns false.
 */
const isRecordWithTypeField = (body: unknown): body is { type: unknown } => {
  return (
    typeof body === 'object' &&
    body !== null &&
    "type" in body
  )
}

const isTestRequest = (body: BasicRequest): body is TestMessage => {
  return body.type == 'test'
}

/**
 * isProcessRequest is a typeguard that verifies that the given body matches the expected type for
 * a process-type request. This does not verify that the content "makes sense", bui
 * @param body 
 * @returns 
 */
const isProcessRequest = (body: BasicRequest): body is EvidenceCreatedMessage => {
  return (
    body.type === 'evidence_created' &&
    hasField(body, "evidenceUuid", isString) &&
    hasField(body, "operationSlug", isString) &&
    hasField(body, "contentType", (v) => (
      isString(v) &&
      // typescript can't compare the string literal of SupportedContentTypes with a random string
      // so we tell it that the random string _is_ one of the types, but `includes` will actually
      // verify that the random string is one of the supported literals.
      SupportedContentTypes.includes(v as typeof SupportedContentTypes[number]))
    )
  )
}

const isString = (v: unknown): v is string => (
  typeof v === 'string'
)

const hasField = (o: Record<string, unknown>, field: string, ofType: (b: unknown) => boolean) => {
  return field in o && ofType(o[field])
}

export const SupportedContentTypes = [
  "http-request-cycle",
  "terminal-recording",
  "codeblock",
  "event",
  "image",
  "none",
] as const

export const SupportedRequestTypes = [
  "evidence_created",
  "test",
] as const

type BasicRequest = Record<string, unknown> & {
  type: typeof SupportedRequestTypes[number]
}

