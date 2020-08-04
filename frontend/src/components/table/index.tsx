// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export const SortAsc = 'asc'
export const SortDesc = 'desc'

export type SortDirection = typeof SortAsc | typeof SortDesc | undefined

export type ColumnData = {
  label: string,
  title: string,
  sortDirection?: SortDirection,
  clickable?: boolean
}

export default (props: {
  children: React.ReactNode,
  className?: string,
  columns: Array<string | ColumnData>,
  onColumnClicked?: (colIndex: number)=>void,
}) => (
    <table className={cx('root', props.className)}>
      <thead>
        <tr>
          {props.columns.map((column, idx) => {
            if (typeof (column) === 'object') {
              const args = {
                key: `${column.label}-${idx}`,
                title: column.title,
                style: column.clickable ? { cursor: "pointer" } : {},
                onClick: () => (column.clickable && props.onColumnClicked) ? props.onColumnClicked(idx) : undefined
              }
              return (
                <th {...args}>
                  {column.label}
                  {column.sortDirection ? <span className={cx('arrow', column.sortDirection == SortAsc ? 'asc' : 'desc')}></span> : null}
                </th>
              )
            }
            return <th key={`${column}-${idx}`}>{column}</th>
          })}
        </tr>
      </thead>
      <tbody>
        {props.children}
      </tbody>
    </table>
  )
