import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { useRef } from 'react'
import { renderHook } from '@testing-library/react'
import Table, { SortAsc, SortDesc, type ColumnData } from './index'

describe('Table', () => {
  it('renders string column headers', () => {
    render(
      <Table columns={['Name', 'Age']}>
        <tr>
          <td>Alice</td>
          <td>30</td>
        </tr>
      </Table>,
    )
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('Age')).toBeInTheDocument()
  })

  it('renders ColumnData object headers', () => {
    const cols: ColumnData[] = [{ label: 'Name', title: 'Sort by name', clickable: true }]
    render(
      <Table columns={cols}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    expect(screen.getByText('Name')).toBeInTheDocument()
  })

  it('renders children in tbody', () => {
    render(
      <Table columns={['Name']}>
        <tr>
          <td>Alice</td>
        </tr>
        <tr>
          <td>Bob</td>
        </tr>
      </Table>,
    )
    expect(screen.getByText('Alice')).toBeInTheDocument()
    expect(screen.getByText('Bob')).toBeInTheDocument()
  })

  it('calls onColumnClicked with column index when a clickable column is clicked', () => {
    const onColumnClicked = vi.fn()
    const cols: ColumnData[] = [
      { label: 'Name', title: 'Sort', clickable: true },
      { label: 'Age', title: 'Sort age', clickable: true },
    ]
    render(
      <Table columns={cols} onColumnClicked={onColumnClicked}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    screen.getByText('Name').click()
    expect(onColumnClicked).toHaveBeenCalledWith(0)
  })

  it('does not call onColumnClicked for non-clickable columns', () => {
    const onColumnClicked = vi.fn()
    const cols: ColumnData[] = [{ label: 'Name', title: 'Name', clickable: false }]
    render(
      <Table columns={cols} onColumnClicked={onColumnClicked}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    screen.getByText('Name').click()
    expect(onColumnClicked).not.toHaveBeenCalled()
  })

  it('renders sort ascending indicator', () => {
    const cols: ColumnData[] = [{ label: 'Name', title: 'Name', sortDirection: SortAsc }]
    const { container } = render(
      <Table columns={cols}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    expect(container.querySelector('.asc')).toBeInTheDocument()
  })

  it('renders sort descending indicator', () => {
    const cols: ColumnData[] = [{ label: 'Name', title: 'Name', sortDirection: SortDesc }]
    const { container } = render(
      <Table columns={cols}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    expect(container.querySelector('.desc')).toBeInTheDocument()
  })

  it('sets tabIndex when onKeyDown is provided', () => {
    const { container } = render(
      <Table columns={['Name']} onKeyDown={() => {}}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    expect(container.querySelector('table')).toHaveAttribute('tabindex', '0')
  })

  it('does not set tabIndex when onKeyDown is not provided', () => {
    const { container } = render(
      <Table columns={['Name']}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )
    expect(container.querySelector('table')).not.toHaveAttribute('tabindex')
  })

  it('fires keydown events via tableRef', () => {
    const onKeyDown = vi.fn()
    const { result } = renderHook(() => useRef<HTMLTableElement | null>(null))
    const tableRef = result.current

    render(
      <Table columns={['Name']} onKeyDown={onKeyDown} tableRef={tableRef}>
        <tr>
          <td>-</td>
        </tr>
      </Table>,
    )

    if (tableRef.current) {
      fireEvent.keyDown(tableRef.current, { key: 'ArrowDown' })
    }
    expect(onKeyDown).toHaveBeenCalled()
  })
})
