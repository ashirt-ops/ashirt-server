// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import classnames from 'classnames/bind'
import { useForm, useFormField } from 'src/helpers/use_form'
import { userChangePassword } from '../services'
import Totp from '../totp'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  username: string,
  authFlags?: Array<string>
}) => (
  <>
    <h1 className={cx('header')}>Settings for local account <span className={cx('user-key')}>{props.username}</span></h1>
    <SettingsSection title="Change Password" width="narrow">
      <ResetPasswordForm username={props.username} />
    </SettingsSection>
    <Totp />
  </>
)
const ResetPasswordForm = (props: {
  username: string,
}) => {
  const oldPassword = useFormField('')
  const newPassword = useFormField('')
  const confirmPassword = useFormField('')

  const profileForm = useForm({
    fields: [oldPassword, newPassword, confirmPassword],
    handleSubmit: () => {
      return userChangePassword({
        username: props.username,
        oldPassword: oldPassword.value,
        newPassword: newPassword.value,
        confirmPassword: confirmPassword.value,
      })
    },
    onSuccessText: 'Password updated',
  })

  return (
    <div>
      <Form submitText="Update Password" {...profileForm}>
        <Input type="password" label="Old Password" {...oldPassword} />
        <Input type="password" label="New Password" {...newPassword} />
        <Input type="password" label="Confirm New Password" {...confirmPassword} />
      </Form>
    </div>
  )
}
