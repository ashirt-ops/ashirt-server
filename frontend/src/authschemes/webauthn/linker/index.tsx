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
  const email = useFormField<string>(props.userData.email)
  const keyName = useFormField<string>('')

  const formComponentProps = useForm({
    fields: [email, keyName],
    onSuccess: () => props.onSuccess(),
    handleSubmit: async () => {
      if (email.value === '') {
        return Promise.reject(new Error("Email must be populated"))
      }
      if (keyName.value === '') {
        return Promise.reject(new Error("Key name must be populated"))
      }

      const reg = await beginLink({
        email: email.value,
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

  return (
    <Form submitText="Link Account" {...formComponentProps}>
      <Input label="Email" {...email} />
      <Input label="Key name" {...keyName} />
    </Form>
  )
}
