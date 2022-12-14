// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { usePaginatedWiredData } from 'src/helpers'

import { UserAdminView } from 'src/global_types'
import { listUsersAdminView } from 'src/services'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from 'src/helpers'

import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import ComplexCheckbox from 'src/components/checkbox_complex'
import Checkbox from 'src/components/checkbox'
import { StandardPager } from 'src/components/paging'
import SettingsSection from 'src/components/settings_section'
import Input from 'src/components/input'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
  setIncludedUsers: (users: Set<string>) => void
  includedUsers: Set<string>
}) => {
  const [withDeleted, setWithDeleted] = React.useState(getIncludeDeletedUsers())
  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const toggleItem = (e: React.ChangeEvent<HTMLInputElement>, userSlug: string): void => {
    const isUserIncluded = e.target.checked
    if (isUserIncluded) {
      const newSet = new Set(props.includedUsers).add(userSlug)
      props.setIncludedUsers(newSet);
    } else {
      const newSet = new Set(props.includedUsers);
      newSet.delete(userSlug);
      props.setIncludedUsers(newSet);
    }
  }

  const columns = Object.keys({})

  const wiredUsers = usePaginatedWiredData<UserAdminView>(
    React.useCallback(page => listUsersAdminView({ page, pageSize: 5, deleted: withDeleted, name: usernameFilterValue }), [usernameFilterValue, withDeleted]),
    (err) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )

  React.useEffect(() => {
    props.onReload(wiredUsers.reload)
    return () => { props.offReload(wiredUsers.reload) }
  })
  React.useEffect(() => { setIncludeDeletedUsers(withDeleted) }, [withDeleted])

  return (
    <SettingsSection title="" width="wide">
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
          {data.map(user =>
            (<tr>
              <td>{`${user.firstName} ${user.lastName}`}</td>
              <td>
                <ComplexCheckbox
                  className={cx('checkbox')}
                  value={props.includedUsers.has(user.slug)}
                  onChange={(e) => toggleItem(e, user.slug)}
                  />
              </td>
            </tr>)
          )}
        </>)}
      </Table>
      <StandardPager className={cx('user-table-pager')} {...wiredUsers.pagerProps} />
    </SettingsSection>
  )
}
