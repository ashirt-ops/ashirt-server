import { describe, it, expect, vi } from 'vitest'
import { render, waitFor, act } from '@testing-library/react'
import { renderHook } from '@testing-library/react'
import { useWiredData, usePaginatedWiredData } from './use_wired_data'

const loadingRenderer = () => <div>loading</div>
const errorRenderer = (err: Error) => <div>error: {err.message}</div>

describe('useWiredData', () => {
  it('starts in loading state', () => {
    const fetchFn = vi.fn(() => new Promise(() => {}))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    expect(result.current.loading).toBe(true)
  })

  it('renders loading renderer while fetch is pending', () => {
    const fetchFn = vi.fn(() => new Promise(() => {}))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    const { getByText } = render(
      <>
        {result.current.render(() => (
          <div>success</div>
        ))}
      </>,
    )
    expect(getByText('loading')).toBeInTheDocument()
  })

  it('transitions to non-loading after fetch resolves', async () => {
    const fetchFn = vi.fn().mockResolvedValue({ name: 'Alice' })
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
  })

  it('renders data after fetch resolves', async () => {
    const fetchFn = vi.fn(() => Promise.resolve({ name: 'Alice' }))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
    const { getByText } = render(
      <>
        {result.current.render((data) => (
          <div>data: {data.name}</div>
        ))}
      </>,
    )
    expect(getByText('data: Alice')).toBeInTheDocument()
  })

  it('renders error renderer when fetch rejects', async () => {
    const fetchFn = vi.fn().mockRejectedValue(new Error('network error'))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
    const { getByText } = render(
      <>
        {result.current.render(() => (
          <div>success</div>
        ))}
      </>,
    )
    expect(getByText('error: network error')).toBeInTheDocument()
  })

  it('calls fetchDataFn exactly once on mount', async () => {
    const fetchFn = vi.fn().mockResolvedValue('data')
    renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(fetchFn).toHaveBeenCalledTimes(1))
  })

  it('expose() does not call exposer when data is not yet loaded', () => {
    const fetchFn = vi.fn(() => new Promise(() => {}))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    const exposer = vi.fn()
    result.current.expose(exposer)
    expect(exposer).not.toHaveBeenCalled()
  })

  it('expose() calls exposer with data after load', async () => {
    const fetchFn = vi.fn().mockResolvedValue({ name: 'Bob' })
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
    const exposer = vi.fn()
    result.current.expose(exposer)
    expect(exposer).toHaveBeenCalledWith({ name: 'Bob' })
  })

  it('reload() triggers a re-fetch', async () => {
    const fetchFn = vi.fn().mockResolvedValue('data')
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(fetchFn).toHaveBeenCalledTimes(1)
    act(() => result.current.reload())
    await waitFor(() => expect(fetchFn).toHaveBeenCalledTimes(2))
  })

  it('uses a custom error renderer', async () => {
    const custom = (err: Error) => <div>custom: {err.message}</div>
    const fetchFn = vi.fn().mockRejectedValue(new Error('oops'))
    const { result } = renderHook(() => useWiredData(fetchFn, custom, loadingRenderer))
    await waitFor(() => expect(result.current.loading).toBe(false))
    const { getByText } = render(
      <>
        {result.current.render(() => (
          <div>ok</div>
        ))}
      </>,
    )
    expect(getByText('custom: oops')).toBeInTheDocument()
  })

  it('uses a custom loading renderer', () => {
    const custom = () => <div>custom loading</div>
    const fetchFn = vi.fn(() => new Promise(() => {}))
    const { result } = renderHook(() => useWiredData(fetchFn, errorRenderer, custom))
    const { getByText } = render(
      <>
        {result.current.render(() => (
          <div>ok</div>
        ))}
      </>,
    )
    expect(getByText('custom loading')).toBeInTheDocument()
  })
})

describe('usePaginatedWiredData', () => {
  it('starts at page 1 with maxPages 1', () => {
    const fetchFn = vi.fn((_page: number) => new Promise<never>(() => {}))
    const { result } = renderHook(() =>
      usePaginatedWiredData(fetchFn, errorRenderer, loadingRenderer),
    )
    expect(result.current.pagerProps.page).toBe(1)
    expect(result.current.pagerProps.maxPages).toBe(1)
  })

  it('calls fetchFn with the current page number', async () => {
    const fetchFn = vi.fn((_page: number) =>
      Promise.resolve({ content: [], totalPages: 1, totalCount: 0, page: 1, pageSize: 10 }),
    )
    const { result } = renderHook(() =>
      usePaginatedWiredData(fetchFn, errorRenderer, loadingRenderer),
    )
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(fetchFn).toHaveBeenCalledWith(1)
  })

  it('parses totalPages from the fetch result', async () => {
    const fetchFn = vi.fn((_page: number) =>
      Promise.resolve({ content: [], totalPages: 7, totalCount: 0, page: 1, pageSize: 10 }),
    )
    const { result } = renderHook(() =>
      usePaginatedWiredData(fetchFn, errorRenderer, loadingRenderer),
    )
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.pagerProps.maxPages).toBe(7)
  })

  it('renders paginated content', async () => {
    const fetchFn = vi.fn((_page: number) =>
      Promise.resolve({ content: ['x', 'y'], totalPages: 1, totalCount: 2, page: 1, pageSize: 10 }),
    )
    const { result } = renderHook(() =>
      usePaginatedWiredData(fetchFn, errorRenderer, loadingRenderer),
    )
    await waitFor(() => expect(result.current.loading).toBe(false))
    const { getByText } = render(
      <>
        {result.current.render((data) => (
          <div>{data.join(',')}</div>
        ))}
      </>,
    )
    expect(getByText('x,y')).toBeInTheDocument()
  })

  it('re-fetches when page changes', async () => {
    const fetchFn = vi.fn((_page: number) =>
      Promise.resolve({ content: [], totalPages: 3, totalCount: 0, page: 1, pageSize: 10 }),
    )
    const { result } = renderHook(() =>
      usePaginatedWiredData(fetchFn, errorRenderer, loadingRenderer),
    )
    await waitFor(() => expect(result.current.loading).toBe(false))
    act(() => result.current.pagerProps.onPageChange(2))
    await waitFor(() => expect(fetchFn).toHaveBeenCalledWith(2))
    expect(result.current.pagerProps.page).toBe(2)
  })
})
