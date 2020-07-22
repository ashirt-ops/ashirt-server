// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import Tag from 'src/components/tag'
import {DeleteTagModal, EditTagModal} from './modals'
import { TagWithUsage } from 'src/global_types'
import {default as Button, ButtonGroup} from 'src/components/button'
import {getTags} from 'src/services'
import {useWiredData, useModal, renderModals} from 'src/helpers'

const TagTable = (props: {
  operationSlug: string,
  tags: Array<TagWithUsage>,
  onUpdate: () => void,
}) => {
  const editTagModal = useModal<{ tag: TagWithUsage }>(modalProps => (
    <EditTagModal {...modalProps} operationSlug={props.operationSlug} onEdited={props.onUpdate} />
  ))
  const deleteTagModal = useModal<{ tag: TagWithUsage }>(modalProps => (
    <DeleteTagModal {...modalProps} operationSlug={props.operationSlug} onDeleted={props.onUpdate} />
  ))

  return <>
    <Table columns={['Tag', '# Evidence Attached To', 'Actions']}>
      {props.tags.map(tag => (
        <tr key={tag.name}>
          <td><Tag name={tag.name} color={tag.colorName} /></td>
          <td>{tag.evidenceCount}</td>
          <td>
            <ButtonGroup>
              <Button small onClick={() => editTagModal.show({tag})}>Edit</Button>
              <Button small onClick={() => deleteTagModal.show({tag})}>Delete</Button>
            </ButtonGroup>
          </td>
        </tr>
      ))}
    </Table>

    {renderModals(editTagModal, deleteTagModal)}
  </>
}

export default (props: {
  operationSlug: string,
}) => {
  const wiredTags = useWiredData(React.useCallback(() => getTags({operationSlug: props.operationSlug}), [props.operationSlug]))

  return (
    <SettingsSection title="Operation Tags">
      {wiredTags.render(tags => (
        <TagTable
          operationSlug={props.operationSlug}
          tags={tags}
          onUpdate={wiredTags.reload}
        />
      ))}
    </SettingsSection>
  )
}
