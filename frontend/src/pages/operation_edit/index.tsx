// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useParams, useNavigate, Routes, Route } from 'react-router-dom'

import Button from 'src/components/button'
import { NavVerticalTabMenu } from 'src/components/tab_vertical_menu'
import OperationEditor from './operation_editor'
import TagEditor from './tag_editor'
import UserPermissionEditor from './user_permission_editor'
import UserGroupPermissionEditor from './user_group_permission_editor'
import DeleteOperationButton from './delete_operation_button'
import BatchRunWorker from './batch_run_worker'

const cx = classnames.bind(require('./stylesheet'))

export const OperationEdit = () => {
  const { slug } = useParams<{ slug: string }>()
  const operationSlug = slug! // useParams puts everything in a partial, so our type above doesn't matter.
  const navigate = useNavigate()

  return (
    <>
      <Button
        className={cx('back-button')}
        icon={require('./back.svg')}
        onClick={() => navigate(-1)}>
        Back
      </Button>
      <NavVerticalTabMenu
        title="Edit Operation"
        tabs={[
          { id: "settings", label: "Settings" },
          { id: "users", label: "Users" },
          { id: "groups", label: "Groups"},
          { id: "tags", label: "Tags" },
          { id: "tasks", label: "Tasks" },
        ]} >
        <Routes>
          <Route path="settings" element={<SettingManagement operationSlug={operationSlug} />} />
          <Route path="users" element={<UserPermissionEditor operationSlug={operationSlug} />} />
          <Route path="groups" element={<UserGroupPermissionEditor operationSlug={operationSlug} />} />
          <Route path="tags" element={<TagEditor operationSlug={operationSlug} />} />
          <Route path="tasks" element={<BatchRunWorker operationSlug={operationSlug} />} />
        </Routes>
      </NavVerticalTabMenu>
    </>
  )
}
export default OperationEdit

const SettingManagement = (props: {
  operationSlug: string
}) => {
  return (<>
    <OperationEditor {...props} />
    <DeleteOperationButton {...props} />
  </>)
}
