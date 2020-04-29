// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { UserOwnView, SupportedAuthenticationScheme, AuthSchemeDetails, RecoveryMetrics } from 'src/global_types'
import { DataSource } from './data_sources/data_source'
import { userOwnViewFromDto } from './data_sources/converters'

export async function getCurrentUser(ds: DataSource): Promise<UserOwnView | null> {
  try {
    return userOwnViewFromDto(await ds.readCurrentUser())
  } catch (err) {
    // Not found indicates a non logged in user
    if (err.status === 404) return null

    // Bubble any other error types up
    throw err
  }
}

// TODO this should be encapsulated in an admin settings component under src/authschemes/local
export async function adminChangePassword(ds: DataSource, i: {
  userSlug: string,
  newPassword: string,
}) {
  if (i.newPassword.length < 3) {
    throw Error("User password must be at least 3 characters long")
  }
  await ds.adminChangePassword(i)
}

export async function logout(ds: DataSource) {
  await ds.logout()
}

export async function getUser(ds: DataSource, i?: { userSlug: string }): Promise<UserOwnView> {
  const user = await (i ? ds.readUser(i) : ds.readCurrentUser())
  return userOwnViewFromDto(user)
}

export async function adminSetUserFlags(ds: DataSource, i: {
  userSlug: string,
  disabled: boolean,
  admin: boolean,
}) {
  await ds.adminSetUserFlags(
    { userSlug: i.userSlug },
    { disabled: i.disabled, admin: i.admin },
  )
}

export async function getSupportedAuthentications(ds: DataSource): Promise<Array<SupportedAuthenticationScheme>> {
  return await ds.listSupportedAuths()
}

export async function getSupportedAuthenticationDetails(ds: DataSource): Promise<Array<AuthSchemeDetails>> {
  const schemes = await ds.listAuthDetails()
  return schemes.map(details => ({
    schemeName: details.schemeName,
    schemeCode: details.schemeCode,
    userCount: details.userCount,
    uniqueUserCount: details.uniqueUserCount,
    labels: details.labels,
    lastUsed: details.lastUsed == null ? null : new Date(details.lastUsed), // incoming value is actually a string that we are focing into a date
  }))
}

export async function adminDeleteUser(ds: DataSource, i: {
  userSlug: string,
}) {
  await ds.adminDeleteUser(i)
}

export async function deleteGlobalAuthScheme(ds: DataSource, i: {
  schemeName: string,
}) {
  await ds.deleteGlobalAuthScheme(i)
}

export async function deleteExpiredRecoveryCodes(ds: DataSource) {
  await ds.deleteExpiredRecoveryCodes()
}

export async function getRecoveryMetrics(ds: DataSource): Promise<RecoveryMetrics> {
  return await ds.getRecoveryMetrics()
}
