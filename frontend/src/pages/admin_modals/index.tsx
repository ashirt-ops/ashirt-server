// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { User, UserAdminView } from 'src/global_types'
import {
  adminChangePassword, adminSetUserFlags, adminDeleteUser, addHeadlessUser,
  deleteGlobalAuthScheme, deleteTotpForUser, adminCreateLocalUser
} from 'src/services'
import AuthContext from 'src/auth_context'
import Button from 'src/components/button'
import ChallengeModalForm from 'src/components/challenge_modal_form'
import Checkbox from 'src/components/checkbox'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import Form from 'src/components/form'
import ModalForm from 'src/components/modal_form'
import { InputWithCopyButton } from 'src/components/text_copiers'
import { useForm, useFormField } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

export const ResetPasswordModal = (props: {
  user: User,
  onRequestClose: () => void,
}) => {
  const tempPassword = useFormField<string>("")
  const formComponentProps = useForm({
    fields: [tempPassword],
    onSuccess: () => props.onRequestClose(),
    handleSubmit: () => adminChangePassword({
      userSlug: props.user.slug,
      newPassword: tempPassword.value,
    }),
  })
  return (
    <ModalForm title="Set Temporary Password" submitText="Update" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <Input label="New Temporary Password" {...tempPassword} />
    </ModalForm>
  )
}

export const AddHeadlessUserModal = (props: {
  onRequestClose: () => void,
}) => {
  const headlessName = useFormField<string>("")
  const contactEmail = useFormField<string>("")
  const formComponentProps = useForm({
    fields: [headlessName, contactEmail],
    onSuccess: () => props.onRequestClose(),
    handleSubmit: () => {
      if (headlessName.value.length == 0) {
        return new Promise((resolve, reject) => reject(Error("Headless users must be given a name")))
      }
      return addHeadlessUser({
        firstName: 'Headless',
        lastName: headlessName.value,
        email: contactEmail.value,
      })
    },
  })
  return (
    <ModalForm title="Create New Headless User" submitText="Create" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <Input label="Headless name" {...headlessName} />
      <Input type="email" label="Contact Email" {...contactEmail} />
    </ModalForm>
  )
}

export const AddUserModal = (props: {
  onRequestClose: () => void,
}) => {
  const firstName = useFormField<string>("")
  const lastName = useFormField<string>("")
  const contactEmail = useFormField<string>("")

  const [username, setUsername] = React.useState<string>("")
  const [password, setPassword] = React.useState<string>("")
  const [isDisabled, setDisabled] = React.useState<boolean>(false)

  const formComponentProps = useForm({
    fields: [firstName, lastName, contactEmail],
    handleSubmit: () => {
      if (firstName.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("Users should have at least a first name")))
      }
      if (contactEmail.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("Users must have an email address")))
      }
      // TODO: this should create the user, then update the form with the new user/password combo
      // to share.
      const runSubmit = async () => {
        const result = await adminCreateLocalUser({
          firstName: firstName.value,
          lastName: lastName.value,
          email: contactEmail.value,
        })
        setUsername(contactEmail.value)
        setPassword(result.temporaryPassword)
        setDisabled(true) // lock the form -- we don't need to allow submits at this time.
      }

      return runSubmit()
    },
  })

  return (
    <Modal title="Create New User" onRequestClose={props.onRequestClose}>
      <Form {...formComponentProps} loading={isDisabled}
        submitText={isDisabled ? undefined : "Submit"}
      >
        <Input label="First Name" {...firstName} disabled={isDisabled} />
        <Input label="Last Name" {...lastName} disabled={isDisabled} />
        <Input label="Email" {...contactEmail} disabled={isDisabled} />
      </Form>
      {isDisabled && (<>
        <div className={cx('success-area')}>
          <p>Below is the new user's initial login credentials:</p>
          <InputWithCopyButton label="Username" value={username} />
          <InputWithCopyButton label="Password" value={password} />
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      }
    </Modal>
  )
}

