import * as React from 'react'
import Checkbox from 'src/components/checkbox'
import Input from 'src/components/input'
import classnames from 'classnames/bind'
import { useWiredData } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default function <T extends { uuid: string }>(props: {
  fetch: (query: string) => Promise<Array<T>>,
  renderRow: (result: T) => React.ReactNode,
  disabled?: boolean,
  onChange: (v: Array<T>) => void,
  placeholder: string,
  value: Array<T>,
  includeSelectAll?: true
}) {
  const [query, setQuery] = React.useState<string>('')
  const dataFetcher = props.fetch
  const wiredResults = useWiredData<Array<T>>(React.useCallback(() => dataFetcher(query), [query, dataFetcher]))
  const fetchedData = React.useRef<Array<T>>()
  const [allSelected, setAllSelected] = React.useState(false)

  const resultsByUuid: { [uuid: string]: T } = {}
  props.value.forEach(result => { resultsByUuid[result.uuid] = result })

  const getOnChangeHandler = (row: T) => (selected: boolean) => {
    if (allSelected) { // turn off to allow (fullset - 1) type operations
      setAllSelected(false)
    }
    const value = props.value.filter(r => r.uuid !== row.uuid)
    if (selected) value.push(row)
    props.onChange(value)
  }

  return (
    <div className={cx({ disabled: props.disabled })}>
      <Input
        placeholder={props.placeholder}
        value={query}
        onChange={setQuery}
        disabled={props.disabled}
        onKeyDown={e => e.key === 'Enter' && e.preventDefault()}
      />
      {
        props.includeSelectAll && (
          <Checkbox value={allSelected} className={cx('select-all-cb')} label='Select All Shown' onChange={(all) => {
            setAllSelected(all)
            if (all) {
              if (fetchedData.current) {
                props.onChange(fetchedData.current)
              }
            }
            else {
              props.onChange([])
            }
          }} />
        )
      }
      <div className={cx('results')}>
        {wiredResults.render(results => {
          fetchedData.current = results
          return (<>
            {results.map(result => (
              <Row
                key={result.uuid}
                selected={resultsByUuid[result.uuid] != null || allSelected}
                onChange={getOnChangeHandler(result)}
                children={props.renderRow(result)}
              />
            ))}
          </>)
        })}
      </div>
    </div>
  )
}

const Row = (props: {
  selected: boolean,
  onChange: (v: boolean) => void,
  children: React.ReactNode,
}) => (
  <div className={cx('row', { selected: props.selected })} onClick={() => props.onChange(!props.selected)}>
    <Checkbox className={cx('checkbox')} value={props.selected} onChange={props.onChange} />
    <div className={cx('children')}>
      {props.children}
    </div>
  </div>
)
