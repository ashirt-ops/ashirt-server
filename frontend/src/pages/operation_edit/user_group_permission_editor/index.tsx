// Copyright 2023, Yahoo Inc.
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
import UserGroupChooser from 'src/components/user_group_chooser'
import classnames from 'classnames/bind'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { UserGroup, UserOwnView, UserRole, userRoleToLabel } from 'src/global_types'
import { getUserGroupPermissions, setUserGroupPermission } from 'src/services'
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

const NewUserGroupForm = (props: {
  operationSlug: string,
  requestReload: () => void
}) => {
  const userGroupField = useFormField<UserGroup | null>(null)
  const roleField = useFormField(UserRole.READ)
  const formProps = useForm({
    fields: [userGroupField, roleField],
    handleSubmit: async () => {
      if (userGroupField.value == null) throw Error("A user group must be selected")
      await setUserGroupPermission({
        operationSlug: props.operationSlug,
        userGroupSlug: userGroupField.value.slug,
        role: roleField.value,
      })
      userGroupField.onChange(null)
      props.requestReload()
    }
  })
  return (
    <Form {...formProps}>
      <div className={cx('inline-form')}>
        <UserGroupChooser operationSlug={props.operationSlug} {...userGroupField} />
        <RoleSelect label="Role" {...roleField} />
        <Button primary loading={formProps.loading}>Add</Button>
      </div>
    </Form>
  )
}

const PermissionTableRow = (props: {
  disabled?: boolean,
  role: UserRole,
  userGroup: UserGroup,
  currentUser?: UserOwnView,
  requestReload: () => void
  updatePermissions: (role: UserRole) => Promise<void>
}) => {
  const currentUserGroup = props?.currentUser
  const isCurrentUser = currentUserGroup ? currentUserGroup.slug === props.userGroup.slug : false

  const removeWarningModal = useModal<{}>(modalProps => (
    <RemoveWarningModal {...modalProps} removeUserGroup={async () => {
      await props.updatePermissions(UserRole.NO_ACCESS)
      props.requestReload()
    }} />
  ))

  const disabled = props.disabled

  return (
    <>
      <tr>
        <td style={{ fontWeight: isCurrentUser ? 800 : 400 }}>
          {props.userGroup.name}
        </td>
        <td><RoleSelect disabled={disabled} value={props.role} onChange={async (r) => {
          await props.updatePermissions(r)
          props.requestReload()
        }} /></td>
        <td><Button danger small disabled={disabled} onClick={() => removeWarningModal.show({})}>Remove</Button></td>
      </tr>
      {renderModals(removeWarningModal)}
    </>
  )
}

const RemoveWarningModal = (props: {
  onRequestClose: () => void,
  removeUserGroup: () => Promise<void>
}) => {
  const warningForm = useForm({
    fields: [],
    handleSubmit: async () => {
      props.removeUserGroup()
    }
  })
  return (
    <Modal title="Remove User Group?" onRequestClose={props.onRequestClose}>
      <Form {...warningForm} submitText={"Remove User Group"} cancelText={"Go Back"} onCancel={props.onRequestClose}>
        <em>Removing this user group will remove all read/write access to all the users in that group from this operation. Do you wish to continue?</em>
      </Form>
    </Modal>
  )
}

const PermissionTable = (props: {
  currentUser?: UserOwnView,
  isAdmin: boolean,
  operationSlug: string,
  requestReload: () => void
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const columns = ['Name', 'Role', 'Remove']
  const itemsPerPage = 10

  const filterField = useFormField("")
  const [page, setPage] = React.useState(1)

  const normalizeName = (userGroup: UserGroup) => `${userGroup.name}`.toLowerCase()
  const normalizedSearchTerm = filterField.value.toLowerCase()

  const wiredPermissions = useWiredData(
    React.useCallback(() => getUserGroupPermissions({ slug: props.operationSlug, name: "" }), [props.operationSlug]),
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
        const matchingUsers = data.filter(({ userGroup }) => normalizeName(userGroup).includes(normalizedSearchTerm))
        const usersInPageRange = matchingUsers.filter((_, i) => {
          const belowUpperBound = i < page * itemsPerPage 
          const aboveLowerBound = i >= (page - 1) * itemsPerPage
          return belowUpperBound && aboveLowerBound
        })

        const notAdmin = !props.isAdmin

        return (
          <>
            <Input label="User Group Filter" {...filterField} />

            <Table columns={columns}>
              {usersInPageRange.map(({ userGroup, role }) => (
                <PermissionTableRow
                  currentUser={props?.currentUser}
                  disabled={notAdmin}
                  key={userGroup.slug}
                  requestReload={props.requestReload}
                  updatePermissions={(r: UserRole) => setUserGroupPermission({ operationSlug: props.operationSlug, userGroupSlug: userGroup.slug, role: r })}
                  userGroup={userGroup}
                  role={role}
                />
              ))}
            </Table>
            <StandardPager
              className={cx('user-table-pager')}
              page={page}
              maxPages={Math.ceil(matchingUsers.length / itemsPerPage)}
              onPageChange={(newPage) => setPage(newPage)}
            />
          </>
        )
      })}
    </>
  )
}

export default (props: {
  operationSlug: string,
  isAdmin: boolean,
}) => {
  const bus = BuildReloadBus()
  const currentUser = React.useContext(AuthContext)?.user

  return (
    <SettingsSection title="Operation User Groups" width="wide">
      {props.isAdmin && (<NewUserGroupForm
        {...bus}
        operationSlug={props.operationSlug}
      />)}
      <PermissionTable
        currentUser={currentUser || undefined}
        isAdmin={props.isAdmin}
        operationSlug={props.operationSlug}
        {...bus}
      />
    </SettingsSection>
  )
}
