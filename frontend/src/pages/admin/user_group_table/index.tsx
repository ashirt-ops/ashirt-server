// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { PaginatedWiredData, usePaginatedWiredData} from 'src/helpers'

import { UserGroupAdminView } from 'src/global_types'
import { listUserGroupsAdminView } from 'src/services'
import AuthContext from 'src/auth_context'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from 'src/helpers'

import { RecoverAccountModal } from 'src/pages/admin_modals'
import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import Checkbox from 'src/components/checkbox'
import { StandardPager } from 'src/components/paging'
import SettingsSection from 'src/components/settings_section'
import { default as Menu } from 'src/components/menu'
import { ClickPopover } from 'src/components/popover'
import Input from 'src/components/input'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const [recoveryCode, setRecoveryCode] = React.useState<null | string>(null)
  const [withDeleted, setWithDeleted] = React.useState(getIncludeDeletedUsers())
  const self = React.useContext(AuthContext).user

  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const columns = Object.keys(rowBuilder(null, <span />))

  const wiredUserGroups = usePaginatedWiredData<UserGroupAdminView>(
    React.useCallback(page => listUserGroupsAdminView({ page, pageSize: 10, deleted: withDeleted }), [usernameFilterValue, withDeleted]),
    (err) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )
  const actionsBuilder = actionsForUserBuilder(self ? self.slug : "", wiredUserGroups)

  React.useEffect(() => {
    props.onReload(wiredUserGroups.reload)
    return () => { props.offReload(wiredUserGroups.reload) }
  })
  React.useEffect(() => { setIncludeDeletedUsers(withDeleted) }, [withDeleted])

  return (
    <SettingsSection title="Group List" width="wide">
      <div className={cx('inline-form')}>
        <Input
          label="Group Filter"
          value={usernameFilterValue}
          onChange={v => { setUsernameFilterValue(v); wiredUserGroups.pagerProps.onPageChange(1) }}
          loading={usernameFilterValue.length > 0 && wiredUserGroups.loading}
        />
        <Checkbox
          label="Include Deleted Groups"
          className={cx('checkbox')}
          value={withDeleted}
          onChange={setWithDeleted} />
      </div>
      <Table className={cx('table')} columns={columns}>
        {wiredUserGroups.render(data => <>
          {data.map(group => <TableRow key={group.slug} data={rowBuilder(group, actionsBuilder(group))} />)}
        </>)}
      </Table>
      <StandardPager className={cx('user-table-pager')} {...wiredUserGroups.pagerProps} />

      {/* {resettingPassword && <ResetPasswordModal user={resettingPassword} onRequestClose={() => setResettingPassword(null)} />}
      {editingUserFlags && <UpdateUserFlagsModal user={editingUserFlags} onRequestClose={() => { setEditingUserFlags(null); wiredUserGroups.reload() }} />}
      {deletingUser && <DeleteUserModal user={deletingUser} onRequestClose={() => { setDeletingUser(null); wiredUserGroups.reload() }} />}
      {deletingTotp && <RemoveTotpModal user={deletingTotp} onRequestClose={() => { setDeletingTotp(null); wiredUserGroups.reload() }} />} */}
      {recoveryCode && <RecoverAccountModal recoveryCode={recoveryCode} onRequestClose={() => setRecoveryCode(null)} />}
    </SettingsSection>
  )
}

const TableRow = (props: { data: Rowdata }) => (
  <tr>
    <td>{props.data["Name"]}</td>
    <td>{props.data["Users"]}</td>
    {/* TODO TN add modify button in next PR */}
  </tr>
)

type Rowdata = {
  "Name": string,
  "Users": JSX.Element,
}

const rowBuilder = (u: UserGroupAdminView | null, actions: JSX.Element): Rowdata => ({
  "Name": u ? u.slug : "",
  "Users": actions,
})

const actionsForUserBuilder = (selfSlug: string,
  wiredUserGroups: PaginatedWiredData<UserGroupAdminView>,
) => (
  u: UserGroupAdminView
) => {
  const userCount = wiredUserGroups.render(data => <span>{data.find(group => group.slug === u.slug)?.userSlugs?.length ?? 0}</span>)
    return (
      <ButtonGroup>
        <ClickPopover className={cx('popover')} closeOnContentClick content={
          <Menu>
            {wiredUserGroups.render(data => {
              const group = data.find(group => u.slug === group.slug)
              const userList = group?.userSlugs?.map(userSlug => <p className={cx('user')}>{userSlug}</p>)
              return <>{userList}</>
        })}
          </Menu>
        }>
          <Button small className={cx('arrow')}><p className={cx('button-text')}>{userCount} Users</p></Button>
        </ClickPopover>
      </ButtonGroup>
    )
  }
