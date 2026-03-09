import { type MutableRefObject, useLayoutEffect, useState, useEffect } from 'react'
import { type GlobalVariableData, type OperationVariableData } from 'src/global_types'

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
export * from './localStorage'
export * from './trim_url'
export * from './highlight_substring'
export * from './codeblock_to_blob'

export function useFocusFirstFocusableChild(ref: MutableRefObject<HTMLDivElement | null>) {
  useLayoutEffect(() => {
    if (!ref.current) return
    const el = ref.current.querySelector('button, [href], input, select, textarea, [tabindex]')
    if (el != null) (el as HTMLElement).focus()
  }, [ref])
}

export function useWindowSize(): { width: number; height: number } {
  const [size, setSize] = useState({ width: window.innerWidth, height: window.innerHeight })

  const onResize = () => {
    setSize({ width: window.innerWidth, height: window.innerHeight })
  }

  useEffect(() => {
    window.addEventListener('resize', onResize)
    return () => {
      window.removeEventListener('resize', onResize)
    }
  }, [])

  return size
}

export function isOperationVariable(
  variable: GlobalVariableData | OperationVariableData,
): variable is OperationVariableData {
  return (variable as OperationVariableData).operationSlug !== undefined
}
