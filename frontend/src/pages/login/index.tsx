// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import {parse as parseQuery, ParsedUrlQuery} from 'querystring'
import {RouteComponentProps} from 'react-router-dom'
import {getSupportedAuthentications} from 'src/services/auth'
import {useAuthFrontendComponent} from 'src/authschemes'
import {useWiredData} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

// This component renders a list of all enabled authscheme login components
// To add a new authentication method add a new authscheme frontend to
// src/auth and ensure it is enabled on the backend
// An optional schemeCode can be provided to only render that auth method
export default (props: RouteComponentProps<{ schemeCode?: string }>) => {
  const query = parseQuery(props.location.search.substr(1))
  const renderOnlyScheme = props.match.params.schemeCode
  const wiredAuthSchemes = useWiredData(getSupportedAuthentications)

  return wiredAuthSchemes.render(supportedAuthSchemes => (
    <div className={cx('login')}>
      {supportedAuthSchemes.map(({schemeCode, schemeFlags}) => {
        if (renderOnlyScheme != null && schemeCode != renderOnlyScheme) return null
        return (
          <AuthSchemeLogin
            key={schemeCode}
            authSchemeCode={schemeCode}
            query={query}
            authSchemeFlags={schemeFlags}
          />
        )
      })}
    </div>
  ))
}


const AuthSchemeLogin = (props: {
  authSchemeCode: string,
  authSchemeFlags: Array<string>
  query: ParsedUrlQuery,
}) => {
  const Login = useAuthFrontendComponent(props.authSchemeCode, 'Login')
  return (
    <div className={cx('auth-scheme-row')}>
      <Login query={props.query} authFlags={props.authSchemeFlags}/>
    </div>
  )
}
