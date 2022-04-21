// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import LoadingSpinner from 'src/components/loading_spinner'
import { SpannedRow } from './spanned_row'


export const LoadingRow = (props: {
  span: number
}) => (
  <SpannedRow span={props.span}>
    <LoadingSpinner />
  </SpannedRow>
)
