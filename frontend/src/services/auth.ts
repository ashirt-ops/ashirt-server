// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import { UserOwnView, SupportedAuthenticationScheme, AuthenticationInfo, AuthSchemeDetails, RecoveryMetrics } from 'src/global_types'

export async function getCurrentUser(): Promise<UserOwnView | null> {
  const res = await fetch('/web/user')
  if ([200, 404].indexOf(res.status) === -1) throw Error(`/web/user returned status ${res.status}`)
  const csrfToken = res.headers.get('X-CSRF-TOKEN')
  if (csrfToken == null) throw Error(`/web/user returned status ${res.status}, but no csrf token`)
  if (res.status !== 200) return null
  return await res.json()
}

// TODO this should be encapsulated in an admin settings component under src/authschemes/local
export async function adminChangePassword(i: {
  userSlug: string,
  newPassword: string,
}) {
  if (i.newPassword.length < 3) {
    throw Error("User password must be at least 3 characters long")
  }
  await req('PUT', '/auth/local/admin/password', i)
}

export async function logout() {
  return req('POST', '/logout', {})
}

export async function getUser(i?: { userSlug: string }): Promise<UserOwnView> {
  const user = await req('GET', `/user`, null, i)

  return {
    email: user.email,
    authSchemes: user.authSchemes.map((scheme: AuthenticationInfo) => ({
      userKey: scheme.userKey,
      schemeCode: scheme.schemeCode,
      lastLogin: scheme.lastLogin == null ? null : new Date(scheme.lastLogin), // last login is actually a string here

    })),
    admin: user.admin,
    headless: user.headless,
    slug: user.slug,
    firstName: user.firstName,
    lastName: user.lastName,
  }
}

export async function adminSetUserFlags(i: {
  userSlug: string,
  disabled: boolean,
  admin: boolean,
}) {
  await req('POST', `/admin/${i.userSlug}/flags`, i)
}

export async function getSupportedAuthentications(): Promise<Array<SupportedAuthenticationScheme>> {
  return await req('GET', '/auths')
}

export async function getSupportedAuthenticationDetails(): Promise<Array<AuthSchemeDetails>> {
  const schemes = await req('GET', '/auths/breakdown')
  return schemes.map((i: AuthSchemeDetails) => ({
    schemeName: i.schemeName,
    schemeCode: i.schemeCode,
    userCount: i.userCount,
    uniqueUserCount: i.uniqueUserCount,
    labels: i.labels,
    lastUsed: i.lastUsed == null ? null : new Date(i.lastUsed), // incoming value is actually a string that we are focing into a date
  }))
}

export async function adminDeleteUser(i: {
  userSlug: string,
}) {
  await req('DELETE', `/admin/user/${i.userSlug}`, i)
}

export async function deleteGlobalAuthScheme(i: {
  schemeName: string,
}) {
  await req('DELETE', `/auths/${i.schemeName}`)
}

export async function deleteExpiredRecoveryCodes() {
  await req('DELETE', '/auth/recovery/expired')
}

export async function getRecoveryMetrics(): Promise<RecoveryMetrics> {
  return await req('GET', '/auth/recovery/metrics')
}
