import { GlobalVar } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createGlobalVar(name: string, value: string | null): Promise<GlobalVar> {
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
  await ds.deleteGlobalVar({ name })
}

export async function updateGlobalVar(name: string, i: { value: string | null, newName: string | null }): Promise<void> {
  await ds.updateGlobalVar({ name, }, i)
}
