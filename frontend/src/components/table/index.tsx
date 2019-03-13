// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children: React.ReactNode,
  className?: string,
  columns: Array<string | {label: string, title: string}>,
}) => (
  <table className={cx('root', props.className)}>
    <thead>
      <tr>
        {props.columns.map((column, idx) => (
          typeof (column) === 'object'
            ? <th key={`${column.label}-${idx}`} title={column.title}>{column.label}</th>
            : <th key={`${column}-${idx}`}>{column}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      {props.children}
    </tbody>
  </table>
)
