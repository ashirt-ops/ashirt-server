// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { UserWithAuth } from 'src/global_types'
import { useDataSource, updateUserProfile } from 'src/services'
import { useForm, useFormField } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'

export default (props: {
  profile: UserWithAuth
  requestReload?: () => void
}) => {
  const ds = useDataSource()

  const firstNameField = useFormField(props.profile.firstName)
  const lastNameField = useFormField(props.profile.lastName)
  const emailField = useFormField(props.profile.email)

  const profileForm = useForm({
    fields: [firstNameField, lastNameField, emailField],
    handleSubmit: () => updateUserProfile(ds, {
      userSlug: props.profile.slug,
      firstName: firstNameField.value,
      lastName: lastNameField.value,
      email: emailField.value,
    }),
    onSuccess: props.requestReload,
    onSuccessText: 'Profile updated',
  })

  return (
    <SettingsSection title="Profile Settings" width="narrow">
      <Form submitText="Update Profile" {...profileForm}>
        <Input name="firstName" label="First Name" {...firstNameField} />
        <Input name="lastName" label="Last Name" {...lastNameField} />
        <Input name="contactEmail" label="Contact Email" {...emailField} />
      </Form>
    </SettingsSection>
  )
}
