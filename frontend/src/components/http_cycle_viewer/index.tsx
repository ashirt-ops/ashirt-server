// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { trimURL } from 'src/helpers'
import { Har, Entry, Response, PostData, Log } from 'har-format'
import { mimetypeToAceLang, requestToRaw, responseToRaw } from './helpers'

import { PrettyHeaders, RawContent } from './components'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { default as TabMenu } from '../tabs'
import { clamp } from 'lodash'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  log: Har
}) => {
  const log: Log = props.log.log
  const [selectedRow, setSelectedRow] = React.useState<number>(-1)

  return (
    <div className={cx('root')} onClick={e => e.stopPropagation()}>
      <EvidenceHeader creator={log.creator.name} version={log.creator.version}/>
      <RequestTable log={log} selectedRow={selectedRow} setSelectedRow={setSelectedRow}/>
      {selectedRow > -1 && <HttpDetails entry={log.entries[selectedRow]} />}
    </div>
  )
}

const EvidenceHeader = (props:{
  creator: string,
  version: string
}) => (
  <div className={cx('header')}>
    From:
    <em className={cx('header-creator')}>
      {props.creator}
    </em>
    @
    <em className={cx('header-version')}>
      {props.version}
    </em>
  </div>
)

const RequestTable = (props: {
  log: Log
  selectedRow: number
  setSelectedRow: (rowNumber: number) => void
}) => {
  const tbodyRef = React.useRef<HTMLTableSectionElement | null>(null)
  const onKeyDown = (e: KeyboardEvent) => {
    if (['ArrowUp', 'ArrowDown'].includes(e.key)) {
      e.stopPropagation()
      e.preventDefault()

      const newIndex = clamp(props.selectedRow + (e.key == 'ArrowDown' ? 1 : -1), 0, props.log.entries.length - 1)
      if (tbodyRef.current != null) {
        // @ts-ignore - typescript is unable to determine that children is an array of HTMLDivElements
        const rows: Array<HTMLTableRowElement> = Array.from(tbodyRef.current.children).filter(el => el instanceof HTMLTableRowElement)
        rows[newIndex].scrollIntoView() // this needs adjustment to be less weird/active
      }
      props.setSelectedRow(newIndex)
    }
  }

  return (
    <div className={cx('table-container')}>
      <Table className={cx('table')} tbodyRef={tbodyRef} onKeyDown={onKeyDown}
        columns={['#', 'Status', 'Method', 'Path', 'Data Size']} >
        {props.log.entries.map((entry, index) => (
          <tr key={index} className={cx(index == props.selectedRow ? ['selected-row', 'render'] : '')}
            onClick={() => props.setSelectedRow(index)} >
            <td>{index + 1}</td>
            <td>{entry.response.status}</td>
            <td>{entry.request.method}</td>
            <td>{trimURL(entry.request.url).trimmedValue}</td>
            <td>{entry.request.postData == null ? "None" : entry.request.postData.text?.length}</td>
          </tr>
        ))}
      </Table>
    </div>
  )
}

const HttpDetails = (props: {
  entry: Entry
}) => (
  <div className={cx('http-details')}>
    <RequestSegment {...props} />
    <ResponseSegment {...props} />
  </div>
)

const RequestSegment = (props: {
  entry: Entry
}) => (
  <SettingsSection className={cx('section-header')} title="Request" width="full-width">
    <TabMenu className={cx('tab-group')}
      tabs={[
        { id: "request-pretty", label: "Pretty Headers", content: <PrettyHeaders headers={props.entry.request.headers} /> },
        { id: "request-raw", label: "Raw Headers", content: <RawContent content={requestToRaw(props.entry.request)} /> },
        { id: "request-content", label: "Post Content", content: <RequestContent data={props.entry.request.postData} /> },
      ]}
    />
  </SettingsSection>
)

const ResponseSegment = (props: {
  entry: Entry
}) => (
  <SettingsSection className={cx('section-header')} title="Response" width="full-width">
    <TabMenu className={cx('tab-group')}
      tabs={[
        { id: "request-pretty", label: "Pretty Headers", content: <PrettyHeaders headers={props.entry.response.headers} /> },
        { id: "response-raw", label: "Raw Headers", content: <RawContent content={responseToRaw(props.entry.response)} /> },
        { id: "response-content", label: "Content", content: <ResponseContent response={props.entry.response} /> },
      ]}
    />
  </SettingsSection>
)

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

