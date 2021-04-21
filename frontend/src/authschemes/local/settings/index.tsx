// Copyright 2020, Verizon Media
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
  userKey: string,
  authFlags?: Array<string>
}) => <>
  <h1 className={cx('header')}>Settings for local account <span className={cx('user-key')}>{props.userKey}</span></h1>
  <SettingsSection title="Change Password" width="narrow">
    <ResetPasswordForm userKey={props.userKey} />
  </SettingsSection>
  <Totp />
</>

const ResetPasswordForm = (props: {
  userKey: string,
}) => {
  const oldPassword = useFormField('')
  const newPassword = useFormField('')
  const confirmPassword = useFormField('')

  const profileForm = useForm({
    fields: [oldPassword, newPassword, confirmPassword],
    handleSubmit: () => {
      return userChangePassword({
        userKey: props.userKey,
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
