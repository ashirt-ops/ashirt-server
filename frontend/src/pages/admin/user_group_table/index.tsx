// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useNavigate } from 'react-router-dom'
import { PaginatedWiredData, usePaginatedWiredData } from 'src/helpers'

import { UserAdminView } from 'src/global_types'
import { listUsersAdminView, createRecoveryCode } from 'src/services'
import AuthContext from 'src/auth_context'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from 'src/helpers'

import {
  ResetPasswordModal, UpdateUserFlagsModal, DeleteUserModal, RecoverAccountModal,
  RemoveTotpModal
} from 'src/pages/admin_modals'
import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import Checkbox from 'src/components/checkbox'
import { StandardPager } from 'src/components/paging'
import SettingsSection from 'src/components/settings_section'
import { default as Menu, MenuItem } from 'src/components/menu'
import { ClickPopover } from 'src/components/popover'
import Input from 'src/components/input'

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
  const navigate = useNavigate()

  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const editUserFn = (u: UserAdminView) => navigate(`/account/profile?user=${u.slug}`)
  const recoverFn = (u: UserAdminView) => createRecoveryCode({ userSlug: u.slug }).then(setRecoveryCode)
  const columns = Object.keys(rowBuilder(null, <span />))

  const wiredUsers = usePaginatedWiredData<UserAdminView>(
    React.useCallback(page => listUsersAdminView({ page, pageSize: 10, deleted: withDeleted, name: usernameFilterValue }), [usernameFilterValue, withDeleted]),
    (err) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )
  const actionsBuilder = actionsForUserBuilder(self ? self.slug : "", wiredUsers) 


  React.useEffect(() => {
    props.onReload(wiredUsers.reload)
    return () => { props.offReload(wiredUsers.reload) }
  })
  React.useEffect(() => { setIncludeDeletedUsers(withDeleted) }, [withDeleted])

  return (
    <SettingsSection title="Group List" width="wide">
      <div className={cx('inline-form')}>
        <Input
          label="Group Filter"
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
  // TODO TN how to ensure the columns are closer together?
  <tr>
    <td>{props.data["Name"]}</td>
    <td>{props.data["Users"]}</td>
    {/* TODO TN where to add modify button? */}
  </tr>
)

type Rowdata = {
  "Name": string,
  "Users": JSX.Element,
}

const rowBuilder = (u: UserAdminView | null, actions: JSX.Element): Rowdata => ({
  "Name": u ? u.firstName : "",
  "Users": actions,
})

const actionsForUserBuilder = (selfSlug: string,
  wiredUsers: PaginatedWiredData<UserAdminView>
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

    const userCount = wiredUsers.render(data => <span>{data.reduce((a, c) => 1 + a, 0)}</span>)


    return (
      <ButtonGroup>
        <ClickPopover className={cx('popover')} closeOnContentClick content={
          <Menu>
            {/* TODO TN use that thing I made so I can load the data without rendering it */}
            {wiredUsers.render(data => <>
            {/* TODO TN should we allow a user to be removed from this interface? */}
          {data.map(user =>  <p className={cx('user')}>{user.slug}</p>)}
        </>)}
          </Menu>
        }>
          <Button small className={cx('arrow')}><p className={cx('button-text')}>{userCount} Users</p></Button>
        </ClickPopover>
        {/* TODO TN make count dynamic */}
        {/* <Button small disabled onClick={() => deleteUserFn(u)} {...canDelete}>8 Users</Button> */}
      </ButtonGroup>
    )
  }
