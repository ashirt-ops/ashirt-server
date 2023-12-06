import * as React from 'react'
import ErrorDisplay from 'src/components/error_display'
import { SpannedRow } from './spanned_row'


export const ErrorRow = (props: {
  span: number
  error: Error
}) => (
  <SpannedRow span={props.span}>
    <ErrorDisplay err={props.error} />
  </SpannedRow>
)
