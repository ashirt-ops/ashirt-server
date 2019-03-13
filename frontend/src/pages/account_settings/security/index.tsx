// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import {useAuthFrontendComponent} from 'src/authschemes'

export default (props: {
}) => {
  const {user} = React.useContext(AuthContext)
  if (user == null) return null

  return <>
    {user.authSchemes.map(authScheme => (
      <AuthSchemeSettings
        key={authScheme.schemeCode}
        authSchemeCode={authScheme.schemeCode}
        userKey={authScheme.userKey}
      />
    ))}
  </>
}

const AuthSchemeSettings = (props: {
  authSchemeCode: string,
  userKey: string,
}) => {
  const Settings = useAuthFrontendComponent(props.authSchemeCode, 'Settings')
  return (
    <Settings userKey={props.userKey} />
  )
}
