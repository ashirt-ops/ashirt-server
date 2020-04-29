// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { Tag as TagType, Evidence } from 'src/global_types'
import { countBy } from 'lodash'
import { useDataSource, getTags, getEvidenceList } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'

import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import Tag from 'src/components/tag'
import { DeleteTagModal, EditTagModal } from './modals'
import { default as Button, ButtonGroup } from 'src/components/button'

const TagTable = (props: {
  operationSlug: string,
  tags: Array<TagType>,
  evidence: Array<Evidence>,
  onUpdate: () => void,
}) => {
  const numEvidenceByTagId: {[id: number]: number} = React.useMemo(() => (
    countBy(props.evidence.flatMap(evi => evi.tags.map(tag => tag.id)))
  ), [props.evidence])

  const editTagModal = useModal<{tag: TagType}>(modalProps => (
    <EditTagModal {...modalProps} operationSlug={props.operationSlug} onEdited={props.onUpdate} />
  ))
  const deleteTagModal = useModal<{tag: TagType}>(modalProps => (
    <DeleteTagModal {...modalProps} operationSlug={props.operationSlug} onDeleted={props.onUpdate} />
  ))

  return <>
    <Table columns={['Tag', '# Evidence Attached To', 'Actions']}>
      {props.tags.map(tag => (
        <tr key={tag.name}>
          <td><Tag name={tag.name} color={tag.colorName} /></td>
          <td>{numEvidenceByTagId[tag.id] || 0}</td>
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
  const ds = useDataSource()
  const wiredTags = useWiredData(React.useCallback(() => Promise.all([
    getTags(ds, { operationSlug: props.operationSlug }),
    getEvidenceList(ds, { operationSlug: props.operationSlug, query: '' }),
  ]), [ds, props.operationSlug]))

  return (
    <SettingsSection title="Operation Tags">
      {wiredTags.render(([tags, evidence]) => (
        <TagTable
          operationSlug={props.operationSlug}
          tags={tags}
          evidence={evidence}
          onUpdate={wiredTags.reload}
        />
      ))}
    </SettingsSection>
  )
}
