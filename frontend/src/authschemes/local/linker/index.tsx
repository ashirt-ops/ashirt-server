// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { linkLocalAccount } from '../services'
import { useForm, useFormField } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'

export default (props: {
  onSuccess: () => void,
  authFlags?: Array<string>,
}) => {
  const username = useFormField<string>('')
  const password = useFormField<string>('')
  const confirmPassword = useFormField<string>('')

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

      return linkLocalAccount({
        username: username.value,
        password: password.value,
        confirmPassword: confirmPassword.value
      })
    }
  })

  return (
    <Form submitText="Link Account" {...formComponentProps}>
      <Input label="Username" {...username} />
      <Input type="password" label="Password" {...password} />
      <Input type="password" label="Confirm Password" {...confirmPassword} />
    </Form>
  )
}
