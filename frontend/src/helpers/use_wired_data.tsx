// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ErrorDisplay from 'src/components/error_display'
import LoadingSpinner from 'src/components/loading_spinner'
import { PaginationResult } from 'src/global_types'

// useWiredData is a react hook helper to make it trivial to load data in a component.
// It takes a `fetchDataFn` which returns a promise, and it returns a render method that will
// call a passed renderer with the specified data when the data is ready
//
// Example:
// const SomeComponent = () => {
//   const wiredUsers = useWiredData(fetchUsers)
//
//   return <div>
//     {wiredUsers.render(users => (
//       <UserTable users={users} />
//     ))}
//   </div>
// }
//
//
// Handling dynamic fetchDataFns
// useWiredData will refetch data any time the fetchDataFn changes. This means using an inline
// function for useWiredData will result in constant fetching. To resolve this, you should use
// `React.useCallback` which memoizes the callback to prevent this. This will have the added
// benefit of automatically reloading data anytime any dependencies change and will be enforced
// by the `react-hooks/exhaustive-deps`  lint rule
//
// Example:
// const SomeComponent = (props: { userId: number }) => {
//   const wiredUser = useWiredData(React.useCallback(() => fetchUser(props.userId), [props.userId]))
//
//   return wiredUser.render(user => <ProfilePicture avatar={user.avatar} />)
// }
//
//
// Custom Error & Loading components
// If the data is not ready or there is an error returned from the fetchDataFn an appropriate
// loading spinner or error message will be rendered in place respectively
// These components can be customized by passing in a custom renderer in to the second and third arguments
//
// Example:
//   const wiredUser = useWiredData(
//     React.useCallback(() => fetchUser(props.userId), [props.userId]),
//     (err: Error) => <CustomErrorRenderer err={err} />,
//     () => <CustomLoadingRenderer />,
//   )
//
//
// Debouncing
// useWiredData will debounce repeated changes to the fetchDataFn at a rate of 250ms.
// This means the first call will be made immediately, however subsequent calls will be delayed
// until there has been at least 250ms of idle time, and then it will call the last update
//
// Example:
//
// fetchDataFn updates:   1..2.........3.4..5..6..........7...
// calls to fetchDataFn:  1.......2....3............6.....7...
//                                     time -->
//
// This allows you to pass in things like user input directly as a dependency to `React.useCallback`
// without sending many intermediate requests

type Renderer<T> = (data: T) => React.ReactElement
export type WiredData<T> = {
  loading: boolean,
  reload: () => void,
  render: (renderer: Renderer<T>) => React.ReactElement
}

export type PaginatedWiredData<T> = WiredData<Array<T>> & {
  pagerProps: {
    page: number,
    maxPages: number,
    onPageChange: (pageNumber: number) => void,
  }
}

const DEBOUNCE_TIMEOUT = 250 // ms

export function useWiredData<T>(
  fetchDataFn: () => Promise<T>,
  errorRenderer: Renderer<Error> = (err) => <ErrorDisplay err={err} />,
  loadingRenderer: Renderer<void> = () => <LoadingSpinner />,
): WiredData<T> {
  const [err, setErr] = React.useState<Error | null>(null)
  const [loading, setLoading] = React.useState(true)
  const [data, setData] = React.useState<{ value: T } | null>(null)
  const [isDirty, makeDirty] = React.useState(false)
  const [debouncedFetchDataFn, setDebouncedFetchDataFn] = React.useState(() => fetchDataFn)
  const isDebouncing = React.useRef(false)

  React.useEffect(() => {
    let timeoutFn: () => void
    if (isDebouncing.current) {
      timeoutFn = () => {
        isDebouncing.current = false
        setDebouncedFetchDataFn(() => fetchDataFn)
      }
    } else {
      setDebouncedFetchDataFn(() => fetchDataFn)
      isDebouncing.current = true
      timeoutFn = () => { isDebouncing.current = false }
    }

    const timeout = setTimeout(timeoutFn, DEBOUNCE_TIMEOUT)
    return () => { clearTimeout(timeout) }
  }, [fetchDataFn])

  React.useEffect(() => {
    setLoading(true)
    debouncedFetchDataFn()
      .then((data: T) => {
        setData({ value: data })
        setErr(null)
        setLoading(false)
      })
      .catch((err: Error) => {
        setErr(err)
        setLoading(false)
      })
  }, [isDirty, debouncedFetchDataFn])

  return {
    loading,
    reload() {
      if (!loading) makeDirty(!isDirty)
    },
    render(renderer: Renderer<T>) {
      if (err != null) return errorRenderer(err)
      if (loading || data == null) return loadingRenderer()
      return renderer(data.value)
    },
  }
}

export function usePaginatedWiredData<T>(
  fetchDataFn: (pageNumber: number) => Promise<PaginationResult<T>>,
  errorRenderer: Renderer<Error> = (err) => <ErrorDisplay err={err} />,
  loadingRenderer: Renderer<void> = () => <LoadingSpinner />,
): PaginatedWiredData<T> {
  const [pageNumber, setPageNumber] = React.useState(1)
  const [totalPages, setTotalPages] = React.useState(1)

  const memoizedFetchDataFn = React.useCallback(async (): Promise<Array<T>> => {
    const data = await fetchDataFn(pageNumber)
    setTotalPages(data.totalPages)
    return data.content
  }, [fetchDataFn, pageNumber])
  const wiredData = useWiredData<Array<T>>(memoizedFetchDataFn, errorRenderer, loadingRenderer)

  return {
    ...wiredData,
    pagerProps: {
      page: pageNumber,
      maxPages: totalPages,
      onPageChange: setPageNumber,
    }
  }
}
