// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useForm, useFormField, useWiredData } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import { beginLink, finishLinking } from '../services'
import { convertToCredentialCreationOptions, encodeAsB64 } from '../helpers'
import { UserOwnView } from 'src/global_types'

export default (props: {
  onSuccess: () => void,
  userData: UserOwnView
  authFlags?: Array<string>,
}) => {
  const initialUsername = props.userData.authSchemes.find(s => s.schemeType == 'local')?.username
  const username = useFormField<string>(initialUsername ?? "")
  const keyName = useFormField<string>('')

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

      const reg = await beginLink({
        username: username.value,
        keyName: keyName.value,
      })
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

  const readonlyUsername = initialUsername !== undefined

  return (
    <Form submitText="Link Account" {...formComponentProps}>
      <Input label="Username" {...username} readOnly={readonlyUsername} />
      <Input label="Key name" {...keyName} />
    </Form>
  )
}
