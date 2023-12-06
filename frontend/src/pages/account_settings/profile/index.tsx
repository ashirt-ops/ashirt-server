import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import { UserWithAuth } from 'src/global_types'
import { updateUserProfile } from 'src/services'
import { useForm, useFormField } from 'src/helpers'

export default (props: {
  profile: UserWithAuth
  requestReload?: () => void
}) => {

  const firstNameField = useFormField(props.profile.firstName)
  const lastNameField = useFormField(props.profile.lastName)
  const emailField = useFormField(props.profile.email)

  const profileForm = useForm({
    fields: [firstNameField, lastNameField, emailField],
    handleSubmit: () => updateUserProfile({
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
        <Input type="email" name="contactEmail" label="Contact Email" {...emailField} />
      </Form>
    </SettingsSection>
  )
}
