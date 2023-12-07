import { ApiKey } from 'src/global_types'
import { apiKeyFromDto } from './data_sources/converters'
import { backendDataSource as ds } from './data_sources/backend'

export async function createApiKey(i: { userSlug: string }): Promise<ApiKey> {
  return apiKeyFromDto(await ds.createApiKey(i))
}

export async function getApiKeys(i?: { userSlug: string }): Promise<Array<ApiKey>> {
  const keys = await ds.listApiKeys(i)
  return keys.map(apiKeyFromDto)
}

export async function deleteApiKey(i: {
  userSlug: string,
  accessKey: string,
}): Promise<void> {
  return await ds.deleteApiKey(i)
}
