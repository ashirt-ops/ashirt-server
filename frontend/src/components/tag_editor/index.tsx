// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { chunk } from 'lodash'
import { useNavigate } from 'react-router-dom'
import { DefaultTag, Tag as TagType, TagWithUsage } from 'src/global_types'
import { useModal, renderModals } from 'src/helpers'

import { StandardPager } from 'src/components/paging'
import { default as Table, EndAlignedColumn, SortAsc, SortDesc, SortDirection } from 'src/components/table'
import Input from 'src/components/input'
import { default as Button, ButtonGroup } from 'src/components/button'
import Tag from 'src/components/tag'
import { DeleteDefaultTagModal, DeleteOperationTagModal, UpsertOperationTagModal, UpsertDefaultTagModal } from './modals'


const cx = classnames.bind(require('./stylesheet'))

export const OperationTagTable = (props: {
  operationSlug: string,
  tags: Array<TagWithUsage>,
  onUpdate: () => void,
}) => {
  const navigate = useNavigate()
  const editTagModal = useModal<{ tag?: TagWithUsage }>(modalProps => (
    <UpsertOperationTagModal {...modalProps} operationSlug={props.operationSlug} onEdited={props.onUpdate} />
  ))
  const deleteTagModal = useModal<{ tag: TagWithUsage }>(modalProps => (
    <DeleteOperationTagModal {...modalProps} operationSlug={props.operationSlug} onDeleted={props.onUpdate} />
  ))
  const extraColumns: Array<tagTableColumn<TagWithUsage>> = [
    {
      title: '', label: '# Evidence Attached To', clickable: true, compareVia: sortNums,
      renderer: (tag: TagWithUsage) => tag.evidenceCount
    },
  ]

  return (
    <>
      <BasicTagTable
        tags={props.tags}
        extraColumns={extraColumns}
        onUpdate={props.onUpdate}
        onDeleteClick={tag => deleteTagModal.show({ tag })}
        onEditClick={tag => editTagModal.show({ tag })}
        onTagClick={tag => navigate(`/operations/${props.operationSlug}/evidence?q=tag:"${tag.name}"`)}
        onCreateClick={() => editTagModal.show({})}
      />
      {renderModals(editTagModal, deleteTagModal)}
    </>
  )
}

export const DefaultTagTable = (props: {
  tags: Array<DefaultTag>,
  onUpdate: () => void,
}) => {
  const editTagModal = useModal<{ tag?: TagType }>(modalProps => (
    <UpsertDefaultTagModal {...modalProps} onEdited={props.onUpdate} />
  ))
  const deleteTagModal = useModal<{ tag: TagType }>(modalProps => (
    <DeleteDefaultTagModal {...modalProps} onDeleted={props.onUpdate} />
  ))

  return (
    <>
      <BasicTagTable
        tags={props.tags}
        onUpdate={props.onUpdate}
        onDeleteClick={tag => deleteTagModal.show({ tag })}
        onEditClick={tag => editTagModal.show({ tag })}
        onCreateClick={() => editTagModal.show({})}
      />
      {renderModals(editTagModal, deleteTagModal)}
    </>
  )
}

