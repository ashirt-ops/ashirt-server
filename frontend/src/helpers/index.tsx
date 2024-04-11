// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { GlobalVariableData, OperationVariableData } from 'src/global_types'

export * from './clamp'
export * from './compute_delta'
export * from './query_parser'
export * from './tag_colors'
export * from './use_async_component'
export * from './use_element_rect'
export * from './use_form'
export * from './use_modal'
export * from './use_user_is_super_admin'
export * from './use_wired_data'
export * from './localStroage'
export * from './trim_url'
export * from './highlight_substring'
export * from './codeblock_to_blob'
export * from './c2event_to_blob'

export function useFocusFirstFocusableChild(ref: React.MutableRefObject<HTMLDivElement|null>) {
  React.useLayoutEffect(() => {
    if (!ref.current) return
    const el = ref
      .current
      .querySelector('button, [href], input, select, textarea, [tabindex]')
    // @ts-ignore - Typescript is unable to know that el is a focusable element
    if (el != null) el.focus()
  }, [ref])
}

export function useWindowSize(): {width: number, height: number} {
  const [size, setSize] = React.useState({width: window.innerWidth, height: window.innerHeight})

  const onResize = () => {
    setSize({width: window.innerWidth, height: window.innerHeight})
  }

  React.useEffect(() => {
    window.addEventListener('resize', onResize)
    return () => { window.removeEventListener('resize', onResize) }
  })

  return size
}

export function isOperationVariable(variable: GlobalVariableData | OperationVariableData): variable is OperationVariableData {
  return (variable as OperationVariableData).operationSlug !== undefined;
}
