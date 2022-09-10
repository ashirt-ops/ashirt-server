// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useForm, useFormField, useWiredData } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import { beginLink, finishLinking } from '../services'
import { convertToCredentialCreationOptions, encodeAsB64 } from '../helpers'
import { UserOwnView } from 'src/global_types'

const alreadyExistsErrorText = "An account for this user already exists"

export default (props: {
  onSuccess: () => void,
  userData: UserOwnView
  authFlags?: Array<string>,
}) => {
  const initialUsername = props.userData.authSchemes.find(s => s.schemeType == 'local')?.username
  const username = useFormField<string>(initialUsername ?? "")
  const keyName = useFormField<string>('')
  const [allowUsernameOverride, setOverride] = React.useState(false)

  const formComponentProps = useForm({
    fields: [username, keyName],
    onSuccess: () => props.onSuccess(),
    handleSubmit: async () => {
      if (username.value === '') {
        return Promise.reject(new Error("Username must be populated"))
      }
      if (keyName.value === '') {
        return Promise.reject(new Error("Key name must be populated"))
      }

      let reg = null
      try {
        reg = await beginLink({
          username: username.value,
          keyName: keyName.value,
        })
      }
      catch (err) {
        if ((err as Error)?.message === alreadyExistsErrorText) {
          setOverride(true)
          throw new Error("This username is taken. Please try another one.")
        }
        throw err
      }

      const credOptions = convertToCredentialCreationOptions(reg)

      const signed = await navigator.credentials.create(credOptions)

      if (signed == null || signed.type != 'public-key') {
        throw new Error("WebAuthn is not supported")
      }
      const pubKeyCred = signed as PublicKeyCredential
      const pubKeyResponse = pubKeyCred.response as AuthenticatorAttestationResponse

      await finishLinking({
        type: 'public-key',
        id: pubKeyCred.id,
        rawId: encodeAsB64(pubKeyCred.rawId),
        response: {
          attestationObject: encodeAsB64(pubKeyResponse.attestationObject),
          clientDataJSON: encodeAsB64(pubKeyResponse.clientDataJSON),
        },
      })
    }
  })

  // yes, this could be rewritten as !allowUsernameOverride && (initialUsername !== undefined)
  // but the variable naming becomes weird, so using a different solution
  const readonlyUsername = allowUsernameOverride ?
    false :
    (initialUsername !== undefined)

  return (
    <Form submitText="Link Account" {...formComponentProps}>
      <Input label="Username" {...username} disabled={readonlyUsername} />
      <Input label="Key name" {...keyName} />
    </Form>
  )
}
