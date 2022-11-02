// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useNavigate } from 'react-router-dom'

import Button from 'src/components/button'
import OperationBadge from 'src/components/operation_badges'
import OperationBadgesModal from 'src/components/operation_badges_modal'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { Operation } from 'src/global_types'
import { getOperationsForAdmin } from 'src/services'
import { renderModals, useModal, useWiredData } from 'src/helpers'

const columns = [
  'Slug',
  'Name',
  'Status',
  'Settings Link',
]

const TableRow = (props: {
  op: Operation,
}) => {
  const navigate = useNavigate()
  const { op } = props

  const moreDetailsModal = useModal<{}>(modalProps => (
    <OperationBadgesModal {...modalProps} topContribs={op.topContribs} evidenceCount={op.evidenceCount} numTags={op.numTags} />
  ))

  const handleDetailsModal = () => moreDetailsModal?.show({})

  return (
    <tr>
      <td>{op.slug}</td>
      <td>{op.name}</td>
      <td><OperationBadge numUsers={op.numUsers} showDetailsModal={handleDetailsModal} /></td>
      <td><Button small onClick={() => navigate(`/operations/${op.slug}/edit/settings`)} >Settings</Button></td>
      {renderModals(moreDetailsModal)}
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
