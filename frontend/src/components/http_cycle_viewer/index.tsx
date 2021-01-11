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

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  log: Har
}) => {
  const log: Log = props.log.log

  const [selectedRow, setSelectedRow] = React.useState<number>(-1)

  return (
    <>
      <div className={cx('root')} onClick={e => e.stopPropagation()}>
        <div className={cx('header')}>From: <em>{log.creator.name} @ {log.creator.version}</em></div>
        <div className={cx('table-container')}>
          <Table columns={['#', 'Status', 'Method', 'Path', 'Data Size']} className={cx('table')}>
            {log.entries.map((entry, index) => (
              <tr key={index} className={cx(index == selectedRow ? ['selected-row', 'render'] : '')} onClick={() => setSelectedRow(index)} >
                <td>{index + 1}</td>
                <td>{entry.response.status}</td>
                <td>{entry.request.method}</td>
                <td>{trimURL(entry.request.url).trimmedValue}</td>
                <td>{entry.request.postData == null ? "None" : entry.request.postData.text?.length}</td>
              </tr>
            ))}
          </Table>
        </div>
        <div className={cx('body')}>
          {selectedRow == -1 ? null : <HttpEntry entry={log.entries[selectedRow]} />}
        </div>
      </div>
    </>
  )
}

const HttpEntry = (props: {
  entry: Entry
}) => (
  <>
    <RequestSegment {...props} />
    <ResponseSegment {...props} />
  </>
)

const RequestSegment = (props: {
  entry: Entry
}) => <>
    <SettingsSection className={cx('section-header')} title="Request" width="full-width">
      <TabMenu className={cx('tab-group')}
        tabs={[
          { id: "request-pretty", label: "Pretty", content: <PrettyHeaders headers={props.entry.request.headers} /> },
          { id: "request-raw", label: "Raw", content: <RawContent content={requestToRaw(props.entry.request)} /> },
          { id: "request-content", label: "Content", content: <RequestContent data={props.entry.request.postData} /> },
        ]}
      />
    </SettingsSection>
  </>

const ResponseSegment = (props: {
  entry: Entry
}) => <>
    <SettingsSection className={cx('section-header')} title="Response" width="full-width">
      <TabMenu className={cx('tab-group')}
        tabs={[
          { id: "request-pretty", label: "Pretty Headers", content: <PrettyHeaders headers={props.entry.response.headers} /> },
          { id: "response-raw", label: "Raw Headers", content: <RawContent content={responseToRaw(props.entry.response)} /> },
          { id: "response-content", label: "Content", content: <ResponseContent response={props.entry.response} /> },
        ]}
      />
    </SettingsSection>
  </>

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

