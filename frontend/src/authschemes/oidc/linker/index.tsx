import { type MouseEvent } from 'react'

import Button from 'src/components/button'
import { type OIDCInstanceConfig } from '..'

// export default (props: {
//   onSuccess: () => void,
//   authFlags?: Array<string>
// }) => (
//   <Button primary onClick={(e) => { e.preventDefault(); window.location.href = "/web/auth/oidc/link" }}>Login with OIDC</Button >
// )

export const makeLinker = (config: OIDCInstanceConfig) => {
  const onClick = (e: MouseEvent<Element>) => {
    e.preventDefault()
    window.location.href = `/web/auth/${config.code}/link`
  }

  return (_props: { onSuccess: () => void; authFlags?: Array<string> }) => (
    <Button primary onClick={onClick}>
      Login with {config.name}
    </Button>
  )
}
