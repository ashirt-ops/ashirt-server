// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { OperationVar } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createOperationVar(operationSlug: string, name: string, value: string | null): Promise<OperationVar> {
  let slug = name.toLowerCase().replace(/[^A-Za-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
  if (slug === "") {
    return (name === ""
      ? Promise.reject(Error("Operation Name must not be empty"))
      : Promise.reject(Error("Operation Name must include letters or numbers"))
    )
  }
  try {
    return await ds.createOperationVar({ operationSlug }, {name, value, varSlug: slug })
  } catch (err) {
    if (err.message.match(/slug already exists/g)) {
      slug += '-' + Date.now()
      return await ds.createOperationVar({ operationSlug }, {varSlug: slug, name, value })
    }
    throw err
  }
}

export async function getOperationVars(operationSlug: string): Promise<Array<OperationVar>> {
  return await ds.listOperationVars({operationSlug})
}

export async function deleteOperationVar(operationSlug: string, varSlug: string): Promise<void> {
  await ds.deleteOperationVar({ operationSlug, varSlug })
}

export async function updateOperationVar(operationSlug: string, varSlug: string, i: { value: string | null, name: string | null }): Promise<void> {
  await ds.updateOperationVar({ operationSlug, varSlug}, i)
}
