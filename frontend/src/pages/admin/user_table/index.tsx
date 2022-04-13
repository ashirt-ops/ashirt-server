// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { usePaginatedWiredData } from 'src/helpers'

import { UserAdminView } from 'src/global_types'
import { listUsersAdminView, createRecoveryCode } from 'src/services'
import AuthContext from 'src/auth_context'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from 'src/helpers'

import {
  ResetPasswordModal, UpdateUserFlagsModal, DeleteUserModal, RecoverAccountModal,
  RemoveTotpModal
} from 'src/pages/admin_modals'
import Table from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import Checkbox from 'src/components/checkbox'
import { StandardPager } from 'src/components/paging'
import ErrorDisplay from 'src/components/error_display'
import LoadingSpinner from 'src/components/loading_spinner'
import SettingsSection from 'src/components/settings_section'
import { default as Menu, MenuItem } from 'src/components/menu'
import { ClickPopover } from 'src/components/popover'
import Input from 'src/components/input'

import { useHistory } from 'react-router-dom'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const [resettingPassword, setResettingPassword] = React.useState<null | UserAdminView>(null)
  const [editingUserFlags, setEditingUserFlags] = React.useState<null | UserAdminView>(null)
  const [deletingUser, setDeletingUser] = React.useState<null | UserAdminView>(null)
  const [deletingTotp, setDeletingTotp] = React.useState<null | UserAdminView>(null)
  const [recoveryCode, setRecoveryCode] = React.useState<null | string>(null)
  const [withDeleted, setWithDeleted] = React.useState(getIncludeDeletedUsers())
  const self = React.useContext(AuthContext).user
  const history = useHistory()

  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const editUserFn = (u: UserAdminView) => history.push(`/account/edit/${u.slug}`)
  const recoverFn = (u: UserAdminView) => createRecoveryCode({ userSlug: u.slug }).then(setRecoveryCode)
  const actionsBuilder = actionsForUserBuilder(self ? self.slug : "", editUserFn, setResettingPassword, setEditingUserFlags, setDeletingUser, recoverFn, setDeletingTotp)
  const columns = Object.keys(rowBuilder(null, <span />))

  const asFullRow = (el: React.ReactElement): React.ReactElement => <tr><td colSpan={columns.length}>{el}</td></tr>

  const wiredUsers = usePaginatedWiredData<UserAdminView>(
    React.useCallback(page => listUsersAdminView({ page, pageSize: 10, deleted: withDeleted, name: usernameFilterValue }), [usernameFilterValue, withDeleted]),
    (err) => asFullRow(<ErrorDisplay err={err} />),
    () => asFullRow(<LoadingSpinner />)
  )

  React.useEffect(() => {
    props.onReload(wiredUsers.reload)
    return () => { props.offReload(wiredUsers.reload) }
  })
  React.useEffect(() => { setIncludeDeletedUsers(withDeleted) }, [withDeleted])

  return (
    <SettingsSection title="User List" width="wide">
      <div className={cx('inline-form')}>
        <Input
          label="User Filter"
          value={usernameFilterValue}
          onChange={v => { setUsernameFilterValue(v); wiredUsers.pagerProps.onPageChange(1) }}
          loading={usernameFilterValue.length > 0 && wiredUsers.loading}
        />
        <Checkbox
          label="Include Deleted Users"
          className={cx('checkbox')}
          value={withDeleted}
          onChange={setWithDeleted} />
      </div>
      <Table className={cx('table')} columns={columns}>
        {wiredUsers.render(data => <>
          {data.map(user => <TableRow key={user.slug} data={rowBuilder(user, actionsBuilder(user))} />)}
        </>)}
      </Table>
      <StandardPager className={cx('user-table-pager')} {...wiredUsers.pagerProps} />

      {resettingPassword && <ResetPasswordModal user={resettingPassword} onRequestClose={() => setResettingPassword(null)} />}
      {editingUserFlags && <UpdateUserFlagsModal user={editingUserFlags} onRequestClose={() => { setEditingUserFlags(null); wiredUsers.reload() }} />}
      {deletingUser && <DeleteUserModal user={deletingUser} onRequestClose={() => { setDeletingUser(null); wiredUsers.reload() }} />}
      {deletingTotp && <RemoveTotpModal user={deletingTotp} onRequestClose={() => { setDeletingTotp(null); wiredUsers.reload() }} />}
      {recoveryCode && <RecoverAccountModal recoveryCode={recoveryCode} onRequestClose={() => setRecoveryCode(null)} />}
    </SettingsSection>
  )
}

