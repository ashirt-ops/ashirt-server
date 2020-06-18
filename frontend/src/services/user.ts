// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'

export async function updateUserProfile(i: {
  userSlug: string,
  firstName: string,
  lastName: string,
  email: string,
}): Promise<void> {
  await req('POST', `/user/profile/${i.userSlug}`, i)
}

export async function deleteUserAuthenticationScheme(i: {
  userSlug: string,
  authSchemeName: string,
}): Promise<void> {
  await req('DELETE', `/user/${i.userSlug}/scheme/${i.authSchemeName}`)
}

export async function addHeadlessUser(i: {
  firstName: string,
  lastName: string,
  email: string,
}): Promise<void> {
  await req('POST', "/admin/user/headless", i)
}

export async function createRecoveryCode(i: {
  userSlug: string
}): Promise<string> {
  const resp = await req('POST', '/auth/recovery/generate', i)
  return resp.code
}
