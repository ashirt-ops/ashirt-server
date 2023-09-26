// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { GlobalVar } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createGlobalVar(name: string, value: string): Promise<GlobalVar> {
  if (name === "") {
    return Promise.reject(Error("Global variable name must not be empty"))
  }
  try {
    return await ds.createGlobalVar({ name, value })
  } catch (err) {
    throw err
  }
}

export async function getGlobalVars(): Promise<Array<GlobalVar>> {
  return await ds.listGlobalVars()
}

export async function deleteGlobalVar(name: string): Promise<void> {
  await ds.deleteGlobalVar({ globalVarName: name })
}

export async function updateGlobalVar(name: string, i: { value: string | null, newName: string | null }): Promise<void> {
  await ds.updateGlobalVar({ globalVarName: name, }, i)
}
