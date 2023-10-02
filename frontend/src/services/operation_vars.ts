// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { OperationVar } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createOperationVar(name: string, value: string | null): Promise<OperationVar> {
  let slug = name.toLowerCase().replace(/[^A-Za-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
  if (slug === "") {
    return (name === ""
      ? Promise.reject(Error("Operation Name must not be empty"))
      : Promise.reject(Error("Operation Name must include letters or numbers"))
    )
  }
  try {
    return await ds.createOperationVar({ slug, name, value })
  } catch (err) {
    if (err.message.match(/slug already exists/g)) {
      slug += '-' + Date.now()
      return await ds.createOperationVar({ slug, name, value })
    }
    throw err
  }
}

export async function getOperationVars(): Promise<Array<OperationVar>> {
  return await ds.listOperationVars()
}

export async function deleteOperationVar(name: string): Promise<void> {
  await ds.deleteOperationVar({ name })
}

export async function updateOperationVar(name: string, i: { value: string | null, newName: string | null }): Promise<void> {
  await ds.updateOperationVar({ name, }, i)
}
