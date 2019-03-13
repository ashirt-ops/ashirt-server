// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import Button from 'src/components/button'
import ErrorDisplay from 'src/components/error_display'
import LoadingSpinner from 'src/components/loading_spinner'
import Form from 'src/components/form'
import Modal from 'src/components/modal'
import Input from 'src/components/input'
import RadioGroup from 'src/components/radio_group'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import UserChooser from 'src/components/user_chooser'
import classnames from 'classnames/bind'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { User, UserRole, userRoleToLabel } from 'src/global_types'
import { getUserPermissions, setUserPermission } from 'src/services'
import { useForm, useFormField } from 'src/helpers/use_form'
import { useModal, renderModals, useWiredData } from 'src/helpers'
import { StandardPager } from 'src/components/paging'
const cx = classnames.bind(require('./stylesheet'))

const RoleSelect = (props: {
  disabled?: boolean,
  onChange: (r: UserRole) => void,
  label?: string,
  value: UserRole,
}) => (
    <RadioGroup
      disabled={props.disabled}
      groupLabel={props.label || ""}
      getLabel={(r: UserRole) => userRoleToLabel[r]}
      options={[UserRole.READ, UserRole.WRITE, UserRole.ADMIN]}
      value={props.value}
      onChange={props.onChange}
    />
  )

const NewUserForm = (props: {
  operationSlug: string,
  requestReload: () => void
}) => {
  const userField = useFormField<User | null>(null)
  const roleField = useFormField(UserRole.READ)
  const formProps = useForm({
    fields: [userField, roleField],
    handleSubmit: async () => {
      if (userField.value == null) throw Error("A user must be selected")
      await setUserPermission({
        operationSlug: props.operationSlug,
        userSlug: userField.value.slug,
        role: roleField.value,
      })
      userField.onChange(null)
      props.requestReload()
    }
  })
  return (
    <Form {...formProps}>
      <div className={cx('inline-form')}>
        <UserChooser {...userField} />
        <RoleSelect label="Role" {...roleField} />
        <Button primary loading={formProps.loading}>Add</Button>
      </div>
    </Form>
  )
}

const PermissionTableRow = (props: {
  disabled?: boolean,
  role: UserRole,
  user: User,
  requestReload: () => void
  updatePermissions: (role: UserRole) => Promise<void>
}) => {
  const currentUser = React.useContext(AuthContext).user
  const isCurrentUser = currentUser ? currentUser.slug === props.user.slug : false
  const isAdmin = currentUser ? currentUser.admin : false

  const removeWarningModal = useModal<void>(modalProps => (
    <RemoveWarningModal {...modalProps} removeUser={async () => {
      await props.updatePermissions(UserRole.NO_ACCESS)
      props.requestReload()
    }} />
  ))

  const disabled = props.disabled || (isCurrentUser && !isAdmin)

  return (
    <>
      <tr>
        <td style={{ fontWeight: isCurrentUser ? 800 : 400 }}>
          {props.user.firstName} {props.user.lastName}
        </td>
        <td><RoleSelect disabled={disabled} value={props.role} onChange={async (r) => {
          await props.updatePermissions(r)
          props.requestReload()
        }} /></td>
        <td><Button danger small disabled={disabled} onClick={() => removeWarningModal.show()}>Remove</Button></td>
      </tr>
      {renderModals(removeWarningModal)}
    </>
  )
}

const RemoveWarningModal = (props: {
  onRequestClose: () => void,
  removeUser: () => Promise<void>
}) => {
  const warningForm = useForm({
    fields: [],
    handleSubmit: async () => {
      props.removeUser()
    }
  })
  return (
    <Modal title="Remove User?" onRequestClose={props.onRequestClose}>
      <Form {...warningForm} submitText={"Remove User"} cancelText={"Go Back"} onCancel={props.onRequestClose}>
        <em>Removing this user will remove all read/write access to this user from this operation. Do you wish to continue?</em>
      </Form>
    </Modal>
  )
}

const PermissionTable = (props: {
  operationSlug: string,
  requestReload: () => void
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const columns = ['Name', 'Role', 'Remove']
  const itemsPerPage = 10

  const filterField = useFormField("")
  const [currentPage, setCurrentPage] = React.useState(1)

  const wiredPermissions = useWiredData(
    React.useCallback(() => getUserPermissions({ slug: props.operationSlug, name: "" }), [props.operationSlug]),
    (err) => <ErrorDisplay err={err} />,
    () => <LoadingSpinner />
  )

  React.useEffect(() => {
    props.onReload(wiredPermissions.reload)
    return () => { props.offReload(wiredPermissions.reload) }
  })

  return (
    <>
      {wiredPermissions.render(data => {
        const normalizeName = (user: User) => `${user.firstName} ${user.lastName}`.toLowerCase()
        const normalizedSearchTerm = filterField.value.toLowerCase()
        const matchingUsers = data.filter(({ user }) => normalizeName(user).includes(normalizedSearchTerm))
        const renderableData = matchingUsers.filter((_, i) => i >= ((currentPage - 1) * itemsPerPage) && i < (itemsPerPage * currentPage))

        return (
          <>
            <Input label="User Filter" {...filterField} />

            <Table columns={columns}>
              {renderableData.map(({ user, role }) => (
                <PermissionTableRow
                  key={user.slug}
                  requestReload={props.requestReload}
                  updatePermissions={(r: UserRole) => setUserPermission({ operationSlug: props.operationSlug, userSlug: user.slug, role: r })}
                  user={user}
                  role={role}
                />
              ))}
            </Table>
            <StandardPager
              className={cx('user-table-pager')}
              page={currentPage}
              maxPages={Math.ceil(matchingUsers.length / itemsPerPage)}
              onPageChange={(newPage) => setCurrentPage(newPage)}
            />
          </>
        )
      })}
    </>
  )
}

export default (props: {
  operationSlug: string,
}) => {
  const bus = BuildReloadBus()

  return (
    <SettingsSection title="Operation Users" width="wide">
      <NewUserForm
        {...bus}
        operationSlug={props.operationSlug}
      />
      <PermissionTable
        operationSlug={props.operationSlug}
        {...bus}
      />
    </SettingsSection>
  )
}
