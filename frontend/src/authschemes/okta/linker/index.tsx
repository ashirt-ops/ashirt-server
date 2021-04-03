// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import Button from 'src/components/button'

export default (props: {
  onSuccess: () => void,
  authFlags?: Array<string>
}) => (
  <Button primary onClick={(e) => { e.preventDefault(); window.location.href = "/web/auth/okta/link" }}>Login with Okta</Button >
)