function BasicTagTable<T extends TagType>(props: {
  onDeleteClick: (tag: T) => void
  onEditClick: (tag: T) => void
  onCreateClick?: () => void
  extraColumns?: Array<tagTableColumn<T>>
  tags: Array<T>,
  onTagClick?: (tag: T) => void,
  onUpdate: () => void,
  operationSlug?: string,
}) {
  const extraCols = props.extraColumns ?? []
  const [tagTableState, dispatch] = React.useReducer(tagTableReducer, TagTableInitialState)

  const columnRenders = extraCols.map(col => col.renderer)
  const columnDefinitions = extraCols.map(col => {
    const { renderer, ...rest } = col
    return rest
  })

  const baseColumns = [
    { title: '', label: 'Tag', clickable: true, compareVia: sortTags },
    ...columnDefinitions,
    { title: '', label: 'Actions', compareVia: sortNone, style: EndAlignedColumn },
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

  const sortedTags = props.tags
    .filter(tag => tag.name.toLowerCase().includes(tagTableState.filterText))
    .sort(tagTableState.sortFunc)
  const paginatedTags = chunk(sortedTags, 10)

  return <>
    <Input
      placeholder={"Filter Tags..."}
      value={tagTableState.filterText}
      onChange={(val) => dispatch({ type: 'filter-text-change', filterText: val })}
    />
    <Table className={cx('table')} columns={baseColumns.map((col, idx) => ({
      ...col,
      sortDirection: (idx == tagTableState.sortColIndex ? tagTableState.sortDir : undefined)
    }))} onColumnClicked={updateColumnSorting}>
      {
        paginatedTags.length == 0
          ? (
            <tr>
              <td colSpan={3} style={{ textAlign: 'center' }}>
                No Matching Tags
              </td>
            </tr>)
          : (paginatedTags[tagTableState.page - 1] ?? []).map(tag => (
            <tr key={tag.name}>
              <td>
                <Tag
                  name={tag.name}
                  color={tag.colorName}
                  onClick={() => props.onTagClick?.(tag)}
                />
              </td>
              {
                columnRenders.map((col, idx) => <td key={idx}>{col(tag)}</td>)
              }
              <td className={cx('button-cell')}>
                <ButtonGroup className={cx('row-buttons')}>
                  <Button small onClick={() => props.onEditClick(tag)}>Edit</Button>
                  <Button small onClick={() => props.onDeleteClick(tag)}>Delete</Button>
                </ButtonGroup>
              </td>
            </tr>
          ))
      }
    </Table>
    <div className={cx('button-block')}>
      {
        props.onCreateClick &&
        <Button className={cx('create-button')} onClick={() => props.onCreateClick?.()}>Create</Button>
      }

      <StandardPager
        page={tagTableState.page}
        onPageChange={pageNum => dispatch({ type: 'page-change', newPage: pageNum })}
        maxPages={paginatedTags.length}
      />
    </div>

  </>
}

const tagTableReducer = (state: TagTableState, action: TagTableAction): TagTableState => {
  if (action.type == 'page-change') {
    return { ...state, page: action.newPage }
  }
  if (action.type == 'sort-column') {
    return {
      ...state,
      ...action,
    }
  }
  if (action.type == 'filter-text-change') {
    return {
      ...state,
      ...action,
      page: 1
    }
  }
  return state
}

type compareableFunc = (l: unknown, r: unknown) => number

const sortNone: compareableFunc = (a: unknown, b: unknown) => 0
const sortNums: compareableFunc = (a: TagWithUsage, b: TagWithUsage) => a.evidenceCount - b.evidenceCount
const sortTags: compareableFunc = (a: TagWithUsage, b: TagWithUsage) => a.name.localeCompare(b.name)

type tagTableColumn<T> = {
  title: string,
  label: string
  clickable?: boolean
  compareVia: compareableFunc
  renderer: (tag: T) => React.ReactNode
}

type TagTableState = {
  page: number
  sortFunc: compareableFunc
  sortDir: SortDirection
  sortColIndex: number
  filterText: string
}

const TagTableInitialState = {
  page: 1,
  sortFunc: sortNone,
  sortDir: undefined,
  sortColIndex: 0,
  filterText: ""
}

type TagTableSortColumn = {
  type: 'sort-column'
  sortFunc: compareableFunc
  sortColIndex: number
  sortDir: SortDirection
}

type TagTableUpdatePage = {
  type: 'page-change'
  newPage: number
}

type TagTableFilterTextChange = {
  type: 'filter-text-change',
  filterText: string
}

type TagTableAction =
  | TagTableSortColumn
  | TagTableUpdatePage
  | TagTableFilterTextChange