const TableRow = (props: { data: Rowdata }) => (
  <tr>
    <td>{props.data["First Name"]}</td>
    <td>{props.data["Last Name"]}</td>
    <td>{props.data["Contact Email"]}</td>
    <td>{props.data["Flags"]}</td>
    <td>{props.data["Actions"]}</td>
  </tr>
)

type Rowdata = {
  "First Name": string,
  "Last Name": string,
  "Contact Email": string,
  "Flags": JSX.Element,
  "Actions": JSX.Element,
}

const rowBuilder = (u: UserAdminView | null, actions: JSX.Element): Rowdata => ({
  "First Name": u ? u.firstName : "",
  "Last Name": u ? u.lastName : "",
  "Contact Email": u ? u.email : "",
  "Flags": u ? <UserFlags user={u} /> : <span />,
  "Actions": actions,
})

const UserFlags = (props: { user: UserAdminView }) => {
  return (
    <>
      {props.user.deleted
        ? <span className={cx('deleted-user')}>Deleted</span>
        : <span>{
          [{ label: "Headless", hasFlag: props.user.headless },
          { label: "Admin", hasFlag: props.user.admin },
          { label: "Disabled", hasFlag: props.user.disabled },
          ].filter(x => x.hasFlag).map(f => f.label).join(", ")
        }</span>
      }
    </>
  )
}


const actionsForUserBuilder = (selfSlug: string,
  editUserFn: (u: UserAdminView) => void,
  resetPwFn: (u: UserAdminView) => void,
  editFlagsFn: (u: UserAdminView) => void,
  deleteUserFn: (u: UserAdminView) => void,
  recoveryFn: (u: UserAdminView) => void,
  deleteTotpFn: (u: UserAdminView) => void,
) => (
  u: UserAdminView
) => {
    const deletedAttrs = { disabled: true, title: "User has been deleted" }
    const notDeletedOrSelf = (msg?: string) => {
      switch (true) {
        case u.deleted: return deletedAttrs
        case (u.slug === selfSlug): return { disabled: true, title: msg }
        default: return {}
      }
    }

    const canReset = () => {
      if (u.deleted) return deletedAttrs
      return {
        disabled: (u.authSchemes && !u.authSchemes.includes('local'))
      }
    }
    const canEditFlags = notDeletedOrSelf()
    const canDelete = notDeletedOrSelf("Admins cannot delete themselves")

    const canRecover = u.deleted ? deletedAttrs : {}
    const canRemoveTotp = () => {
      if (u.deleted) return deletedAttrs
      if (!u.hasLocalTotp) {
        return { disabled: true, title: "User does not have multi-factor Authentication enabled" }
      }
      return {}
    }

    return (
      <ButtonGroup>
        <ClickPopover className={cx('popover')} closeOnContentClick content={
          <Menu>
            <MenuItem onClick={() => editUserFn(u)}>Edit User</MenuItem>
            <MenuItem onClick={() => resetPwFn(u)} {...canReset()}>Reset Password</MenuItem>
            <MenuItem onClick={() => editFlagsFn(u)} {...canEditFlags}>Edit User Flags</MenuItem>
            <MenuItem onClick={() => deleteTotpFn(u)} {...canRemoveTotp()}>Remove Multi-Factor Authentication</MenuItem>
            <MenuItem onClick={() => recoveryFn(u)} {...canRecover}>Generate Recovery Code</MenuItem>
          </Menu>
        }>
          <Button small className={cx('arrow')} />
        </ClickPopover>
        <Button danger small onClick={() => deleteUserFn(u)} {...canDelete}>Delete User</Button>
      </ButtonGroup>
    )
  }
