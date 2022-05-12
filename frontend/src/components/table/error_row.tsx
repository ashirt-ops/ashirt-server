// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

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
