import req from 'src/services/data_sources/backend/request_helper'

import {
  CompletedLoginChallenge,
  CredentialEntry,
  CredentialList,
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
}, discoverable: boolean): Promise<ProvidedCredentialCreationOptions> {
  return await req('POST', '/auth/webauthn/register/begin', i, { discoverable })
}

export async function finishRegistration(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/register/finish', i)
}

export async function beginLogin(i: {
  username: string,
}, discoverable: boolean): Promise<ProvidedCredentialRequestOptions> {
  return await req('POST', '/auth/webauthn/login/begin', i, { discoverable })
}

export async function finishLogin(i: CompletedLoginChallenge, discoverable: boolean): Promise<void> {
  return await req('POST', '/auth/webauthn/login/finish', i, { discoverable })
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

export async function beginAddCredential(i: {
  credentialName: string
}): Promise<ProvidedCredentialCreationOptions> {
  return await req('POST', '/auth/webauthn/credential/add/begin', i)
}

export async function finishAddCredential(i: WebAuthNRegisterConfirmation) {
  return await req('POST', '/auth/webauthn/credential/add/finish', i)
}

export async function listWebauthnCredentials(): Promise<CredentialList> {
  const data: CredentialList = await req('GET', '/auth/webauthn/credentials')

  return {
    credentials: data.credentials.map((credential: CredentialEntry) => ({
      ...credential,
      dateCreated: new Date(credential.dateCreated)
    }))
  }

}

export async function deleteWebauthnCredential(i: { credentialId: string }): Promise<CredentialList> {
  return await req('DELETE', `/auth/webauthn/credential/${i.credentialId}`)
}

export async function modifyCredentialName(i: { credentialName: string, newCredentialName: string }): Promise<void> {
  return await req('PUT', `/auth/webauthn/credential`, i)
}
