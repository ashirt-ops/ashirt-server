// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { trimURL, useAsyncComponent } from 'src/helpers'

import Table from 'src/components/table'
import { Har, Entry, Request, Response, Header, PostData } from 'har-format'
import { default as TabMenu } from '../tabs'
import SettingsSection from 'src/components/settings_section'

const cx = classnames.bind(require('./stylesheet'))
const importAceEditorAsync = () => import('../code_block/ace_editor').then(module => module.default)

export default (props: {
  log: Har
}) => {
  const parsedLog: Har = props.log

  const [selectedRow, setSelectedRow] = React.useState<number>(-1)

  return (
    <>
      <div className={cx('root')} onClick={e => e.stopPropagation()}>
        <div className={cx('header')}>From: <em>{parsedLog.log.creator.name} @ {parsedLog.log.creator.version}</em></div>
        <div className={cx('table-container')}>
          <Table columns={['#', 'Status', 'Method', 'Path', 'Data Size']} className={cx('table')}>
            {parsedLog.log.entries.map((entry, index) => (
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
          {selectedRow == -1 ? null : <HttpEntry entry={parsedLog.log.entries[selectedRow]} />}
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

const PrettyHeaders = (props: {
  headers: Array<Header>
}) => {
  let content
  if (props.headers.length == 0) {
    content = [<em className={cx('pretty-headers-no-content')}>No Captured Headers</em>]
  }
  else {
    content = props.headers
      .sort((a, b) => a.name.toLowerCase().localeCompare(b.name.toLowerCase()))
      .map((h, i) => (
        <div key={i} className={cx('pretty-headers-entry')}>
          <em className={cx('pretty-headers-key')}>{h.name}:</em>
          <span className={cx('pretty-headers-value')}>{h.value}</span>
        </div>
      ))
  }

  return (
    <div className={cx('pretty-headers-outer-container')}>
      <div className={cx('pretty-headers-container')}>{...content}</div>
    </div>
  )
}

const RawContent = (props: {
  content: string
  language?: string
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('ace-container')}>
      <div className={cx('ace')}>
        <AceEditor
          readOnly
          mode={props.language ? props.language : ''}
          value={props.content}
        />
      </div>
    </div>
  )
}

const ResponseContent = (props: {
  response: Response
}) => {

  const length = props.response.content.size
  const rawText = props.response.content.text || ''
  const mimetype = props.response.content.mimeType || ''

  const content = (rawText == '' && length > 0)
    ? `Content is ${length} bytes long, but no data/text was captured`
    : rawText

  return <RawContent content={content} language={mimetypeToAceLang(mimetype)} />
}

const RequestContent = (props: {
  data?: PostData
}) => {

  if (props.data == null) {
    return <RawContent content="No Post Data captured" />
  }

  const mimetype = props.data.mimeType

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

const mimetypeToAceLang = (mimetype:string) => {
  if (mimetype.includes("text/javascript") || mimetype.includes('application/json')) {
    return 'javascript'
  }
  else if (mimetype.includes('text/html')) {
    return 'html'
  }
  else if (mimetype.includes('text/css')) {
    return 'css'
  }
  else if (mimetype.includes('text/xml')) {
    return 'xml'
  }
  return ''
}

const requestToRaw = (req: Request) => {
  const parsedUrl = new URL(req.url)
  const reqSummary = req.method + " " + parsedUrl.pathname + parsedUrl.search + " " + req.httpVersion + "\n"

  return reqSummary + req.headers.map(h => `${ h.name }: ${ h.value } `).join("\n")
}

const responseToRaw = (resp: Response) => {
  return resp.headers.map(h => `${ h.name }: ${ h.value } `).join("\n")
}
