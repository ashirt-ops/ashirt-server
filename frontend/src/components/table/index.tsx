import { type CSSProperties, type ReactNode, type MutableRefObject, useEffect } from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export * from './error_row'
export * from './loading_row'
export * from './spanned_row'

export const SortAsc = 'asc'
export const SortDesc = 'desc'

export type SortDirection = typeof SortAsc | typeof SortDesc | undefined

export type ColumnData = {
  label: string
  title: string
  sortDirection?: SortDirection
  clickable?: boolean
  style?: CSSProperties
}

export const StartAlignedColumn: CSSProperties = {
  display: 'flex',
  justifyContent: 'flex-start',
}
export const EndAlignedColumn: CSSProperties = { display: 'flex', justifyContent: 'flex-end' }
export const CenterAlignedColumn: CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
}

export default function Table(props: {
  children: ReactNode
  className?: string
  columns: Array<string | ColumnData>
  onColumnClicked?: (colIndex: number) => void
  onKeyDown?: (e: KeyboardEvent) => void
  tableRef?: MutableRefObject<HTMLTableElement | null>
}) {
  const noop = () => {}
  useEffect(() => {
    const curRootRef = props.tableRef?.current
    if (!curRootRef) {
      return
    }
    curRootRef.addEventListener('keydown', props.onKeyDown || noop)
    return () => {
      curRootRef.removeEventListener('keydown', props.onKeyDown || noop)
    }
  }, [props.onKeyDown, props.tableRef])

  return (
    <table
      className={cx('root', props.className)}
      {...(props.onKeyDown ? { tabIndex: 0 } : {})}
      {...(props.tableRef ? { ref: props.tableRef } : {})}
    >
      <thead>
        <tr>
          {props.columns.map((column, idx) => {
            if (typeof column === 'object') {
              const args = {
                key: `${column.label}-${idx}`,
                title: column.title,
                style: {
                  ...column.style,
                  ...(column.clickable ? { cursor: 'pointer' } : {}),
                },
                onClick: () =>
                  column.clickable && props.onColumnClicked
                    ? props.onColumnClicked(idx)
                    : undefined,
              }
              return (
                <th {...args}>
                  {column.label}
                  {column.sortDirection ? (
                    <span
                      className={cx('arrow', column.sortDirection == SortAsc ? 'asc' : 'desc')}
                    ></span>
                  ) : null}
                </th>
              )
            }
            return <th key={`${column}-${idx}`}>{column}</th>
          })}
        </tr>
      </thead>
      <tbody>{props.children}</tbody>
    </table>
  )
}
