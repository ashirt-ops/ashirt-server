// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { Tag, User, ViewName } from 'src/global_types'

import { default as Button, ButtonGroup } from 'src/components/button'
import { FilterFieldsGrid } from 'src/components/filter_fields/filter-field-grid'
import Input from 'src/components/input'
import { stringToSearch, SearchOptions, stringifySearch } from 'src/components/filter_fields/helpers'
import { useWiredData, useModal, renderModals } from 'src/helpers'
import { getTags, listEvidenceCreators } from 'src/services'
import { CreateButtonPosition } from '..'
import { SearchHelpModal } from './search_modal'

const cx = classnames.bind(require('./stylesheet'))


export const Toolbar = (props: {
  operationSlug: string,
  expandedView: boolean,
  query: string,
  queryName?: string
  setExpandedView: (expand: boolean) => void,
  viewName: ViewName,
  onSearch: (query: string) => void
  onRequestCreateFinding: () => void
  onRequestCreateEvidence: () => void
  requestQueriesReload?: () => void
  showCreateButtons: CreateButtonPosition
}) => {
  const [queryString, setQueryString] = React.useState<string>(props.query)
  const inputRef = React.useRef<HTMLInputElement>(null)

  React.useEffect(() => {
    setQueryString(props.query)
  }, [props.query])

  const wiredData = useWiredData<[Array<Tag>, Array<User>]>(
    React.useCallback(() =>
      Promise.all([
        getTags({ operationSlug: props.operationSlug }),
        listEvidenceCreators({ operationSlug: props.operationSlug }),
      ]), [props.operationSlug]
    ))

  return (
    <div className={cx('toolbar-root')}>
      {wiredData.render(([tags, users]) => {
        const searchOptions = stringToSearch(queryString, tags, users)
        return (
          <>
            <div className={cx('toolbar-flex')}>
              <Button
                className={cx('tb-expand-button')}
                onClick={() => props.setExpandedView(!props.expandedView)}
              />

              <div className={cx('tb-content')}>
                {!props.expandedView && (
                  <SearchInput
                    queryString={queryString}
                    inputRef={inputRef}
                    setQueryString={setQueryString}
                    onSearch={props.onSearch}
                  />
                )}
              </div>

              {props.showCreateButtons === 'filter' && (
                <ButtonGroup className={cx('tb-create-buttons')}>
                  <Button onClick={props.onRequestCreateEvidence}>Create Evidence</Button>
                  <Button onClick={props.onRequestCreateFinding}>Create Finding</Button>
                </ButtonGroup>
              )}
            </div>
            {props.expandedView && (
              <div className={cx('toolbar-overlay')}>
                <ExpandedSearch
                  {...props}
                  searchOptions={searchOptions}
                  setQueryString={setQueryString}
                  requestQueriesReload={props.requestQueriesReload}
                  queryName={props.queryName}
                />
              </div>
            )}

          </>
        )
      })}
    </div>
  )
}

const ExpandedSearch = (props: {
  operationSlug: string
  viewName: ViewName
  searchOptions: SearchOptions
  queryName?: string
  setExpandedView: (expand: boolean) => void
  onSearch: (query: string) => void
  setQueryString: (query: string) => void
  requestQueriesReload?: () => void
}) => {
  const { operationSlug, viewName, searchOptions,
    setQueryString, onSearch, setExpandedView, requestQueriesReload, } = props

  return (
    <FilterFieldsGrid
      className={cx('filter-grid')}
      operationSlug={operationSlug}
      viewName={viewName}
      queryName={props.queryName}
      withButtonRow
      value={searchOptions}
      requestQueriesReload={requestQueriesReload}
      onChange={(v) => setQueryString(stringifySearch(v))}
      onCanceled={() => setExpandedView(false)}
      onSubmit={(options) => {
        setExpandedView(false)
        const updatedQueryString = stringifySearch(options)
        setQueryString(updatedQueryString)
        onSearch(updatedQueryString)
      }}
    />
  )
}

const SearchInput = (props: {
  queryString: string
  inputRef: React.RefObject<HTMLInputElement>
  onSearch: (query: string) => void
  setQueryString: (query: string) => void
}) => {
  const { inputRef, queryString, setQueryString } = props
  const helpModal = useModal<void>(modalProps => <SearchHelpModal {...modalProps} />)
  return (
    <>
      <div className={cx('tb-search-container')}>
        <Input
          ref={inputRef}
          className={cx('tb-search')}
          inputClassName={cx('tb-search-input')}
          value={queryString}
          onChange={setQueryString}
          placeholder="Filter Timeline"
          icon={require('./search.svg')}
          onKeyDown={e => {
            if (e.key === 'Enter') {
              inputRef.current?.blur()
              props.onSearch(queryString)
            }
          }}
        />
        <a className={cx('search-help-icon')} onClick={() => helpModal.show()} title="Search Help"></a>
      </div>
      {renderModals(helpModal)}
    </>
  )
}
