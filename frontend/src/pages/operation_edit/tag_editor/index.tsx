// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import SettingsSection from 'src/components/settings_section'
import { default as Table, SortAsc, SortDesc, SortDirection, ColumnData } from 'src/components/table'
import Tag from 'src/components/tag'
import { DeleteTagModal, EditTagModal } from './modals'
import { TagWithUsage } from 'src/global_types'
import { default as Button, ButtonGroup } from 'src/components/button'
import { getTags } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'

// @ts-ignore - npm package @types/react-router-dom needs to be updated (https://github.com/DefinitelyTyped/DefinitelyTyped/issues/40131)
import { useHistory } from 'react-router-dom'

type compareableFunc = (l: unknown, r: unknown) => number

const sortNone: compareableFunc = (a: unknown, b: unknown) => 0
const sortNums: compareableFunc = (a: TagWithUsage, b: TagWithUsage) => a.evidenceCount - b.evidenceCount
const sortTags: compareableFunc = (a: TagWithUsage, b: TagWithUsage) => a.name.localeCompare(b.name)

const TagTable = (props: {
  operationSlug: string,
  tags: Array<TagWithUsage>,
  onUpdate: () => void,
}) => {
  const history = useHistory()
  const editTagModal = useModal<{ tag: TagWithUsage }>(modalProps => (
    <EditTagModal {...modalProps} operationSlug={props.operationSlug} onEdited={props.onUpdate} />
  ))
  const deleteTagModal = useModal<{ tag: TagWithUsage }>(modalProps => (
    <DeleteTagModal {...modalProps} operationSlug={props.operationSlug} onDeleted={props.onUpdate} />
  ))

  // React has a special meaning for passing functions into useState, so re-wrapping the function to
  // simply store the desired sorting function (we do the same with the setSortFunc value)
  const [sortFunc, setSortFunc] = React.useState<compareableFunc>(() => sortNone)

  const [columns, setColumns] = React.useState<Array<ColumnData & { compareVia: compareableFunc }>>([
    { title: '', label: 'Tag', clickable: true, compareVia: sortTags },
    { title: '', label: '# Evidence Attached To', clickable: true, compareVia: sortNums },
    { title: '', label: 'Actions', compareVia: sortNone },
  ])

  const updateColumnSorting = (index: number) => {
    const sortDirections: Array<{ compare: compareableFunc, dir: SortDirection }> = [
      { dir: SortAsc, compare: columns[index].compareVia },
      { dir: SortDesc, compare: (a, b) => columns[index].compareVia(b, a) },
      { dir: undefined, compare: sortNone }
    ]
    const matchIndex = (sortDirections.findIndex(x => x.dir == columns[index].sortDirection) + 1) % sortDirections.length
    setColumns(columns.map((col, idx) => ({ ...col, sortDirection: (idx == index) ? sortDirections[matchIndex].dir : undefined })))
    setSortFunc(() => sortDirections[matchIndex].compare)
  }

  return <>
    <Table columns={columns} onColumnClicked={updateColumnSorting}>
      {[...props.tags].sort(sortFunc).map(tag => (
        <tr key={tag.name}>
          <td><Tag name={tag.name} color={tag.colorName} onClick={() => history.push(`/operations/${props.operationSlug}/evidence?q=tag:"${tag.name}"`)} /></td>
          <td>{tag.evidenceCount}</td>
          <td>
            <ButtonGroup>
              <Button small onClick={() => editTagModal.show({ tag })}>Edit</Button>
              <Button small onClick={() => deleteTagModal.show({ tag })}>Delete</Button>
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
  const wiredTags = useWiredData(React.useCallback(() => getTags({ operationSlug: props.operationSlug }), [props.operationSlug]))

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
