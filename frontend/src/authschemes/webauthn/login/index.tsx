// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import { beginLogin, beginRegistration, finishLogin, finishRegistration } from '../services'
import { useForm, useFormField } from 'src/helpers/use_form'
import { useModal, renderModals, OnRequestClose } from 'src/helpers'
import { convertToCredentialCreationOptions, convertToPublicKeyCredentialRequestOptions, encodeAsB64 } from '../helpers'
import { getResultState } from 'src/helpers/is_success_result'


export default (props: {
  query: URLSearchParams,
  authFlags?: Array<string>
}) => {
  return (
    <Login authFlags={props.authFlags} />
  )
}

const Login = (props: {
  authFlags?: Array<string>
}) => {
  const emailField = useFormField('')

  const loginForm = useForm({
    fields: [emailField],
    handleSubmit: async () => {
      const protoOptions = await beginLogin({ email: emailField.value })
      const credOptions = convertToPublicKeyCredentialRequestOptions(protoOptions)

      const cred = await navigator.credentials.get({
        publicKey: credOptions
      })
      if (cred == null || cred.type != 'public-key') {
        throw new Error("WebAuthn is not supported")
      }
      const pubKeyCred = cred as PublicKeyCredential
      const pubKeyResponse = pubKeyCred.response as AuthenticatorAssertionResponse

      await finishLogin({
        id: pubKeyCred.id,
        rawId: encodeAsB64(pubKeyCred.rawId),
        type: pubKeyCred.type,
        response: {
          authenticatorData: encodeAsB64(pubKeyResponse.authenticatorData),
          clientDataJSON: encodeAsB64(pubKeyResponse.clientDataJSON),
          signature: encodeAsB64(pubKeyResponse.signature),
          userHandle: pubKeyResponse.userHandle == null ? "" : encodeAsB64(pubKeyResponse.userHandle),
        }
      })
      window.location.pathname = '/' // TODO: is this what we want to do?
    },
  })

  const registerModal = useModal<void>(modalProps => <RegisterModal {...(modalProps as OnRequestClose)} />)

  const allowRegister = props.authFlags?.includes("open-registration")
  // const registerProps = allowRegister
  //   ? { cancelText: "Register", onCancel: () => registerModal.show() }
  //   : {}

  const registerProps = { cancelText: "Register", onCancel: () => registerModal.show() }

  return (
    <div>
      {window.PublicKeyCredential && (
        <div style={{ minWidth: 300 }}>
          <Form submitText="Login with WebAuthN" {...registerProps} {...loginForm}>
            <Input label="Email" {...emailField} />
          </Form>
          {renderModals(registerModal)}
        </div>
      )}
    </div>
  )
}

const RegisterModal = (props: {
  onRequestClose: () => void,
}) => {
  const firstNameField = useFormField('')
  const lastNameField = useFormField('')
  const emailField = useFormField('')
  const keyNameField = useFormField('')

  const registerForm = useForm({
    onSuccessText: "Successfully registered",
    fields: [
      firstNameField,
      lastNameField,
      emailField,
      keyNameField,
    ],
    handleSubmit: async () => {
      if (getResultState(registerForm.result) === 'success' ) {
        props.onRequestClose()
        return
      }
      const reg = await beginRegistration({
        firstName: firstNameField.value,
        lastName: lastNameField.value,
        email: emailField.value,
        keyName: keyNameField.value,
      })
      const credOptions = convertToCredentialCreationOptions(reg)

      const signed = await navigator.credentials.create(credOptions)

      if (signed == null || signed.type != 'public-key') {
        throw new Error("WebAuthn is not supported")
      }
      const pubKeyCred = signed as PublicKeyCredential
      const pubKeyResponse = pubKeyCred.response as AuthenticatorAttestationResponse

      await finishRegistration({
        type: 'public-key',
        id: pubKeyCred.id,
        rawId: encodeAsB64(pubKeyCred.rawId),
        response: {
          attestationObject: encodeAsB64(pubKeyResponse.attestationObject),
          clientDataJSON: encodeAsB64(pubKeyResponse.clientDataJSON),
        },
      })
    },
  })

  return (
    <Modal title="Register" onRequestClose={props.onRequestClose}>
      <Form
        submitText={getResultState(registerForm.result) == 'success' ? "Close" : "Create Account"}
        {...registerForm}
      >
        <Input label="First Name" {...firstNameField} />
        <Input label="Last Name" {...lastNameField} />
        <Input label="Email" {...emailField} />
        <Input label="Key Name" {...keyNameField} />
      </Form>
    </Modal>
  )
}
