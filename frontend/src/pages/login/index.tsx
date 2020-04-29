// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { RouteComponentProps } from 'react-router-dom'
import { parse as parseQuery, ParsedUrlQuery } from 'querystring'
import { useAuthFrontendComponent } from 'src/authschemes'
import { useDataSource, getSupportedAuthentications } from 'src/services'
import { useWiredData } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

// This component renders a list of all enabled authscheme login components
// To add a new authentication method add a new authscheme frontend to
// src/auth and ensure it is enabled on the backend
// An optional schemeCode can be provided to only render that auth method
export default (props: RouteComponentProps<{ schemeCode?: string }>) => {
  const ds = useDataSource()
  const query = parseQuery(props.location.search.substr(1))
  const renderOnlyScheme = props.match.params.schemeCode
  const wiredAuthSchemes = useWiredData(React.useCallback(() => (
    getSupportedAuthentications(ds)
  ), [ds]))

  return wiredAuthSchemes.render(supportedAuthSchemes => (
    <div className={cx('login')}>
      {supportedAuthSchemes.map(({schemeCode}) => {
        if (renderOnlyScheme != null && schemeCode != renderOnlyScheme) return null
        return (
          <AuthSchemeLogin
            key={schemeCode}
            authSchemeCode={schemeCode}
            query={query}
          />
        )
      })}
    </div>
  ))
}


const AuthSchemeLogin = (props: {
  authSchemeCode: string,
  query: ParsedUrlQuery,
}) => {
  const Login = useAuthFrontendComponent(props.authSchemeCode, 'Login')
  return (
    <div className={cx('auth-scheme-row')}>
      <Login query={props.query} />
    </div>
  )
}
