// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import { OIDCInstanceConfig } from '..'

// const loginWithOidc = () => {
//   window.location.href = "/web/auth/oidc/login"
// }

// export default (props: {
//   authFlags?: Array<string>
// }) => (
//   <div style={{ textAlign: 'right' }}>
//     <Button primary onClick={loginWithOidc}>Login With OIDC</Button>
//   </div>
// )

// use the below

const makeLoginFn = (code: string) => {
  return () => {
    window.location.href = `/web/auth/${code}/login`
  }
}

export const makeLogin = (config: OIDCInstanceConfig) => {
  const loginFn = makeLoginFn(config.code)

  return (_props: {
    authFlags?: Array<string>
  }) => (
    <div style={{ textAlign: 'right' }}>
      <Button primary onClick={loginFn}>Login With {config.name}</Button>
    </div>
  )
}
