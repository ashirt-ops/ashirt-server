// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React  from 'react'
import classnames from 'classnames/bind'

import { ApiKey, GlobalVar, User, UserAdminView, UserGroupAdminView } from 'src/global_types'
import {
  adminChangePassword, adminSetUserFlags, adminDeleteUser, addHeadlessUser,
  deleteGlobalAuthScheme, deleteTotpForUser, adminCreateLocalUser,
  adminInviteUser,
  createApiKey,
  createUserGroup,
  deleteUserGroup,
  modifyUserGroup,
  deleteGlobalVar,
  updateGlobalVar,
  createGlobalVar
} from 'src/services'
import SimpleUserTable from './simple_user_table'
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
import { NewApiKeyModalContents } from 'src/pages/account_settings/api_keys/modals'
import { BuildReloadBus } from 'src/helpers/reload_bus'

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
  const [apiKey, setApiKey] = React.useState<ApiKey | null>(null)
  const [newUserSlug, setNewUserSlug] = React.useState<string | null>(null)

  const headlessName = useFormField<string>("")
  const contactEmail = useFormField<string>("")
  const doCreateApiKey = useFormField(true)

  const handleSubmit = async () => {
    if (headlessName.value.length == 0) {
      throw new Error("Headless users must be given a name")
    }

    let createdSlug = newUserSlug
    if (createdSlug == null) {
      const newUser = await addHeadlessUser({
        firstName: 'Headless',
        lastName: headlessName.value,
        email: contactEmail.value,
      })
      setNewUserSlug(newUser.slug)
      createdSlug = newUser.slug
    }

    if (doCreateApiKey.value) {
      setApiKey(await createApiKey({
        userSlug: createdSlug,
      }))
    }
  }

  const formComponentProps = useForm({
    fields: [headlessName, contactEmail, doCreateApiKey],
    handleSubmit,
    onSuccess: () => {
      if (!doCreateApiKey.value) {
        props.onRequestClose()
      }
    }
  })
  return (
    <ModalForm
      title="Create New Headless User"
      submitText="Create"
      cancelText={apiKey == null ? "Cancel" : "Close"}
      onRequestClose={props.onRequestClose}
      {...formComponentProps}
      disableSubmit={apiKey != null}
    >
      <Input label="Headless name" {...headlessName} disabled={apiKey != null} />
      <Input type="email" label="Contact Email" {...contactEmail} disabled={apiKey != null} />
      <Checkbox label="Also create API key" {...doCreateApiKey} />
      {
        apiKey && <NewApiKeyModalContents apiKey={apiKey} />
      }
      {
        (apiKey == null && newUserSlug != null) && (
          <div>User created, but received an error creating key. Please try again.</div>
        )
      }
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
      const runSubmit = async () => {
        const result = await adminCreateLocalUser({
          firstName: firstName.value,
          lastName: lastName.value,
          email: contactEmail.value,
          username: contactEmail.value,
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
        <Input type="email" label="Email" {...contactEmail} disabled={isDisabled} />
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

export const AddUserGroupModal = (props: {
  onRequestClose: () => void,
}) => {
  const [isCompleted, setIsCompleted] = React.useState<boolean>(false)
  const [includedUsers, setIncludedUsers] = React.useState<Set<string>>(() => new Set());

  const name = useFormField<string>("")
  const userSlugs = Array.from(includedUsers as Set<string>)
  const formComponentProps = useForm({
    fields: [name],
    handleSubmit: () => {
      if (name.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("User group should have a name")))
      }
      const runSubmit = async () => {
        await createUserGroup({
          name: name.value,
          userSlugs: userSlugs
        })
        setIsCompleted(true)
      }
      return runSubmit()
    },
  })

  const bus = BuildReloadBus()
  return (
    <Modal title="Create New Group" onRequestClose={props.onRequestClose}>
      {isCompleted ? (<>
        <div className={cx('success-area')}>
          <p>Group has been created successfully!</p>
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      :
      (<>
      <h1 className={cx('header')}>Users</h1>
      <SimpleUserTable {...bus} setIncludedUsers={setIncludedUsers} includedUsers={includedUsers as Set<string>} />
      <Form {...formComponentProps} loading={isCompleted}
        submitText={isCompleted ? undefined : "Submit"}
      >
        <h1 className={cx('header')}>Name<span className={cx('optional')}>*</span></h1>
        <Input label="" {...name} disabled={isCompleted} />
      </Form>
      </>)
      }
    </Modal>
  )
}

export const ModifyUserGroupModal = (props: {
  userGroup: UserGroupAdminView,
  onRequestClose: () => void,
}) => {
  const [isCompleted, setIsCompleted] = React.useState<boolean>(false)
  const slugs = props.userGroup?.userSlugs ? props.userGroup.userSlugs : []
  const [includedUsers, setIncludedUsers] = React.useState(() => new Set([...slugs]));

  const name = useFormField<string>(props.userGroup.name)
  const formComponentProps = useForm({
    fields: [name],
    handleSubmit: () => {
      if (name.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("User goup should have a name")))
      }

      const slugsToAdd: Array<string> = []
      const slugsToRemove: Array<string> = []
      const initialSlugs = new Set([...slugs])
      const newSlugs = includedUsers as Set<string>
      initialSlugs.forEach((slug) => !newSlugs.has(slug) && slugsToRemove.push(slug))
      newSlugs.forEach((slug) => !initialSlugs.has(slug) && slugsToAdd.push(slug))

      const nameOrNull = name.value.toLowerCase() !== props.userGroup.name.toLowerCase() ? name.value : null
      const runSubmit = async () => {
        await modifyUserGroup({
          slug: props.userGroup.slug,
          newName: nameOrNull,
          userSlugsToAdd: slugsToAdd,
          userSlugsToRemove: slugsToRemove,
        })
        setIsCompleted(true)
      }
      return runSubmit()
    },
  })

  const bus = BuildReloadBus()
  return (
    <Modal title="Modify Group" onRequestClose={props.onRequestClose}>
      {isCompleted ? (<>
        <div className={cx('success-area')}>
          <p>Group has been modified successfully!</p>
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      :
      (<>
      <h1 className={cx('header')}>Users</h1>
      <SimpleUserTable {...bus} setIncludedUsers={setIncludedUsers} includedUsers={includedUsers as Set<string>} />
      <Form {...formComponentProps} loading={isCompleted}
        submitText={isCompleted ? undefined : "Submit"}
      >
        <h1 className={cx('header')}>Name<span className={cx('optional')}>*</span></h1>
        <Input label="" {...name} disabled={isCompleted} />
      </Form>
      </>)
      }
    </Modal>
  )
}

export const DeleteUserGroupModal = (props: {
  userGroup: UserGroupAdminView,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Delete User"
    warningText="This will remove the user group from the system. All user group information will be lost."
    submitText="Delete"
    challengeText={props.userGroup.slug}
    handleSubmit={() => deleteUserGroup({ userGroupSlug: props.userGroup.slug })}
    onRequestClose={props.onRequestClose}
  />

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

export const InviteUserModal = (props: {
  onRequestClose: () => void,
}) => {
  const firstName = useFormField<string>("")
  const lastName = useFormField<string>("")
  const contactEmail = useFormField<string>("")

  const [url, setUrl] = React.useState<string>("")
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
        const result = await adminInviteUser({
          firstName: firstName.value,
          lastName: lastName.value,
          email: contactEmail.value,
        })
        const url = `${window.location.origin}/web/auth/recovery/login?code=${result.code}`

        setUrl(url)
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
        <Input type="email" label="Email" {...contactEmail} disabled={isDisabled} />
      </Form>
      {isDisabled && (<>
        <div className={cx('success-area')}>
          <p>The user can login with the link below to configure their account:</p>
          <InputWithCopyButton label="Recovery Code" value={url} />
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      }
    </Modal>
  )
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

export const AddGlobalVarModal = (props: {
  onRequestClose: () => void,
}) => {
  const [isCompleted, setIsCompleted] = React.useState<boolean>(false)

  const name = useFormField<string>("")
  const value = useFormField<string>("")
  const formComponentProps = useForm({
    fields: [name, value],
    handleSubmit: () => {
      if (name.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("Global Variable should have a name")))
      }

      // TODO TN does this line up with what the API is expecting?
      const valOrNull = value.value === "" ? null : value.value
      const runSubmit = async () => {
        await createGlobalVar(name.value, valOrNull) 
        setIsCompleted(true)
      }
      return runSubmit()
    },
  })

  return (
    <Modal title="Add Variable" onRequestClose={props.onRequestClose}>
      {isCompleted ? (<>
        <div className={cx('success-area')}>
          <p>Variable has been added successfully!</p>
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      :
      (<>
      <Form {...formComponentProps} loading={isCompleted}
        submitText={isCompleted ? undefined : "Submit"}
      >
        <h1 className={cx('header')}>Name</h1>
        <Input label="" {...name} disabled={isCompleted} />
        <h1 className={cx('header')}>Value<span className={cx('optional')}>*</span></h1>
        <Input label="" {...value} disabled={isCompleted} />
      </Form>
      </>)
      }
    </Modal>
  )
}

export const DeleteGlobalVarModal = (props: {
  globalVar: GlobalVar,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Delete Global Variable"
    warningText="This will remove the global variable from the system."
    submitText="Delete"
    challengeText={props.globalVar.name}
    handleSubmit={() => deleteGlobalVar(props.globalVar.name)}
    onRequestClose={props.onRequestClose}
  />

export const ModifyGlobalVarModal = (props: {
  globalVar: GlobalVar,
  onRequestClose: () => void,
}) => {
  const [isCompleted, setIsCompleted] = React.useState<boolean>(false)

  const name = useFormField<string>(props.globalVar.name)
  const value = useFormField<string>(props.globalVar.value)
  const formComponentProps = useForm({
    fields: [name, value],
    handleSubmit: () => {
      if (name.value.length == 0) {
        return new Promise((_resolve, reject) => reject(Error("Global Variable should have a name")))
      }

      const nameOrNull = name.value.toLowerCase() !== props.globalVar.name.toLowerCase() ? name.value : null
      const valOrNull = value.value.toLowerCase() !== props.globalVar.value.toLowerCase() ? value.value : null
      const somethingChanged = nameOrNull !== null || valOrNull !== null
      const runSubmit = async () => {
        somethingChanged && await updateGlobalVar(props.globalVar.name, {
          value: valOrNull,
          newName: nameOrNull,
        }) 
        setIsCompleted(true)
      }
      return runSubmit()
    },
  })

  return (
    <Modal title="Modify Variable" onRequestClose={props.onRequestClose}>
      {isCompleted ? (<>
        <div className={cx('success-area')}>
          <p>Variable has been modified successfully!</p>
          <Button className={cx('success-close-button')} primary onClick={props.onRequestClose} >Close</Button>
        </div>
      </>)
      :
      (<>
      <Form {...formComponentProps} loading={isCompleted}
        submitText={isCompleted ? undefined : "Submit"}
      >
        <h1 className={cx('header')}>Name<span className={cx('optional')}>*</span></h1>
        <Input label="" {...name} disabled={isCompleted} />
        <h1 className={cx('header')}>Value<span className={cx('optional')}>*</span></h1>
        <Input label="" {...value} disabled={isCompleted} />
      </Form>
      </>)
      }
    </Modal>
  )
}
