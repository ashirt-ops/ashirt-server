
/**
 * ProvidedCredentialCreationOptions tries to be a mirror of CredentialCreationOptions, but
 * with all base64 encoded fields marked as string
 *
 * This seems to be the cleanest version that extends CredentialCreationOptions that typescript
 * supports. The definition amounts to "copy this defintion, but exclude this bit", then
 * "tack on this new bit", which ultimately replaces the type from the original type
 */
export type ProvidedCredentialCreationOptions =
  Omit<CredentialCreationOptions, "publicKey">
  & {
    publicKey: Omit<PublicKeyCredentialCreationOptions, "challenge" | "excludeCredentials" | "user"> &
    {
      challenge: string // base64 encoded string (url encoded)
      excludeCredentials?: Array<Omit<PublicKeyCredentialDescriptor, "id"> & {
        id: string // base64 encoded string (url encoded)
      }>
      user: Omit<PublicKeyCredentialUserEntity, "id"> & {
        id: string // base64 encoded string (url encoded)
      }
    }
  }

export type ProvidedCredentialRequestOptions = {
  publicKey: Omit<PublicKeyCredentialRequestOptions, "challenge" | "allowCredentials"> & {
    challenge: string
    allowCredentials: Array<Omit<PublicKeyCredentialDescriptor, "id"> & {
      id: string
    }>
  }
}

export type WebAuthNRegisterConfirmation = {
  id: string
  rawId: string // base64
  type: 'public-key'
  response: {
    clientDataJSON: string // base64
    attestationObject: string // base64
  }
}

export type CompletedLoginChallenge = {
  id: string
  type: string
  rawId: string
  response: {
    authenticatorData: string
    clientDataJSON: string
    signature: string
    userHandle: string
  }
}

export type CredentialList = {
  credentials: Array<CredentialEntry>
}

export type CredentialEntry = {
  credentialName: string
  dateCreated: Date
  credentialId: string
}
