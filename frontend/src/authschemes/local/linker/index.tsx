import * as React from 'react'
import { linkLocalAccount } from '../services'
import { useForm, useFormField } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import { UserOwnView } from 'src/global_types'

const alreadyExistsErrorText = "An account for this user already exists"

export default (props: {
  onSuccess: () => void,
  userData: UserOwnView
  authFlags?: Array<string>,
}) => {
  const initialUsername = props.userData.authSchemes.find(s => s.schemeType == 'webauthn')?.username
  const username = useFormField<string>(initialUsername ?? "")
  const password = useFormField<string>('')
  const confirmPassword = useFormField<string>('')
  const [allowUsernameOverride, setOverride] = React.useState(false)

  const formComponentProps = useForm({
    fields: [password, confirmPassword],
    onSuccess: () => props.onSuccess(),
    handleSubmit: () => {
      if (username.value === '') {
        return Promise.reject("Username must be populated")
      }
      if (password.value === '') {
        return Promise.reject("Password must be populated")
      }

      try {
        return linkLocalAccount({
          username: username.value,
          password: password.value,
          confirmPassword: confirmPassword.value
        })
      }
      catch (err) {
        if ((err as Error)?.message === alreadyExistsErrorText) {
          setOverride(true)
          throw new Error("This username is taken. Please try another one.")
        }
        throw err
      }
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
      <Input type="password" label="Password" {...password} />
      <Input type="password" label="Confirm Password" {...confirmPassword} />
    </Form>
  )
}
