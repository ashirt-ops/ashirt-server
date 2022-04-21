// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

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
  const obj = JSON.parse(config) as Record<string, unknown> | Array<unknown> | JSONPrimitive

  if (typeof obj === 'object' && !Array.isArray(obj) && obj !== null) {
    if (isWebConfig(obj)) {
      return obj
    }
  }
  return null
}

const isWebConfig = (json: Record<string, unknown>): json is ServiceWorkerWeb => {
  const modelConfig: ServiceWorkerWeb = {
    type: 'web',
    url: '',
    version: 1,
  }

  if (
    hasValue(json, "type", modelConfig.type)
    && hasValue(json, "version", modelConfig.version)
    && hasValueType(json, "url", typeof modelConfig.url)
  ) {
    return (
      hasSomeValue(json, 'synchronous', [true, undefined])
      && hasOptionalRecord(json, "headers")
    )
  }

  return false
}

const hasValue = (
  object: Record<string, unknown>,
  field: string,
  value: unknown
) => {
  return object[field] == value
}

const hasSomeValue = (
  object: Record<string, unknown>,
  field: string,
  valueSet: Array<unknown>
) => {
  return valueSet.includes(object[field])
}

const hasValueType = (
  object: Record<string, unknown>,
  field: string,
  valueType: string
) => {
  return typeof (object[field]) == valueType
}

const hasRecord = (
  object: Record<string, unknown>,
  field: string,
  valueType: string
) => {
  return isRecord(object[field])
}

const hasOptionalArray = (
  object: Record<string, unknown>,
  field: string,
) => {
  if (field in object) {
    return Array.isArray(object[field])
  }
  return true
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
