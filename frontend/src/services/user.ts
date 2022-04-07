// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { NewUser } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function updateUserProfile(i: {
  userSlug: string,
  firstName: string,
  lastName: string,
  email: string,
}): Promise<void> {
  await ds.updateUser({ userSlug: i.userSlug }, {
    firstName: i.firstName,
    lastName: i.lastName,
    email: i.email,
  })
}

export async function deleteUserAuthenticationScheme(i: {
  userSlug: string,
  authSchemeName: string,
}): Promise<void> {
  await ds.deleteUserAuthScheme(i)
}

export async function addHeadlessUser(i: {
  firstName: string,
  lastName: string,
  email: string,
}): Promise<NewUser> {
  return ds.adminCreateHeadlessUser(i)
}

export async function createRecoveryCode(i: {
  userSlug: string
}): Promise<string> {
  const resp = await ds.createRecoveryCode(i)
  return resp.code
}

export async function getTotpForUser(i:{
  userSlug: string
}): Promise<boolean> {
  return await ds.getTotpForUser(i)
}

export async function deleteTotpForUser(i: {
  userSlug: string
}): Promise<void> {
  await ds.deleteTotpForUser(i)
}
