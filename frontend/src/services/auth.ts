// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { UserOwnView, SupportedAuthenticationScheme, AuthSchemeDetails, RecoveryMetrics } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'
import { userOwnViewFromDto } from './data_sources/converters'

export async function getCurrentUser(): Promise<UserOwnView | null> {
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
export async function adminChangePassword(i: {
  userSlug: string,
  newPassword: string,
}) {
  if (i.newPassword.length < 3) {
    throw Error("User password must be at least 3 characters long")
  }
  await ds.adminChangePassword(i)
}

// TODO this should be encapsulated in an admin settings component under src/authschemes/local
export async function adminCreateLocalUser(i: {
  firstName: string,
  lastName?: string,
  email: string,
}) {
  return await ds.adminCreateLocalUser(i)
}

export async function adminInviteUser(i: {
  firstName: string,
  lastName?: string,
  email: string,
}) {
  return await ds.adminInviteUser(i)
}

export async function logout() {
  await ds.logout()
}

export async function getUser(i?: { userSlug: string }): Promise<UserOwnView> {
  const user = await (i ? ds.readUser(i) : ds.readCurrentUser())
  return userOwnViewFromDto(user)
}

export async function adminSetUserFlags(i: {
  userSlug: string,
  disabled: boolean,
  admin: boolean,
}) {
  await ds.adminSetUserFlags(
    { userSlug: i.userSlug },
    { disabled: i.disabled, admin: i.admin },
  )
}

export async function getSupportedAuthentications(): Promise<Array<SupportedAuthenticationScheme>> {
  return await ds.listSupportedAuths()
}

export async function getSupportedAuthenticationDetails(): Promise<Array<AuthSchemeDetails>> {
  const schemes = await ds.listAuthDetails()
  return schemes.map(details => ({
    schemeName: details.schemeName,
    schemeCode: details.schemeCode,
    schemeType: details.schemeType,
    schemeFlags: details.schemeFlags,
    userCount: details.userCount,
    uniqueUserCount: details.uniqueUserCount,
    labels: details.labels,
    lastUsed: details.lastUsed == null ? null : new Date(details.lastUsed), // incoming value is actually a string that we are focing into a date
  }))
}

export async function adminDeleteUser(i: {
  userSlug: string,
}) {
  await ds.adminDeleteUser(i)
}

export async function deleteGlobalAuthScheme(i: {
  schemeName: string,
}) {
  await ds.deleteGlobalAuthScheme(i)
}

export async function deleteExpiredRecoveryCodes() {
  await ds.deleteExpiredRecoveryCodes()
}

export async function getRecoveryMetrics(): Promise<RecoveryMetrics> {
  return await ds.getRecoveryMetrics()
}
