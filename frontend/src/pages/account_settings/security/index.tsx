// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useAuthFrontendComponent } from 'src/authschemes'
import { SupportedAuthenticationScheme, UserOwnView } from 'src/global_types'

export default (props: {
  user: UserOwnView | null
}) => {
  const { user } = props
  if (user == null) return null

  return <>
    {user.authSchemes.map(authScheme => (
      <AuthSchemeSettings
        key={authScheme.schemeCode}
        authSchemeDetails={authScheme.authDetails}
        authSchemeType={authScheme.schemeType}
        userKey={authScheme.userKey}
      />
    ))}
  </>
}

const AuthSchemeSettings = (props: {
  authSchemeDetails?: SupportedAuthenticationScheme
  authSchemeType: string,
  userKey: string,
}) => {
  const Settings = useAuthFrontendComponent(props.authSchemeType, 'Settings', props.authSchemeDetails)
  return (
    <Settings userKey={props.userKey} authFlags={props.authSchemeDetails?.schemeFlags || []} />
  )
}
