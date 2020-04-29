// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { ApiKey } from 'src/global_types'
import { apiKeyFromDto } from './data_sources/converters'
import { DataSource } from './data_sources/data_source'

export async function createApiKey(ds: DataSource, i: { userSlug: string }): Promise<ApiKey> {
  return apiKeyFromDto(await ds.createApiKey(i))
}

export async function getApiKeys(ds: DataSource, i?: { userSlug: string }): Promise<Array<ApiKey>> {
  const keys = await ds.listApiKeys()
  return keys.map(apiKeyFromDto)
}

export async function deleteApiKey(ds: DataSource, i: {
  userSlug: string,
  accessKey: string,
}): Promise<void> {
  return await ds.deleteApiKey(i)
}
