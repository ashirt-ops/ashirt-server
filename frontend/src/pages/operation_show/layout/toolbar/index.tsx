// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import DateRangePicker from 'src/components/date_range_picker'
import Input from 'src/components/input'
import classnames from 'classnames/bind'
import {default as Button, ButtonGroup} from 'src/components/button'
import {getDateRangeFromQuery, addOrUpdateDateRangeInQuery} from 'src/helpers'
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
    </div>
  )
}
