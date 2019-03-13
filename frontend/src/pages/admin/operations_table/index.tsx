// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import OperationBadge from 'src/components/operation_badges'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { Operation } from 'src/global_types'
import { getOperationsForAdmin } from 'src/services'
import { useWiredData } from 'src/helpers'

// @ts-ignore - npm package @types/react-router-dom needs to be updated (https://github.com/DefinitelyTyped/DefinitelyTyped/issues/40131)
import { useHistory } from 'react-router-dom'

const columns = [
  'Slug',
  'Name',
  'Status',
  'Settings Link',
]

const TableRow = (props: {
  op: Operation,
}) => {
  const history = useHistory()

  return (
    <tr>
      <td>{props.op.slug}</td>
      <td>{props.op.name}</td>
      <td><OperationBadge numUsers={props.op.numUsers} status={props.op.status} /></td>
      <td><Button small onClick={() => history.push(`/operations/${props.op.slug}/edit/settings`)} >Settings</Button></td>
    </tr>
  )
}

export default (props: {
}) => {
  const wiredOps = useWiredData<Array<Operation>>(getOperationsForAdmin)

  return (
    <SettingsSection title="Operation List">
      {wiredOps.render(data =>
        <Table columns={columns}>
          {data.map(op => <TableRow key={op.slug} op={op} />)}
        </Table>
      )}
    </SettingsSection>
  )
}
