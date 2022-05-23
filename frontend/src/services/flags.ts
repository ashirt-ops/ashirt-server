// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { backendDataSource as ds } from './data_sources/backend'

let backendFlags: Array<string> | null = null

export async function getFlags(force?: true): Promise<Array<string>> {
  if (force || backendFlags == null) {
    const { flags } = await ds.flags()
    backendFlags = flags
  }
  return backendFlags
}

export async function hasFlag(flagName: string): Promise<boolean> {
  return (await getFlags()).includes(flagName)
}
