// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { Har, Entry, Log, Header, PostData, Response } from 'har-format'
import { EvidenceViewHint } from 'src/global_types'
import { trimURL, clamp } from 'src/helpers'
import { mimetypeToAceLang } from './helpers'
import { EvidenceHeader, RawContent } from './components'
import { ExpandableSection } from '../expandable_area'
import { Tab } from '../tab_vertical_menu'
import { default as TabMenu } from '../tabs'
import Table from 'src/components/table'

const cx = classnames.bind(require('./stylesheet'))


export * from './is_valid_har'

type ExpandedState = { [tabId: string]: boolean }
type EntryState = {
  activeTab: number
  expandedAreas: ExpandedState
}

export const HarViewer = (props: {
  log: Har
  viewHint?: EvidenceViewHint
  disableKeyHandler?: boolean
}) => {
  const log: Log = props.log.log
  const [selectedRow, setSelectedRow] = React.useState<number>(-1)
  const [entryState, setEntryState] = React.useState<EntryState>({
    activeTab: -1,
    expandedAreas: {}
  })

  return (
    <div className={cx('root')} onClick={e => {
      if (!props.disableKeyHandler ) {
        e.stopPropagation() // prevent lightbox from showing
      }
    }}>
      <EvidenceHeader creator={log.creator.name} version={log.creator.version} />
      <div className={cx('columns')}>
        <RequestTable log={log} selectedRow={selectedRow} setSelectedRow={props.disableKeyHandler ? (_)=>{} : setSelectedRow} />
        {selectedRow > -1  &&
          <EntryData entry={log.entries[selectedRow]} state={entryState} setState={setEntryState} />
        }
      </div>
    </div>
  )
}

const EntryData = (props: {
  entry: Entry,
  state: EntryState,
  setState: (newState: EntryState) => void
}) => {
  const onTabChanged = (_: Tab, idx: number) => props.setState({ ...props.state, activeTab: idx })

  return <div className={cx('entry-root', 'column')}>
    <TabMenu className={cx('tab-menu-group')} initialActiveTab={props.state.activeTab} onTabChanged={onTabChanged}
      tabs={[
        {
          id: 'entry-headers', label: 'Headers',
          content: (
            <EntryHeadersData entry={props.entry}
              expandedAreas={props.state.expandedAreas}
              setExpandedAreas={(s) => props.setState({ ...props.state, expandedAreas: s })}
            />
          )
        },
        {
          id: 'entry-request', label: 'Request',
          content: <RequestContent data={props.entry.request.postData} />
        },
        {
          id: 'entry-response', label: 'Response',
          content: <ResponseContent response={props.entry.response} />
        },
      ]}
    />
  </div>
}

const EntryHeadersData = (props: {
  entry: Entry
  expandedAreas: ExpandedState
  setExpandedAreas: (newState: ExpandedState) => void
}) => {
  const expandedAreaProps = (id: string) => ({
    onExpanded: (include: boolean) => props.setExpandedAreas({ ...props.expandedAreas, [id]: include }),
    initiallyExpanded: props.expandedAreas[id] || false
  })

  return (
    <div className={cx('headers-grouping')}>
      <ExpandableSection {...expandedAreaProps('request-info')} content={
        <RequestInfo entry={props.entry} />
      }>
        Request Info
    </ExpandableSection>
      <ExpandableSection {...expandedAreaProps('request-headers')} content={
        <SectionDefintions definitions={formatHeaders(props.entry.request.headers)} />
      }>
        Request Headers
    </ExpandableSection>
      <ExpandableSection {...expandedAreaProps('response-headers')} content={
        <SectionDefintions definitions={formatHeaders(props.entry.response.headers)} />
      }>
        Response Headers
    </ExpandableSection>
    </div>
  )
}

const formatHeaders = (headers: Array<Header>): Array<[string, string]> => headers
  .sort((a, b) => a.name.toLowerCase().localeCompare(b.name.toLowerCase()))
  .map(header => [header.name, header.value])

