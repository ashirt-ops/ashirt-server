import * as React from 'react'

import BulletChooser, { BulletProps } from 'src/components/bullet_chooser'
import { SupportedEvidenceType } from 'src/global_types'

export type EvidenceTypeOption = BulletProps & {
  id: SupportedEvidenceType
}

export const supportedEvidenceCount: Array<EvidenceTypeOption> = [
  { name: 'Screenshot', id: 'image' },
  { name: 'Code Block', id: 'codeblock' },
  { name: 'Terminal Recording', id: 'terminal-recording' },
  { name: 'HTTP Request/Response', id: 'http-request-cycle' },
  { name: 'Events', id: 'event' },
  { name: 'No Content', id: 'none' },
]

export const EvidenceTypeChooser = (props: {
  label: string
  value: Array<BulletProps>
  onChange: (types: Array<BulletProps>) => void
  className?: string
  disabled?: boolean
  enableNot?: boolean
}) => {
  return (
    <BulletChooser
      label={props.label}
      options={supportedEvidenceCount}
      value={props.value}
      onChange={props.onChange}
      enableNot={props.enableNot}
    />
  )
}

export default EvidenceTypeChooser
