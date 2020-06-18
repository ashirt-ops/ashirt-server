// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { ApiKey } from 'src/global_types'
import req from './request_helper'

function formatApiKey(apiKey: any): ApiKey {
  return {
    ...apiKey,
    lastAuth: apiKey.lastAuth ? new Date(apiKey.lastAuth) : null,
  }
}

export async function createApiKey(i: { userSlug: string }): Promise<ApiKey> {
  return formatApiKey(await req('POST', `/user/${i.userSlug}/apikeys`))
}

export async function getApiKeys(i?: { userSlug: string }): Promise<Array<ApiKey>> {
  return (await req('GET', '/user/apikeys', null, i)).map(formatApiKey)
}

export async function deleteApiKey(i: {
  userSlug: string,
  accessKey: string,
}): Promise<void> {
  return await req('DELETE', `/user/${i.userSlug}/apikeys/${i.accessKey}`)
}
