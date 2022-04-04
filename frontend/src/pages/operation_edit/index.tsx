// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useParams, useNavigate } from 'react-router-dom'

import Button from 'src/components/button'
import { NavVerticalTabMenu } from 'src/components/tab_vertical_menu'
import OperationEditor from './operation_editor'
import TagEditor from './tag_editor'
import UserPermissionEditor from './user_permission_editor'
import DeleteOperationButton from './delete_operation_button'

const cx = classnames.bind(require('./stylesheet'))

export default () => {
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
          {
            id: "settings", label: "Settings", content: (<>
              <OperationEditor operationSlug={operationSlug} />
              <DeleteOperationButton operationSlug={operationSlug} />
            </>)
          },
          { id: "users", label: "Users", content: <UserPermissionEditor operationSlug={operationSlug} /> },
          { id: "tags", label: "Tags", content: <TagEditor operationSlug={operationSlug} /> },
        ]} />
    </>
  )
}
