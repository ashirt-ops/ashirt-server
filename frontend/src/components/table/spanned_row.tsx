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
