// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useLocation, useParams } from 'react-router-dom'
import { getSupportedAuthentications } from 'src/services/auth'
import { useAuthFrontendComponent } from 'src/authschemes'
import { useWiredData } from 'src/helpers'
import { SupportedAuthenticationScheme } from 'src/global_types'
import Button from 'src/components/button'
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
  const [currentAuth, setCurrentAuth] = React.useState<SupportedAuthenticationScheme | null>(null)

  return wiredAuthSchemes.render(supportedAuthSchemes => {
    return (
      <div className={cx('login')}>
        {currentAuth == null
          ? <AuthButtons schemes={supportedAuthSchemes} onSelected={setCurrentAuth} renderOnlyScheme={renderOnlyScheme} />
          : (
            <div>
              <AuthSchemeLogin
                key={currentAuth.schemeCode}
                authSchemeType={currentAuth.schemeType}
                authScheme={currentAuth}
                query={query}
              />
              <Button
                icon={require('./back.svg')}
                onClick={() => setCurrentAuth(null)}
              >
                Choose different method
              </Button>
            </div>
          )
        }
      </div>
    )
  })
}

const AuthButtons = (props: {
  schemes: Array<SupportedAuthenticationScheme>
  onSelected: (scheme: SupportedAuthenticationScheme | null) => void
  renderOnlyScheme?: string
}) => {
  return (
    <div>
      {
        props.schemes.map(schemeDetails => {
          const { schemeCode, schemeName } = schemeDetails
          if (props.renderOnlyScheme != null && schemeCode != props.renderOnlyScheme) {
            return null
          }

          return (
            <Button onClick={() => props.onSelected(schemeDetails)} >{schemeName}</Button>
          )
        })
      }
    </div>

  )
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
