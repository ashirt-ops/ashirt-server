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

  const [tagTableState, dispatch] = React.useReducer(tagTableReducer, TagTableInitialState)

  const baseColumns = [
    { title: '', label: 'Tag', clickable: true, compareVia: sortTags },
    { title: '', label: '# Evidence Attached To', clickable: true, compareVia: sortNums },
    { title: '', label: 'Actions', compareVia: sortNone },
  ]

  const updateColumnSorting = (index: number) => {
    const sortDirections: Array<{ compare: compareableFunc, dir: SortDirection }> = [
      { dir: SortAsc, compare: baseColumns[index].compareVia },
      { dir: SortDesc, compare: (a, b) => baseColumns[index].compareVia(b, a) },
      { dir: undefined, compare: sortNone }
    ]
    const matchIndex = index != tagTableState.sortColIndex
      ? 0
      : (sortDirections.findIndex(x => x.dir == tagTableState.sortDir) + 1) % sortDirections.length

    const sortDirIndex = sortDirections[matchIndex]

    dispatch({
      type: 'sort-column',
      sortFunc: sortDirIndex.compare,
      sortColIndex: index,
      sortDir: sortDirIndex.dir
    })
  }

  const sortedTags = [...props.tags].sort(tagTableState.sortFunc)
  return <>
    <Table columns={baseColumns.map((col, idx) => ({
      ...col,
      sortDirection: (idx == tagTableState.sortColIndex ? tagTableState.sortDir : undefined)
    }))} onColumnClicked={updateColumnSorting}>
      {sortedTags.map(tag => (
        <tr key={tag.name}>
          <td>
            <Tag
              name={tag.name}
              color={tag.colorName}
              onClick={() => history.push(`/operations/${props.operationSlug}/evidence?q=tag:"${tag.name}"`)}
            />
          </td>
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

const tagTableReducer = (state: TagTableState, action: TagTableAction): TagTableState => {
  if (action.type == 'sort-column') {
    return {
      ...state,
      ...action
    }
  }
  return state
}

type TagTableState = {
  sortFunc: compareableFunc
  sortDir: SortDirection
  sortColIndex: number
}

const TagTableInitialState = {
  sortFunc: sortNone,
  sortDir: undefined,
  sortColIndex: 0
}

type TagTableSortColumn = {
  type: 'sort-column'
  sortFunc: compareableFunc
  sortColIndex: number
  sortDir: SortDirection
}

type TagTableAction =
  | TagTableSortColumn
