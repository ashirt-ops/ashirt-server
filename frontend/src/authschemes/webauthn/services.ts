// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from 'src/services/data_sources/backend/request_helper'

import {
  CompletedLoginChallenge,
  KeyEntry,
  KeyList,
  ProvidedCredentialCreationOptions,
  ProvidedCredentialRequestOptions,
  WebAuthNRegisterConfirmation,
} from "./types"

export async function beginRegistration(i: {
  email: string,
  username: string,
  firstName: string,
  lastName: string
  credentialName: string
}): Promise<ProvidedCredentialCreationOptions> {
  return await req('POST', '/auth/webauthn/register/begin', i)
}

export async function finishRegistration(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/register/finish', i)
}

export async function beginLogin(i: {
  username: string,
}): Promise<ProvidedCredentialRequestOptions> {
  return await req('POST', '/auth/webauthn/login/begin', i)
}

export async function finishLogin(i: CompletedLoginChallenge): Promise<void> {
  return await req('POST', '/auth/webauthn/login/finish', i)
}

export async function beginLink(i: {
  username: string,
  credentialName: string
}): Promise<ProvidedCredentialCreationOptions> {
  return await req('POST', '/auth/webauthn/link/begin', i)
}

export async function finishLinking(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/link/finish', i)
}

export async function beginAddKey(i: {
  credentialName: string
}): Promise<ProvidedCredentialCreationOptions> {
  return await req('POST', '/auth/webauthn/key/add/begin', i)
}

export async function finishAddKey(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/key/add/finish', i)
}

export async function listWebauthnKeys(): Promise<KeyList> {
  const data: KeyList = await req('GET', '/auth/webauthn/keys')

  return {
    keys: data.keys.map((key: KeyEntry) => ({
      ...key,
      dateCreated: new Date(key.dateCreated)
    }))
  }

}

export async function deleteWebauthnKey(i: { credentialName: string }): Promise<KeyList> {
  return await req('DELETE', `/auth/webauthn/key/${i.credentialName}`)
}
