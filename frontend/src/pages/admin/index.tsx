// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { RouteComponentProps } from 'react-router-dom'

import AuthTable from './auth_table'
import HeadlessButton from './add_headless'
import CreateUserButton from "./add_user"
import NavVerticalTab from 'src/components/tab_vertical_menu'
import OperationsTable from './operations_table'
import FindingCategoriesTable from "./finding_categories_table"
import RecoveryMetrics from './recovery_metrics'
import UserTable from './user_table'

import { BuildReloadBus } from 'src/helpers/reload_bus'

const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps) => {
  const bus = BuildReloadBus()

  return (
    <div className={cx('root')}>
      <NavVerticalTab {...props}
        title="Admin Tools"
        tabs={[
          {
            id: "users", label: "User Management", content: <>
              <UserTable {...bus} />
              <HeadlessButton {...bus} />
              <CreateUserButton {...bus} />
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
                <FindingCategoriesTable />
              </>
            )
          },
        ]}
      />
    </div>
  )
}
