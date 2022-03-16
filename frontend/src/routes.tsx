// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import AuthContext from 'src/auth_context'
import ErrorDisplay from 'src/components/error_display'
import { NavLinkButton } from './components/button'
import { Route, Switch, Redirect, RouteComponentProps } from 'react-router-dom'
import { useAsyncComponent, useUserIsSuperAdmin } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

const AsyncLogin = makeAsyncPage(() => import('src/pages/login'))
const AsyncOperationList = makeAsyncPage(() => import('src/pages/operation_list'))
const AsyncOperationEdit = makeAsyncPage(() => import('src/pages/operation_edit'))
const AsyncOperationOverview = makeAsyncPage(() => import('src/pages/operation_overview'))
const AsyncFindingShow = makeAsyncPage(() => import('src/pages/operation_show/finding_show'))
const AsyncEvidenceList = makeAsyncPage(() => import('src/pages/operation_show/evidence_list'))
const AsyncFindingList = makeAsyncPage(() => import('src/pages/operation_show/finding_list'))
const AsyncAdminSettings = makeAsyncPage(() => import('src/pages/admin'))
const AsyncAccountSettings = makeAsyncPage(() => import('src/pages/account_settings'))
const AsyncNotFound = makeAsyncPage(() => import('src/pages/not_found'))

export default () => {
  const user = React.useContext(AuthContext).user
  const isSuperAdmin = useUserIsSuperAdmin()

  if (user == null) return (
    <Switch>
      <Route exact path="/login" >
        {(props: RouteComponentProps) => <AsyncLogin {...props} />}
      </Route>
      <Route exact path="/login/:schemeCode">
        {(props: RouteComponentProps) => <AsyncLogin {...props} />}
      </Route>

      <Route exact path="/autherror/recoveryfailed">
        <AuthRecoveryFailed />
      </Route>
      <Route exact path="/autherror/noaccess">
        <AuthNoAccess />
      </Route>
      <Route exact path="/autherror/noverify">
        <AuthNoVerify />
      </Route>
      <Route exact path="/autherror/incomplete">
        <AuthIncomplete />
      </Route>
      <Route exact path="/autherror/disabled">
        <AuthDisabled />
      </Route>
      <Route exact path="/autherror/registrationdisabled">
        <AuthNoRegistration />
      </Route>

      <Route render={() => <Redirect to="/login" />} />
    </Switch>
  )

  return (
    <Switch>
      <Route exact path="/login" render={() => <Redirect to="/operations" />} />
      <Route exact path="/" render={() => <Redirect to="/operations" />} />

      {/* AuthError routes that an admin might reach if testing */}
      <Route exact path="/autherror/recoveryfailed" render={NoAccess} />

      <Route exact path="/operations" component={AsyncOperationList} />

      {/* Operation edit */}
      <Route exact path="/operations/:slug/edit/:view(settings|users|tags)" component={AsyncOperationEdit} />
      <Redirect from="/operations/:slug/edit" to="/operations/:slug/edit/settings" />

      {/* Operation overview */}
      <Route exact path="/operations/:slug/overview" component={AsyncOperationOverview} />

      {/* Operation show */}
      <Route exact path="/operations/:slug/findings" component={AsyncFindingList} />
      <Route exact path="/operations/:slug/findings/:uuid" component={AsyncFindingShow} />
      <Route exact path="/operations/:slug/evidence" component={AsyncEvidenceList} />
      <Redirect exact path="/operations/:slug/evidence/:uuid" to="/operations/:slug/evidence?q=uuid%3A:uuid" />
      <Redirect from="/operations/:slug" to="/operations/:slug/evidence" />

      {/* Account Settings */}
      <Route exact path="/account/:view(profile|security|apikeys|authmethods)" component={AsyncAccountSettings} />
      <Redirect exact from="/account" to={`/account/profile`} />

      {isSuperAdmin && <Route exact path="/account/:view(profile|apikeys|authmethods)/:slug" component={AsyncAccountSettings} />}
      {isSuperAdmin && <Redirect exact from="/account/edit/:slug" to="/account/profile/:slug" />}

      {/* Admin Settings */}
      {isSuperAdmin && <Route exact path="/admin/:view(users|operations|authdata|findings|tags)" component={AsyncAdminSettings} />}
      {isSuperAdmin && <Redirect from="/admin/" to="/admin/users" />}


      <Route component={AsyncNotFound} />
    </Switch>
  )
}

// makeAsyncPage turns a `() => import('path/to/page/component')` into a react component ready to be passed into <Route />
// It uses useAsyncComponent to properly render a loading spinner or error as appropriate
//
// This is used to break up each page into its own bundle to prevent the main entry bundle from becoming too large and allows
// page javascript to load on demand.
function makeAsyncPage(page: () => Promise<{ default: React.FunctionComponent<RouteComponentProps> }>) {
  const defaultPage = () => page().then(module => module.default)
  return (props: RouteComponentProps) => {
    const Page = useAsyncComponent(defaultPage);
    return <Page {...props} />
  }
}

const makeErrorDisplay = (title: string, message: string, withLoginLink = false) => () => (
  <ErrorDisplay title={title} err={new Error(message)}>
    {
      withLoginLink && (
        <>
          <br />
          <NavLinkButton primary className={cx('return-button')} to={"/login"}>Return to login</NavLinkButton>
        </>
      )
    }
  </ErrorDisplay>
)
const makeAuthErr = (body: string, addLoginLink?: boolean) => makeErrorDisplay("Authentication Error", body, addLoginLink ?? true)

const AuthRecoveryFailed = makeAuthErr("Account recovery failed. The recovery code may be expired or incorrect. Please contact an administrator to provide a new url.")
const AuthNoAccess = makeAuthErr("This user is not permitted to use this service", false)
const AuthNoVerify = makeAuthErr("Unable to verify user account. Please try again.")
const AuthIncomplete = makeAuthErr("The system could not complete the login process. Please retry, and if the issue persists, please contact a system administrator.")
const AuthDisabled = makeAuthErr("This account has been disabled. Please contact an adminstrator if you think this is an error", false)
const AuthNoRegistration = makeAuthErr("Registration has been disabled. Please contract an administrator to request access.")
const NoAccess = makeErrorDisplay("Access Error", "This url only works for users that are not logged in.", true)
