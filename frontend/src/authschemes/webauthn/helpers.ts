// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { ProvidedCredentialCreationOptions, ProvidedCredentialRequestOptions } from "./types"

export const encodeAsB64 = (ab: ArrayBuffer) => {
  return base64UrlEncode(arrayBufferToString(ab))
}

export const arrayBufferToString = (a: ArrayBuffer) => {
  return String.fromCharCode.apply(null, new Uint8Array(a))
}

export const base64UrlEncode = (value: string) => {
  return btoa(value)
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "")
}

export const toKeycode = (c: string) => c.charCodeAt(0)

export const toByteArray = (s: string) => Uint8Array.from(atob(s), toKeycode)

export const convertToCredentialCreationOptions = (
  input: ProvidedCredentialCreationOptions
): CredentialCreationOptions => {

  const output: CredentialCreationOptions = {
    ...input,
    publicKey: {
      ...input.publicKey,
      challenge: toByteArray(input.publicKey.challenge),
      user: {
        ...input.publicKey.user,
        id: toByteArray(input.publicKey.user.id)
      },
      excludeCredentials: input.publicKey.excludeCredentials?.map(
        cred => ({ ...cred, id: toByteArray(cred.id) })
      )
    }
  }

  return output
}

export const convertToPublicKeyCredentialRequestOptions = (input: ProvidedCredentialRequestOptions): PublicKeyCredentialRequestOptions => {
  const output: PublicKeyCredentialRequestOptions = {
    ...input.publicKey,
    challenge: toByteArray(input.publicKey.challenge),
    allowCredentials: input.publicKey.allowCredentials.map(
      listItem => ({ ...listItem, id: toByteArray(listItem.id) })
    )
  }

  return output
}
