export type BasicServiceWorkerConfig = {
  type: string
  version: number //int only
  synchronous?: true
}

export type ServiceWorkerWeb = {
  type: 'web'
  version: 1
  url: string
  headers?: Record<string, string> // not allowing multiple values for a header for now
  synchronous?: true
}

export type ServiceWorkerConfig =
  | ServiceWorkerWeb

type JSONPrimitive = string | boolean | number | null

export const parseServiceConfig = (config: string): ServiceWorkerConfig | null => {
  // restrict JSON.parse type into something reasonable. You might be able to trick this, but...
  const obj = JSON.parse(config) as Record<string, unknown> | Array<unknown> | JSONPrimitive

  if (typeof obj === 'object' && obj !== null && !Array.isArray(obj) && isWebConfig(obj)) {
    return obj
  }
  return null
}

const isWebConfig = (json: Record<string, unknown>): json is ServiceWorkerWeb => {
  const modelConfig: ServiceWorkerWeb = {
    type: 'web',
    url: '',
    version: 1,
  }

  return (
    hasValue(json, "type", modelConfig.type)
    && hasValue(json, "version", modelConfig.version)
    && hasValueType(json, "url", typeof modelConfig.url)
    && hasOptionalRecord(json, "headers")
  )
}

const hasValue = (
  object: Record<string, unknown>,
  field: string,
  value: unknown
) => {
  return object[field] == value
}

const hasValueType = (
  object: Record<string, unknown>,
  field: string,
  valueType: string
) => {
  return typeof (object[field]) == valueType
}

const hasOptionalRecord = (
  object: Record<string, unknown>,
  field: string,
) => {
  if (field in object) {
    return isRecord(object[field])
  }
  return true
}

const isRecord = (o: unknown): o is Record<string, unknown> => {
  return typeof o === "object" && o !== null && !Array.isArray(o)
}
