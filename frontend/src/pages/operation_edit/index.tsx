// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { RouteComponentProps } from 'react-router-dom'

import Button from 'src/components/button'
import NavVerticalTab from 'src/components/tab_vertical_menu'
import OperationEditor from './operation_editor'
import TagEditor from './tag_editor'
import UserPermissionEditor from './user_permission_editor'

const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{ slug: string }>) => {
  const { slug } = props.match.params

  return (
    <>
      <Button className={cx('back-button')} icon={require('./back.svg')} onClick={() => props.history.goBack()}>Back</Button>
      <NavVerticalTab {...props}
        title="Edit Operation"
        tabs={[
          { id: "settings", label: "Settings", content: <OperationEditor operationSlug={slug} /> },
          { id: "users", label: "Users", content: <UserPermissionEditor operationSlug={slug} /> },
          { id: "tags", label: "Tags", content: <TagEditor operationSlug={slug} /> },
        ]} />
    </>
  )
}
