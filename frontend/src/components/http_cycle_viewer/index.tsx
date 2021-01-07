// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { trimURL } from 'src/helpers'

import Table from 'src/components/table'
import { Har, Entry } from 'har-format'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  log: string
}) => {
  const parsedLog: Har = JSON.parse(props.log)  // TODO: moveoutside of this component

  const [selectedRow, setSelectedRow] = React.useState<number>(-1)

  return (
    <>
      <div className={cx('header')}>From: <em>{parsedLog.log.creator.name} @ {parsedLog.log.creator.version}</em></div>
      <div className={cx('table-container')}>
        <Table columns={['#', 'Status', 'Method', 'Path', 'Data Size']} className={cx('table')}>
          {parsedLog.log.entries.map((entry, index) => (
            <tr key={index} className={cx(index == selectedRow ? ['selected-row', 'render'] : '')} onClick={() => setSelectedRow(index)} >
              <td>{index + 1}</td>
              <td>{entry.response.status}</td>
              <td>{entry.request.method}</td>
              <td>{trimURL(entry.request.url).trimmedValue}</td>
              <td>{entry.request.postData == null ? "No Data" : entry.request.postData.text?.length}</td>
            </tr>
          ))}
        </Table>
      </div>
      <div className={cx('body')}>
        {selectedRow == -1 ? null : <HttpEntry entry={parsedLog.log.entries[selectedRow]} />}
      </div>
    </>
  )
}

const HttpEntry = (props: {
  entry: Entry
}) => {
  return <div>Stuff!</div>
}
