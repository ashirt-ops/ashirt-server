// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { DataSource } from './data_sources/data_source'

export async function updateUserProfile(ds: DataSource, i: {
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

export async function deleteUserAuthenticationScheme(ds: DataSource, i: {
  userSlug: string,
  authSchemeName: string,
}): Promise<void> {
  await ds.deleteUserAuthScheme(i)
}

export async function addHeadlessUser(ds: DataSource, i: {
  firstName: string,
  lastName: string,
  email: string,
}): Promise<void> {
  await ds.adminCreateHeadlessUser(i)
}

export async function createRecoveryCode(ds: DataSource, i: {
  userSlug: string
}): Promise<string> {
  const resp = await ds.createRecoveryCode(i)
  return resp.code
}
