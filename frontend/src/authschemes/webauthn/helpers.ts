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

const base64UrlToBase64 = (input: string) => {
  // Replace non-url compatible chars with base64 standard chars
  let standardCharInput = input
      .replace(/-/g, '+')
      .replace(/_/g, '/');

  // Pad out with standard base64 required padding characters
  const pad = standardCharInput.length % 4;
  if(pad) {
    if(pad === 1) {
      throw new Error('InvalidLengthError: Input base64url string is the wrong length to determine padding');
    }
    standardCharInput += new Array(5-pad).join('=');
  }
  return standardCharInput;
}

export const toByteArrayFromB64URL = (s: string) => Uint8Array.from(atob(base64UrlToBase64(s)), toKeycode)

export const convertToCredentialCreationOptions = (
  input: ProvidedCredentialCreationOptions
): CredentialCreationOptions => {

  const output: CredentialCreationOptions = {
    ...input,
    publicKey: {
      ...input.publicKey,
      challenge: toByteArrayFromB64URL(input.publicKey.challenge),
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
    challenge: toByteArrayFromB64URL(input.publicKey.challenge),
    allowCredentials: input.publicKey.allowCredentials.map(
      listItem => ({ ...listItem, id: toByteArray(listItem.id) })
    )
  }

  return output
}
