// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import AuthTable from './auth_table'
import HeadlessButton from './add_headless'
import { NavVerticalTabMenu } from 'src/components/tab_vertical_menu'
import CreateUserButton from "./add_user"
import InviteuserButton from "./invite_user"
import OperationsTable from './operations_table'
import FindingCategoriesTable from "./finding_categories_table"
import RecoveryMetrics from './recovery_metrics'
import UserTable from './user_table'

import { BuildReloadBus } from 'src/helpers/reload_bus'
import { DefaultTagEditor } from './default_tag_editor'
import { TagPorter } from './tag_porter'

const cx = classnames.bind(require('./stylesheet'))

export default () => {
  const bus = BuildReloadBus()

  return (
    <div className={cx('root')}>
      <NavVerticalTabMenu
        title="Admin Tools"
        tabs={[
          {
            id: "users", label: "User Management", content: <>
              <UserTable {...bus} />
              <HeadlessButton {...bus} />
              <CreateUserButton {...bus} />
              <InviteuserButton {...bus} />
            </>
          },
          {
            id: "authdata", label: "Authentication Overview", content: (
              <>
                <AuthTable />
                <RecoveryMetrics />
              </>
            )
          },
          {
            id: "operations", label: "Operation Management", content: (
              <>
                <OperationsTable />
              </>
            )
          },
          {
            id: "tags", label: "Tag Management", content: (
              <>
                <DefaultTagEditor {...bus}/>
                <TagPorter {...bus}/>
              </>
            )
          },
          {
            id: "findings", label: "Finding Categories", content: (
              <>
                <FindingCategoriesTable />
              </>
            )
          },
        ]}
      />
    </div>
  )
}
