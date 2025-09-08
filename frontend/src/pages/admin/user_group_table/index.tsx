import * as React from 'react'
import classnames from 'classnames/bind'
import { WiredData} from 'src/helpers'

import { UserGroupAdminView } from 'src/global_types'
import { listUserGroupsAdminView } from 'src/services'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from 'src/helpers'

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
import { DeleteUserGroupModal, ModifyUserGroupModal } from 'src/pages/admin_modals'
import { useWiredData } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const [deletingUserGroup, setDeletingUserGroup] = React.useState<null | UserGroupAdminView>(null)
  const [modifyingUserGroup, setModifyingUserGroup] = React.useState<null | UserGroupAdminView>(null)
  const [withDeleted, setWithDeleted] = React.useState(getIncludeDeletedUsers())
  const itemsPerPage = 10
  const [page, setPage] = React.useState(1)
  const [pageLength, setPageLength] = React.useState(0)

  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const columns = Object.keys(rowBuilder(null, <span />, <span />))

  const wiredUserGroups = useWiredData<UserGroupAdminView[]>(
    React.useCallback(() => listUserGroupsAdminView({  deleted: withDeleted }), [usernameFilterValue, withDeleted]),
    (err: Error) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )

  React.useEffect(() => {
    props.onReload(wiredUserGroups.reload)
    return () => { props.offReload(wiredUserGroups.reload) }
  })
  React.useEffect(() => { setIncludeDeletedUsers(withDeleted) }, [withDeleted])
  React.useEffect(() => {
    wiredUserGroups.expose(data => setPageLength(Math.ceil(data.length / itemsPerPage)))
  }, [wiredUserGroups])

  return (
    <SettingsSection title="Group List" width="wide">
      <div className={cx('inline-form')}>
        <Input
          label="Group Filter"
          value={usernameFilterValue}
          onChange={v => { setUsernameFilterValue(v); }}
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
          {data?.map((group, i) => {
            const belowUpperBound = i < page * itemsPerPage 
            const aboveLowerBound = i >= (page - 1) * itemsPerPage
            const inPageRange = belowUpperBound && aboveLowerBound
            return inPageRange && <TableRow key={group.slug} data={rowBuilder(group, usersInGroup(wiredUserGroups, group), modifyActions(group, setDeletingUserGroup, setModifyingUserGroup))} />
          })}
        </>)}
      </Table>
      <StandardPager className={cx('user-table-pager')} page={page} maxPages={pageLength} onPageChange={setPage} />

      {deletingUserGroup && <DeleteUserGroupModal userGroup={deletingUserGroup} onRequestClose={() => { setDeletingUserGroup(null); wiredUserGroups.reload() }} />}
      {modifyingUserGroup && <ModifyUserGroupModal userGroup={modifyingUserGroup} onRequestClose={() => { setModifyingUserGroup(null); wiredUserGroups.reload() }} />}
    </SettingsSection>
  )
}

const TableRow = (props: { data: Rowdata }) => (
  <tr>
    <td>{props.data["Name"]}</td>
    <td>{props.data["Users"]}</td>
    <td>{props.data["Flags"]}</td>
    <td>{props.data["Actions"]}</td>
  </tr>
)

type Rowdata = {
  "Name": string,
  "Users": React.JSX.Element,
  "Flags": React.JSX.Element,
  "Actions": React.JSX.Element,
}

const rowBuilder = (u: UserGroupAdminView | null, users: React.JSX.Element, actions: React.JSX.Element): Rowdata => ({
  "Name": u ? u.name : "",
  "Users": users,
  "Flags": (u && u.deleted) ? <span className={cx('deleted-user')}>Deleted</span> : <span />,
  "Actions": actions,
})

const usersInGroup = (
  wiredUserGroups: WiredData<UserGroupAdminView[]>,
  u: UserGroupAdminView
) => {
  const userCount = wiredUserGroups.render(data => <span>{data.find(group => group.slug === u.slug)?.userSlugs?.length ?? 0}</span>)
  return (
    <ButtonGroup>
      <ClickPopover className={cx('popover')} closeOnContentClick content={
        <Menu>
          {wiredUserGroups.render(data => {
            const group = data.find(group => u.slug === group.slug)
            const userList = group?.userSlugs?.map(userSlug => <p key={userSlug} className={cx('user')}>{userSlug}</p>)
            return <>{userList}</>
      })}
        </Menu>
      }>
        <Button small className={cx('arrow')}><p className={cx('button-text')}>{userCount} Users</p></Button>
      </ClickPopover>
    </ButtonGroup>
  )
}

const modifyActions = (
  u: UserGroupAdminView,
  onDeleteClick: (u: UserGroupAdminView) => void,
  onEditClick: (u: UserGroupAdminView) => void
) => {
  return (
    <ButtonGroup className={cx('row-buttons')}>
      <Button small disabled={u.deleted} onClick={() => onEditClick(u)}>Edit</Button>
      <Button small disabled={u.deleted} onClick={() => onDeleteClick(u)}>Delete</Button>
    </ButtonGroup>
  )
}
