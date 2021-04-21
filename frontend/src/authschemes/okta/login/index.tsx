// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'

const loginWithOkta = () => {
  window.location.href = "/web/auth/okta/login"
}

export default (props: {
  authFlags?: Array<string>
}) => (
  <div style={{textAlign: 'right'}}>
    <Button primary onClick={loginWithOkta}>Login With Okta</Button>
  </div>
)
