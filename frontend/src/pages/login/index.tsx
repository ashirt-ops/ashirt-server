// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useLocation, useParams } from 'react-router-dom'

import { useAuthFrontendComponent } from 'src/authschemes'
import Button from 'src/components/button'
import { SupportedAuthenticationScheme } from 'src/global_types'
import { useWiredData } from 'src/helpers'
import { getSupportedAuthentications } from 'src/services/auth'

import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

// This component renders a list of all enabled authscheme login components
// To add a new authentication method add a new authscheme frontend to
// src/auth and ensure it is enabled on the backend
// An optional schemeCode can be provided to only render that auth method
export default () => {
  const { schemeCode: renderOnlyScheme } = useParams<{ schemeCode?: string }>()
  const location = useLocation()
  const wiredAuthSchemes = useWiredData(getSupportedAuthentications)

  const query = new URLSearchParams(location.search)

  return wiredAuthSchemes.render(supportedAuthSchemes => {
    if (supportedAuthSchemes.length === 0) {
      return (
        <div className={cx('login')}>
          <NoAuthsWarning />
        </div>
      )
    }

    const oneScheme = (supportedAuthSchemes.length === 1)
      ? supportedAuthSchemes[0]
      : (renderOnlyScheme !== undefined)
        ? supportedAuthSchemes.find(s => s.schemeCode == renderOnlyScheme)
        : null

    return (
      <div className={cx('login')}>
        {oneScheme
          ? <AuthSchemeLogin key={oneScheme.schemeCode} authScheme={oneScheme} query={query} />
          : <LoginMenu authSchemes={supportedAuthSchemes} query={query} />
        }
      </div>
    )
  })
}

const LoginMenu = (props: {
  authSchemes: Array<SupportedAuthenticationScheme>
  query: URLSearchParams
}) => {
  const [currentAuth, setCurrentAuth] = React.useState<SupportedAuthenticationScheme | null>(null)

  return (
    <div className={cx('login-menu')}>
      {currentAuth == null
        ? <MenuHeader title="How do you want to authenticate?" />
        : (
          <MenuHeader
            title={currentAuth.schemeName}
            backButtonText="Choose a different method"
            onBackPressed={() => setCurrentAuth(null)}
          />
        )
      }

      <hr className={cx('menu-divider')} />
      {currentAuth == null
        ? <AuthButtons schemes={props.authSchemes} onSelected={setCurrentAuth} />
        : <AuthSchemeLogin authScheme={currentAuth} query={props.query} />
      }
    </div>
  )
}

const MenuHeader = (props: {
  title: string
  backButtonText?: string
  onBackPressed?: () => void
}) => (
  <>
    <h1 className={cx('menu-title')}>{props.title}</h1>
    {props.onBackPressed && (
      <Button
        className={cx('full-width-button')}
        icon={require('./back.svg')}
        onClick={props.onBackPressed}
      >
        {props.backButtonText}
      </Button>
    )}
  </>
)

const NoAuthsWarning = (props: {}) => (
  <div className={cx('no-auths-warning')}>This instance of AShirt has no way to authenticate users.</div>
)

const AuthButtons = (props: {
  schemes: Array<SupportedAuthenticationScheme>
  onSelected: (scheme: SupportedAuthenticationScheme | null) => void
}) => (
  <div className={cx('auth-buttons')}>
    {
      props.schemes.map(schemeDetails => {
        const { schemeCode, schemeName } = schemeDetails

        return (
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

const AuthSchemeLogin = (props: {
  authScheme: SupportedAuthenticationScheme,
  query: URLSearchParams,
}) => {
  const Login = useAuthFrontendComponent(props.authScheme.schemeType, 'Login', props.authScheme)
  return (
    <div className={cx('auth-scheme-row')}>
      <Login query={props.query} authFlags={props.authScheme.schemeFlags} />
    </div>
  )
}