export const UpdateUserFlagsModal = (props: {
  user: UserAdminView,
  onRequestClose: () => void,
}) => {
  const fullContext = React.useContext(AuthContext)
  const adminSlug = fullContext.user ? fullContext.user.slug : ""

  const isAdmin = useFormField(props.user.admin)
  const isDisabled = useFormField(props.user.disabled)

  const formComponentProps = useForm({
    fields: [isAdmin, isDisabled],
    onSuccess: () => props.onRequestClose(),
    handleSubmit: () => {
      return adminSetUserFlags({ userSlug: props.user.slug, disabled: isDisabled.value, admin: isAdmin.value })
    }
  })

  const badAdmin = { disabled: true, title: "Admins cannot alter this flag on themselves" }
  const adminIsTargetUser = adminSlug === props.user.slug
  const isHeadlessUser = props.user.headless
  const canAlterDisabled = adminIsTargetUser ? badAdmin : {}
  const canAlterAdmin = () => {
    if (adminIsTargetUser) { return badAdmin }
    if (isHeadlessUser) { return { disabled: true, title: "Headless users cannot be admins" } }
    return {}
  }

  const mergedAdminProps = { ...isAdmin, ...canAlterAdmin() }
  const mergedDisabledProps = { ...isDisabled, ...canAlterDisabled }

  return <ModalForm title="Set User Flags" submitText="Update" onRequestClose={props.onRequestClose} {...formComponentProps}>
    <em className={cx('warning')}>Updating these values will log out the user</em>
    <Checkbox label="Admin" {...mergedAdminProps} />
    <Checkbox label="Disabled" {...mergedDisabledProps} />
  </ModalForm>
}

export const DeleteUserModal = (props: {
  user: UserAdminView,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Delete User"
    warningText="This will remove the user from the system. All user information will be lost."
    submitText="Delete"
    challengeText={props.user.slug}
    handleSubmit={() => adminDeleteUser({ userSlug: props.user.slug })}
    onRequestClose={props.onRequestClose}
  />

export const DeleteGlobalAuthSchemeModal = (props: {
  schemeCode: string,
  uniqueUsers: number,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Remove Users from Authentication Scheme"
    warningText={`This will unlink/remove this authentication scheme from all users.${
      props.uniqueUsers == 0 ? "" : ` Note that this will effectively disable ${props.uniqueUsers} accounts.`
      }`}
    submitText="Remove All"
    challengeText={props.schemeCode}
    handleSubmit={() => deleteGlobalAuthScheme({ schemeName: props.schemeCode })}
    onRequestClose={props.onRequestClose}
  />

export const RecoverAccountModal = (props: {
  recoveryCode: string
  onRequestClose: () => void,
}) => {
  const url = `${window.location.origin}/web/auth/recovery/login?code=${props.recoveryCode}`
  return <Modal title="Recovery URL" onRequestClose={props.onRequestClose}>
    <div className={cx('recovery-code-modal')}>
      <p>
        Below is the recovery URL. Provide this to the user, and they will be able
        to log in without the need to authenticate.
      </p>
      <InputWithCopyButton label="Recovery URL" value={url} />
      <Button primary onClick={() => props.onRequestClose()}>Close</Button>
    </div>
  </Modal>
}

export const RemoveTotpModal = (props: {
  user: UserAdminView,
  onRequestClose: () => void,
}) => {
  const formComponentProps = useForm({
    fields: [],
    onSuccess: () => props.onRequestClose(),
    handleSubmit: () => deleteTotpForUser({userSlug: props.user.slug}),
  })

  return <ModalForm title="Disable Multi-Factor Authentication" submitText="Continue" onRequestClose={props.onRequestClose} {...formComponentProps}>
    <em className={cx('warning')}>
      Multi-factor Authentication provides an extra layer of security for this user.
      Removing this factor should only be done if the user has lost the device or the mechansim to authenticate.
      Are you sure you want to continue?
    </em>
  </ModalForm>
}
