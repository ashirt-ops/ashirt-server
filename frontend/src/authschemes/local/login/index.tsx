import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import classnames from 'classnames/bind'
import { login, register, userResetPassword, totpLogin } from '../services'
import { useForm, useFormField } from 'src/helpers/use_form'
import { useModal, renderModals, OnRequestClose } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

async function handleLoginStepPromise(promise: Promise<void>): Promise<void> {
  try {
    await promise
    window.location.href = '/'
  } catch (err) {
    if (err.message === 'PASSWORD_RESET_REQUIRED') {
      window.location.href = '/login/local?step=reset'
      return
    }
    else if (err.message === 'TOTP_REQUIRED') {
      window.location.href = '/login/local?step=totp'
      return
    }
    throw err
  }
}

// Used to pull a value out of a password field and clear the field value for security
function getValueAndClear(field: { value: string, onChange: (s: string) => void }): string {
  const { value } = field
  field.onChange('')
  return value
}

export default (props: {
  query: URLSearchParams,
  authFlags?: Array<string>
}) => {
  switch (props.query.get('step')) {
    case 'reset': return <ResetPassword />
    case 'totp': return <EnterTotp />
    default: return <Login authFlags={props.authFlags} />
  }
}

const Login = (props: {
  authFlags?: Array<string>
}) => {
  const usernameField = useFormField('')
  const passwordField = useFormField('')

  const loginForm = useForm({
    fields: [usernameField, passwordField],
    handleSubmit: () => (
      handleLoginStepPromise(login(usernameField.value, getValueAndClear(passwordField)))
    ),
  })

  const registerModal = useModal<void>(modalProps => <RegisterModal {...(modalProps as OnRequestClose)} />)

  const allowRegister = props.authFlags?.includes("open-registration")
  const registerProps = allowRegister
    ? { cancelText: "Register", onCancel: () => registerModal.show() }
    : {}

  return (
    <div style={{ minWidth: 300 }}>
      <Form submitText="Login" {...registerProps} {...loginForm}>
        <Input label="Username" autoFocus {...usernameField} />
        <Input label="Password" type="password" {...passwordField} />
      </Form>
      {renderModals(registerModal)}
    </div>
  )
}

const RegisterModal = (props: {
  onRequestClose: () => void,
}) => {
  const firstNameField = useFormField('')
  const lastNameField = useFormField('')
  const usernameField = useFormField('')
  const emailField = useFormField('')
  const passwordField = useFormField('')
  const confirmPasswordField = useFormField('')

  const registerForm = useForm({
    fields: [
      firstNameField,
      lastNameField,
      usernameField,
      emailField,
      passwordField,
      confirmPasswordField,
    ],
    handleSubmit: async () => {
      await register({
        firstName: firstNameField.value,
        lastName: lastNameField.value,
        username: usernameField.value,
        email: emailField.value,
        password: getValueAndClear(passwordField),
        confirmPassword: getValueAndClear(confirmPasswordField),
      })
      window.location.pathname = '/'
    },
  })

  return (
    <Modal title="Register" onRequestClose={props.onRequestClose}>
      <Form submitText="Create Account" {...registerForm}>
        <Input label="First Name" {...firstNameField} />
        <Input label="Last Name" {...lastNameField} />
        <Input label="Contact Email" {...emailField} />
        <Input label="Desired Username" {...usernameField} />
        <Input label="Password" type="password" {...passwordField} />
        <Input label="Confirm Password" type="password" {...confirmPasswordField} />
      </Form>
    </Modal>
  )
}

const ResetPassword = (props: {
}) => {
  const passwordField = useFormField('')
  const confirmPasswordField = useFormField('')

  const resetPasswordForm = useForm({
    fields: [passwordField, confirmPasswordField],
    handleSubmit: () => (
      handleLoginStepPromise(userResetPassword({
        newPassword: getValueAndClear(passwordField),
        confirmPassword: getValueAndClear(confirmPasswordField),
      }))
    ),
  })

  return <>
    <div className={cx('messagebox')}>
      You have been given a temporary password. You must change this password before you can continue using this application.
    </div>
    <Form submitText="Update Password" {...resetPasswordForm}>
      <Input label="New Password" type="password" {...passwordField} />
      <Input label="Confirm New Password" type="password" {...confirmPasswordField} />
    </Form>
  </>
}

const EnterTotp = (props: {}) => {
  const totpField = useFormField('')

  const totpForm = useForm({
    fields: [totpField],
    handleSubmit: () => handleLoginStepPromise(
      totpLogin(totpField.value)
    ),
  })

  return (<>
    <h2 className={cx('title')}>Multi-factor Authentication</h2>
    <Form submitText="Submit" {...totpForm}>
      <Input label="Passcode" {...totpField} />
    </Form>
  </>)
}
