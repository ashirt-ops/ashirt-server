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
import UserGroupChooser from 'src/components/user_group_chooser'
import classnames from 'classnames/bind'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { UserGroup, UserOwnView, UserRole, userRoleToLabel } from 'src/global_types'
import { getUserGroupPermissions, getUserPermissions, setUserGroupPermission, setUserPermission } from 'src/services'
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
      if (userGroupField.value == null) throw Error("A user must be selected")
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
        <UserGroupChooser {...userGroupField} />
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
    <RemoveWarningModal {...modalProps} removeUser={async () => {
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
  removeUser: () => Promise<void>
}) => {
  const warningForm = useForm({
    fields: [],
    handleSubmit: async () => {
      props.removeUser()
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
  setIsOperationAdmin: (isOperationAdmin: boolean) => void,
  operationSlug: string,
  requestReload: () => void
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const columns = ['Name', 'Role', 'Remove']
  const itemsPerPage = 10

  const filterField = useFormField("")
  const [currentPage, setCurrentPage] = React.useState(1)
  const [isOperationAdmin, setLocalOperationAdmin] = React.useState(false)

  const normalizeName = (userGroup: UserGroup) => `${userGroup.name}`.toLowerCase()
  const normalizedSearchTerm = filterField.value.toLowerCase()

  const wiredPermissions = useWiredData(
    React.useCallback(() => getUserGroupPermissions({ slug: props.operationSlug, name: "" }), [props.operationSlug]),
    (err) => <ErrorDisplay err={err} />,
    () => <LoadingSpinner />
  )

  React.useEffect(() => {
    props.onReload(wiredPermissions.reload)
    wiredPermissions.expose(data => {
      const matchingUsers = data.filter(({ userGroup }) => normalizeName(userGroup).includes(normalizedSearchTerm))
      const renderableData = matchingUsers.filter((_, i) => i >= ((currentPage - 1) * itemsPerPage) && i < (itemsPerPage * currentPage))

      setLocalOperationAdmin(renderableData.find(datum => datum.userGroup.slug === props.currentUser?.slug)?.role === UserRole.ADMIN)
      props.setIsOperationAdmin(isOperationAdmin)
    })
    return () => { props.offReload(wiredPermissions.reload) }
  })

  // TODO TN add something so there's not an error message when there are not user groups
  return (
    <>
      {wiredPermissions.render(data => {
        const matchingUsers = data.filter(({ userGroup }) => normalizeName(userGroup).includes(normalizedSearchTerm))
        const renderableData = matchingUsers.filter((_, i) => i >= ((currentPage - 1) * itemsPerPage) && i < (itemsPerPage * currentPage))

        const notAdmin = !props.isAdmin && !isOperationAdmin

        return (
          <>
            <Input label="User Group Filter" {...filterField} />

            <Table columns={columns}>
              {renderableData.map(({ userGroup, role }) => (
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
  // isAdmin: boolean,
}) => {
  const bus = BuildReloadBus()

  // if (!props.isAdmin) {
  //   return <Navigate to="/operations" replace />;
  // }

  // TODO TN - ask if non sys admins should even be able to see this?

  const [isOperationAdmin, setIsOperationAdmin] = React.useState(false)
  const currentUser = React.useContext(AuthContext)?.user
  const isSysAdmin = currentUser ? currentUser?.admin : false
  const isAdmin = isSysAdmin || isOperationAdmin
  // const isAdmin = props.isAdmin || isOperationAdmin

  return (
    <SettingsSection title="Operation User Groups" width="wide">
      {isAdmin && (<NewUserGroupForm
        {...bus}
        operationSlug={props.operationSlug}
      />)}
      <PermissionTable
        currentUser={currentUser || undefined}
        isAdmin={isAdmin}
        setIsOperationAdmin={setIsOperationAdmin}
        operationSlug={props.operationSlug}
        {...bus}
      />
    </SettingsSection>
  )
}
