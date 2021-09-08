// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import Button from 'src/components/button'
import { OIDCInstanceConfig } from '..';

// export default (props: {
//   onSuccess: () => void,
//   authFlags?: Array<string>
// }) => (
//   <Button primary onClick={(e) => { e.preventDefault(); window.location.href = "/web/auth/oidc/link" }}>Login with OIDC</Button >
// )

export const makeLinker = (config: OIDCInstanceConfig) => {
  const onClick = (e: React.MouseEvent<Element, MouseEvent>) => {
    e.preventDefault()
    window.location.href = `/web/auth/${config.code}/link`
  }

  return (_props: {
    onSuccess: () => void,
    authFlags?: Array<string>
  }) => (
    <Button primary onClick={onClick}>
      Login with {config.name}
    </Button >
  )
}
