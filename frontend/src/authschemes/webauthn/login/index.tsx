import * as React from 'react'
import Button, { ButtonGroup } from 'src/components/button'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import { beginLogin, beginRegistration, finishLogin, finishRegistration } from '../services'
import { useForm, useFormField } from 'src/helpers/use_form'
import { useModal, renderModals, OnRequestClose } from 'src/helpers'
import { convertToCredentialCreationOptions, convertToPublicKeyCredentialRequestOptions, encodeAsB64 } from '../helpers'
import { getResultState } from 'src/helpers/is_success_result'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))


export default (props: {
  query: URLSearchParams,
  authFlags?: Array<string>,
}) => {
  return (
    <Login authFlags={props.authFlags} />
  )
}

const Login = (props: {
  authFlags?: Array<string>,
}) => {
  const usernameField = useFormField('')
  const [isDiscoverable, setIsDiscoverable] = React.useState(true)

  const loginForm = useForm({
    fields: [usernameField],
    handleSubmit: async () => {
      const protoOptions = await beginLogin({ username: usernameField.value }, isDiscoverable)
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
      }, isDiscoverable)
      window.location.pathname = '/'
    },
  })

  const registerModal = useModal<void>(modalProps => <RegisterModal onRequestClose={() => modalProps.onRequestClose()} isDiscoverable={isDiscoverable} />)

  const allowRegister = props.authFlags?.includes("open-registration") // TODO: this isn't being used

  const registerProps = allowRegister
    ? { cancelText: "Register", onCancel: () => registerModal.show() }
    : {}

  const makeDiscoverable = () => setIsDiscoverable(true)
  const makeNonDiscoverable = () => setIsDiscoverable(false)

  return (
    <div>
      {window.PublicKeyCredential && (
        <div className={cx('login-container')}>
          <div className={cx('mode-buttons')}>
            <ButtonGroup className={cx('row-buttons')}>
              <Button active={isDiscoverable} className={cx('mode-button-right')} onClick={makeDiscoverable}>Discoverable</Button>
              <Button active={!isDiscoverable} className={cx('mode-button-left')} onClick={makeNonDiscoverable}>Username</Button>
            </ButtonGroup>
          </div>
          {isDiscoverable ? (
            <div>
            <Form submitText="Login" {...registerProps} {...loginForm} autoFocus={true}>
            </Form>
            {renderModals(registerModal)}
          </div>
          ) : (
            <div>
            <Form submitText="Login" {...registerProps} {...loginForm}>
              <Input label="Username" {...usernameField} />
            </Form>
            {renderModals(registerModal)}
          </div>
          )}
        </div>
      )}
    </div>
  )
}

const RegisterModal = (props: {
  onRequestClose: () => void,
  isDiscoverable: boolean,
}) => {
  const firstNameField = useFormField('')
  const lastNameField = useFormField('')
  const emailField = useFormField('')
  const usernameField = useFormField('')
  const keyNameField = useFormField('')

  const registerForm = useForm({
    onSuccessText: "Successfully registered",
    fields: [
      firstNameField,
      lastNameField,
      emailField,
      usernameField,
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
        username: usernameField.value,
        credentialName: keyNameField.value,
      }, props.isDiscoverable)
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
        <Input label="Desired Username" {...usernameField} />
        <Input label="Credential Name" {...keyNameField} />
      </Form>
    </Modal>
  )
}
