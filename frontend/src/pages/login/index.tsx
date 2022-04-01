// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useLocation, useParams } from 'react-router-dom'
import { getSupportedAuthentications } from 'src/services/auth'
import { useAuthFrontendComponent } from 'src/authschemes'
import { useWiredData } from 'src/helpers'
import { SupportedAuthenticationScheme } from 'src/global_types'
const cx = classnames.bind(require('./stylesheet'))

// This component renders a list of all enabled authscheme login components
// To add a new authentication method add a new authscheme frontend to
// src/auth and ensure it is enabled on the backend
// An optional schemeCode can be provided to only render that auth method
export default () => {
  const { schemeCode: renderOnlyScheme } = useParams<{ schemeCode?: string }>()
  const location = useLocation()
  const query = new URLSearchParams(location.search)
  const wiredAuthSchemes = useWiredData(getSupportedAuthentications)

  return wiredAuthSchemes.render(supportedAuthSchemes => (
    <div className={cx('login')}>
      {supportedAuthSchemes.map((schemeDetails) => {
        const { schemeCode, schemeType } = schemeDetails
        if (renderOnlyScheme != null && schemeCode != renderOnlyScheme) return null
        return (
          <AuthSchemeLogin
            key={schemeCode}
            authSchemeType={schemeType}
            authScheme={schemeDetails}
            query={query}
          />
        )
      })}
    </div>
  ))
}


const AuthSchemeLogin = (props: {
  authSchemeType: string,
  authScheme: SupportedAuthenticationScheme,
  query: URLSearchParams,
}) => {
  const Login = useAuthFrontendComponent(props.authSchemeType, 'Login', props.authScheme)
  return (
    <div className={cx('auth-scheme-row')}>
      <Login query={props.query} authFlags={props.authScheme.schemeFlags} />
    </div>
  )
}
