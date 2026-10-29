import { type ReactNode, type ReactElement } from 'react'

export const SpannedRow = (props: { span: number; children?: ReactNode }): ReactElement => (
  <tr>
    <td colSpan={props.span}>{props.children}</td>
  </tr>
)
