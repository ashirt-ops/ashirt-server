// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { useLocation, useParams } from 'react-router-dom'

import { SupportedAuthenticationScheme } from 'src/global_types'
import { getSupportedAuthentications } from 'src/services/auth'
import { useAuthFrontendComponent } from 'src/authschemes'
import { useWiredData } from 'src/helpers'
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
    if (supportedAuthSchemes.length === 0) {
      return (
        <div className={cx('login')}>
          <NoAuthsWarning />
        </div>
      )
    }

    if (supportedAuthSchemes.length === 1) {
      const scheme = supportedAuthSchemes[0]
      return (
        <div className={cx('login')}>
          <AuthSchemeLogin
            key={scheme.schemeCode}
            authSchemeType={scheme.schemeType}
            authScheme={scheme}
            query={query}
          />
        </div>
      )
    }

    return (
      <div className={cx('login')}>
        <LoginMenu selectedAuth={currentAuth} resetAuth={() => setCurrentAuth(null)} >
          {currentAuth == null
            ? <AuthButtons schemes={supportedAuthSchemes} onSelected={setCurrentAuth} renderOnlyScheme={renderOnlyScheme} />
            : (
              <AuthSchemeLogin
                key={currentAuth.schemeCode}
                authSchemeType={currentAuth.schemeType}
                authScheme={currentAuth}
                query={query}
              />
            )
          }
        </LoginMenu>
      </div >
    )
  })
}

const LoginMenu = (props: {
  selectedAuth: SupportedAuthenticationScheme | null
  resetAuth: () => void
  children: React.ReactNode
}) => {
  return (
    <div className={cx('login-wrapper')}>
      {props.selectedAuth === null
        ? (
          <h1 className={cx('login-wrapper-title')}>How do you want to authenticate?</h1>
        )
        : (<>
          <h1 className={cx('login-wrapper-title')}>{props.selectedAuth.schemeName}</h1>
          <div>
            <Button
              className={cx('full-width-button')}
              icon={require('./back.svg')}
              onClick={props.resetAuth}
            >
              Choose different method
            </Button>
          </div>
        </>)

      }
      <hr className={cx('login-wrapper-divider')} />
      <div className={cx('login-wrapper-children')}>
        {props.children}
      </div>
    </div>
  )
}

const NoAuthsWarning = (props: {}) => (
  <div className={cx('no-auths-warning')}>This instance of AShirt has no way to authenticate users.</div>
)

const AuthButtons = (props: {
  schemes: Array<SupportedAuthenticationScheme>
  onSelected: (scheme: SupportedAuthenticationScheme | null) => void
  renderOnlyScheme?: string
}) => {
  return (
    <div className={cx('auth-buttons')}>
      {
        props.schemes.map(schemeDetails => {
          const { schemeCode, schemeName } = schemeDetails
          if (props.renderOnlyScheme != null && schemeCode != props.renderOnlyScheme) {
            return null
          }

          return (
            // todo: fix direction of arrow
            <Button
              className={cx('full-width-button')}
              key={schemeCode}
              afterIcon={require('./forward.svg')}
              onClick={() => props.onSelected(schemeDetails)}
            >
              {schemeName}
            </Button>
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
