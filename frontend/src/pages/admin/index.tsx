// Copyright 2022, Yahoo Inc.
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
import ServiceWorkerTable from './service_worker_table'
import AddServiceWorker from './service_worker_table/add_service_button'

import { BuildReloadBus, BusSupportedService } from 'src/helpers/reload_bus'
import { DefaultTagEditor } from './default_tag_editor'
import { TagPorter } from './tag_porter'
import { Route, Routes } from 'react-router-dom'

const cx = classnames.bind(require('./stylesheet'))

export const AdminTools = () => {
  const bus = BuildReloadBus()

  return (
    <>
      <div className={cx('root')}>
        <NavVerticalTabMenu
          title="Admin Tools"
          tabs={[
            { id: "users", label: "User Management" },
            { id: "authdata", label: "Authentication Overview" },
            { id: "operations", label: "Operation Management" },
            { id: "tags", label: "Tag Management" },
            { id: "findings", label: "Finding Categories" },
            { id: "services", label: "Service Workers" },
          ]}
        >
          <Routes>
            <Route path="users" element={<UserManagement {...bus} />} />
            <Route path="authdata" element={<AuthOverview />} />
            <Route path="operations" element={<OperationsTable />} />
            <Route path="tags" element={<TagManagement {...bus} />} />
            <Route path="findings" element={<FindingCategoriesTable />} />
            <Route path="services" element={<ServiceWorkers {...bus}/>} />
          </Routes>
        </NavVerticalTabMenu>
      </div>
    </>
  )
}

export default AdminTools

const UserManagement = (props: BusSupportedService) => (
  <>
    <UserTable {...props} />
    <HeadlessButton {...props} />
    <CreateUserButton {...props} />
    <InviteuserButton {...props} />
  </>
)

const TagManagement = (props: BusSupportedService) => (
  <>
    <DefaultTagEditor {...props} />
    <TagPorter {...props} />
  </>
)

const AuthOverview = () => (
  <>
    <AuthTable />
    <RecoveryMetrics />
  </>
)

const ServiceWorkers = (props: BusSupportedService) => (
  <>
    <ServiceWorkerTable {...props} />
    <AddServiceWorker {...props} />
  </>
)