const RequestInfo = (props: {
  entry: Entry
}) => {
  const requestUrl = new URL(props.entry.request.url)
  return (
    <SectionDefintions
      definitions={[
        [props.entry.request.method, ''],
        ['Host', requestUrl.hostname + (requestUrl.port == '80' ? '' : `:${requestUrl.port}`)],
        ['Path', requestUrl.pathname],
        ['Server IP', props.entry.serverIPAddress || ''],
        ['Duration', props.entry.time.toFixed(2) + ' ms'],
      ]}
    />
  )
}

const SectionDefintions = (props: {
  definitions: Array<[key: string, value: string]>
}) => (
  <section className={cx('section-container')}>
    {props.definitions.map(([key, value]) => (
      <div className={cx('section-entry')}>
        <em className={cx('section-key')}>{key}{value && ':'}</em>
        <div className={cx('section-value')}>{value}</div>
      </div>
    ))}
  </section>
)

const RequestTable = (props: {
  log: Log
  selectedRow: number
  setSelectedRow: (rowNumber: number) => void
}) => {

  const tableRef = React.useRef<HTMLTableElement | null>(null)

  const onKeyDown = (e: KeyboardEvent) => {
    if (['ArrowUp', 'ArrowDown'].includes(e.key)) {
      e.stopPropagation()
      e.preventDefault()

      const direction = (e.key == 'ArrowDown' ? 1 : -1)
      const newIndex = clamp(props.selectedRow + direction, 0, props.log.entries.length - 1)
      if (tableRef.current != null) {

        if (props.selectedRow == 0 && newIndex == 0) {
          tableRef.current.tHead?.scrollIntoView({ block: 'nearest' })
        }
        else {
          tableRef.current.tBodies.item(0)?.rows.item(newIndex)?.scrollIntoView({ block: 'nearest' })
        }
      }
      props.setSelectedRow(newIndex)
    }
  }

  const onRowSelected = (index: number) => (e: React.MouseEvent<HTMLTableRowElement, MouseEvent>) => {
    props.setSelectedRow(props.selectedRow == index ? -1 : index)
  }

  return (
    <div className={cx('table-container', 'column', props.selectedRow > -1 ? '' : 'full-width-column')}>
      <Table className={cx('table')} columns={['#', 'Status', 'Method', 'Path']}
        onKeyDown={onKeyDown} tableRef={tableRef}>
        {props.log.entries.map((entry, index) => (
          <tr key={index} className={cx(index == props.selectedRow ? ['selected-row', 'render'] : '')}
            onClick={(e) => onRowSelected(index)(e)} >
            <td>{index + 1}</td>
            <td>{entry.response.status}</td>
            <td>{entry.request.method}</td>
            <td>{trimURL(entry.request.url).trimmedValue}</td>
          </tr>
        ))}
      </Table>
    </div>
  )
}

const RequestContent = (props: {
  data?: PostData
}) => {
  if (props.data == null) {
    return <RawContent content="No Post Data captured" />
  }

  const mimetype = (props.data?.mimeType) || ''

  // Per the draft HAR v1.2 standard, text and params are mutually exclusive.
  // However, in practice they are not (see chrome form data har export). Opting to prefer text
  // over params
  let body = ''

  if (props.data.text != null) {
    body = props.data.text
  }
  else {
    body = 'Parameters:\n'
    for (let p of props.data.params) {
      body += `  ${p.name}${(p.value ? ': ' + p.value : '')}\n`
    }
  }

  return <RawContent content={body} language={mimetypeToAceLang(mimetype)} />
}

const ResponseContent = (props: {
  response: Response
}) => {
  const length = props.response.content.size
  const rawText = props.response.content.text || ''

  const content = (rawText == '' && length > 0)
    ? `Content is ${length} bytes long, but no data/text was captured`
    : rawText

  return <RawContent content={content} language={mimetypeToAceLang(props.response.content.mimeType)} />
}
