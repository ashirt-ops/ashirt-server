// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { OperationWithExportData, ExportStatus } from 'src/global_types'
import { formatDistanceToNow } from 'date-fns'
import { useDataSource, getOperationsForAdmin, queueOperationExport } from 'src/services'
import { useWiredData } from 'src/helpers'

import Button, { ButtonGroup } from 'src/components/button'
import OperationBadge from 'src/components/operation_badges'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'

// @ts-ignore - npm package @types/react-router-dom needs to be updated (https://github.com/DefinitelyTyped/DefinitelyTyped/issues/40131)
import { useHistory } from 'react-router-dom'

const columns = [
  'Slug',
  'Name',
  'Status',
  'Last Archived',
  'Actions',
]

const TableRow = (props: {
  op: OperationWithExportData,
  reload: () => void
}) => {
  const ds = useDataSource()
  const history = useHistory()
  const [status, lastExport] = [props.op.exportStatus, props.op.lastCompletedExport]
  const slug = props.op.slug

  const archiveAttrs = {
    disabled: status == ExportStatus.PENDING || status == ExportStatus.IN_PROGRESS,
    title: (
      status == ExportStatus.PENDING ? "Archive has been queued" :
        status == ExportStatus.IN_PROGRESS ? "Operation is being archived now" : ""
    ),
  }

  const archiveClicked = async () => {
    await queueOperationExport(ds, slug)
    props.reload()
  }

  return (
    <tr>
      <td>{slug}</td>
      <td>{props.op.name}</td>
      <td><OperationBadge numUsers={props.op.numUsers} status={props.op.status} /></td>
      <td>{lastExport == null ? "Never" : formatDistanceToNow(lastExport, { addSuffix: true })}</td>
      <td>
        <ButtonGroup>
          <Button small onClick={() => history.push(`/operations/${slug}/edit/settings`)}>Settings</Button>
          <Button small {...archiveAttrs} onClick={archiveClicked}>Archive</Button>
        </ButtonGroup>
      </td>
    </tr>
  )
}

export default (props: {
}) => {
  const ds = useDataSource()
  const wiredOps = useWiredData<Array<OperationWithExportData>>(React.useCallback(() => (
    getOperationsForAdmin(ds)
  ), [ds]))

  return (
    <SettingsSection title="Operation List">
      {wiredOps.render(data =>
        <Table columns={columns}>
          {data.map(op => <TableRow key={op.slug} op={op} reload={wiredOps.reload} />)}
        </Table>
      )}
    </SettingsSection>
  )
}
