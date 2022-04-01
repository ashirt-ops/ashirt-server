// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import AuthContext from 'src/auth_context'
import ErrorDisplay from 'src/components/error_display'
import { NavLinkButton } from './components/button'
import { Route, Routes, Navigate, useParams, useRoutes, generatePath } from 'react-router-dom'
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

function Redirect(props: {
  to: string
}) {
  const params = useParams()
  return <Navigate to={generatePath(props.to, params)} replace />
}

export default () => {
  const user = React.useContext(AuthContext).user
  const isSuperAdmin = useUserIsSuperAdmin()

  if (user == null) return (
    <Routes>
      <Route path="/login" element={<AsyncLogin />} />
      <Route path="/login/:schemeCode" element={<AsyncLogin />} />
      <Route path="/autherror/recoveryfailed" element={<AuthRecoveryFailed />} />
      <Route path="/autherror/noaccess" element={<AuthNoAccess />} />
      <Route path="/autherror/noverify" element={<AuthNoVerify />} />
      <Route path="/autherror/incomplete" element={<AuthIncomplete />} />
      <Route path="/autherror/disabled" element={<AuthDisabled />} />
      <Route path="/autherror/registrationdisabled" element={<AuthNoRegistration />} />

      <Route path="*" element={<Redirect to="/login" />} />
    </Routes>
  )

  return (
    <Routes>
      <Route path="/login" element={<Redirect to="/operations" />} />
      <Route path="/" element={<Redirect to="/operations" />} />

      {/* AuthError routes that an admin might reach if testing */}
      <Route path="/autherror/recoveryfailed" element={<NoAccess />} />

      <Route path="/operations" element={<AsyncOperationList />} />

      {/* Operation edit */}
      <Route path="/operations/:slug/edit/:view" element={<AsyncOperationEdit />} />
      <Route
        path="/operations/:slug/edit"
        element={<Redirect to={`/operations/:slug/edit/settings`}/>}
      />

      {/* Operation overview */}
      <Route path="/operations/:slug/overview" element={<AsyncOperationOverview />} />

      {/* Operation show */}
      <Route path="/operations/:slug/findings" element={<AsyncFindingList />} />
      <Route path="/operations/:slug/findings/:uuid" element={<AsyncFindingShow />} />
      <Route path="/operations/:slug/evidence" element={<AsyncEvidenceList />} />

      <Route path="/operations/:slug/evidence/:uuid" element={
        <Redirect to={`/operations/:slug/evidence?q=uuid%3A:uuid`} />
      } />
      <Route path="/operations/:slug" element={<Redirect to={`/operations/:slug/evidence`}/>} />

      {/* Account Settings */}
      <Route path="/account/:view" element={<AsyncAccountSettings />} />
      <Route path="/account" element={<Redirect to="/account/profile" />} />

      {isSuperAdmin && (
        // For some reason, we can't navigate to this route directly -- only through page links
        <Route path="/account/:view/:slug" element={<AsyncAccountSettings />} />
      )}
      {isSuperAdmin && (
        <Route path="/account/edit/:slug" element={<Redirect to={`/account/profile/:slug`} />} />
      )}

      {/* Admin Settings */}
      {isSuperAdmin && (
        <Route path="/admin/:view" element={<AsyncAdminSettings />} />
      )}
      {isSuperAdmin && (
        <Route path="/admin/*" element={<Redirect to="/admin/users" />} />
      )}

      <Route path="*" element={<AsyncNotFound />} />
    </Routes>
  )
}

// makeAsyncPage turns a `() => import('path/to/page/component')` into a react component ready to be passed into <Route />
// It uses useAsyncComponent to properly render a loading spinner or error as appropriate
//
// This is used to break up each page into its own bundle to prevent the main entry bundle from becoming too large and allows
// page javascript to load on demand.
function makeAsyncPage(page: () => Promise<{ default: React.FunctionComponent }>) {
  const defaultPage = () => page().then(module => module.default)
  return () => {
    const Page = useAsyncComponent(defaultPage);
    return <Page />
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
