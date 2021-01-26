// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import DateRangePicker from 'src/components/date_range_picker'
import Input from 'src/components/input'
import classnames from 'classnames/bind'
import { default as Button, ButtonGroup } from 'src/components/button'

import { stringToSearch, SearchType, SearchOptions, stringifySearch } from 'src/components/search_query_builder/helpers'
import Modal from 'src/components/modal'
import { default as SearchQueryBuilder } from 'src/components/search_query_builder'

import { getDateRangeFromQuery, addOrUpdateDateRangeInQuery, useModal, renderModals } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onRequestCreateFinding: () => void,
  onRequestCreateEvidence: () => void,
  onSearch: (query: string) => void,
  query: string,
}) => {
  const [queryInput, setQueryInput] = React.useState<string>("")
  React.useEffect(() => {
    setQueryInput(props.query)
  }, [props.query])

  const inputRef = React.useRef<HTMLInputElement>(null)
  const builderModal = useModal<{searchText: string}>(modalProps => (
    <SearchBuilderModal
      {...modalProps}
      onChanged={(result: string) => setQueryInput(result)}
      operationSlug={"HPCoS"} // TODO
      searchType={SearchType.EVIDENCE_SEARCH} // TODO
    />
  ))

  return (
    <div className={cx('root')}>
      <Input
        ref={inputRef}
        className={cx('search')}
        value={queryInput}
        onChange={setQueryInput}
        placeholder="Filter Timeline"
        icon={require('./search.svg')}
        onKeyDown={e => {
          if (e.which == 13) {
            inputRef.current?.blur()
            props.onSearch(queryInput)
          }
        }}
      />
      <Button onClick={() => builderModal.show({searchText: queryInput})} >Help Me!</Button>
      <DateRangePicker
        range={getDateRangeFromQuery(queryInput)}
        onSelectRange={r => {
          const newQuery = addOrUpdateDateRangeInQuery(queryInput, r)
          setQueryInput(newQuery)
          props.onSearch(newQuery)
        }}
      />

      <ButtonGroup>
        <Button onClick={props.onRequestCreateFinding}>Create Finding</Button>
        <Button onClick={props.onRequestCreateEvidence}>Create Evidence</Button>
      </ButtonGroup>
      {renderModals(builderModal)}
    </div>
  )
}

const SearchBuilderModal = (props: {
  searchText: string,
  operationSlug: string,
  searchType: SearchType,
  onRequestClose: () => void,
  onChanged: (resultString: string) => void,
}) => {

  return <Modal title="Query Builder" onRequestClose={props.onRequestClose}>
    <SearchQueryBuilder
      searchOptions={stringToSearch(props.searchText)}
      onChanged={(result:SearchOptions) => {
        props.onChanged(stringifySearch(result))
        props.onRequestClose()
      }}
      operationSlug={props.operationSlug}
      searchType={props.searchType}
    />
  </Modal>
}
