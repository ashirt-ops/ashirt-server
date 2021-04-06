// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button, { ButtonGroup } from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { FindingCategory } from 'src/global_types'
import { getFindingCategories } from 'src/services'
import { useModal, useWiredData, renderModals } from 'src/helpers'

import {
  DeleteFindingCategoryModal,
  EditFindingCategoryModal,
} from './modals'

const columns = [
  'Name',
  'Actions',
]

const TableRow = (props: {
  category: FindingCategory,
  onUpdate: () => void
}) => {

  const editModal = useModal<void>(modalProps => (
    <EditFindingCategoryModal {...modalProps} onEdited={props.onUpdate} category={props.category} />
  ))
  const deleteModal = useModal<void>(modalProps => (
    <DeleteFindingCategoryModal {...modalProps} onDeleted={props.onUpdate} category={props.category} />
  ))

  return (
    <tr>
      <td>{props.category.category}</td>
      <td>
        <ButtonGroup>
          <Button small onClick={() => editModal.show()}>Edit</Button>
          <Button small danger onClick={() => deleteModal.show()}>Delete</Button>
        </ButtonGroup>
        {renderModals(editModal, deleteModal)}
      </td>
    </tr>
  )
}

export default (props: {
}) => {
  const wiredCategories = useWiredData<Array<FindingCategory>>(getFindingCategories)

  const createModal = useModal<void>(modalProps => (
    <EditFindingCategoryModal {...modalProps} onEdited={wiredCategories.reload} />
  ))

  return (
    <SettingsSection title="Finding Categories">
      {wiredCategories.render(data => (
        <>
          <Table columns={columns}>
            {data.map(category => (
              <TableRow key={category.id} category={category} onUpdate={wiredCategories.reload} />
            ))}
          </Table>
          <Button primary onClick={() => createModal.show()}>Add</Button>
        </>
      ))}
      {renderModals(createModal)}
    </SettingsSection>
  )
}
