// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

export const SpannedRow = (props: {
  span: number
  children?: React.ReactNode
}): React.ReactElement => (
  <tr>
    <td colSpan={props.span}>
      {props.children}
    </td>
  </tr >
)
