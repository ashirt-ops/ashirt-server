// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import AuthContext from 'src/auth_context'
import ErrorDisplay from 'src/components/error_display'
import { NavLinkButton } from './components/button'
import { Route, Routes, Navigate, useParams, Params } from 'react-router-dom'
import { useAsyncComponent, useUserIsSuperAdmin } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

const AsyncLogin = makeAsyncPage(() => import('src/pages/login'))
const AsyncOperationList = makeAsyncPage(() => import('src/pages/operation_list'))
const AsyncOperationEdit = makeAsyncPage(() => import('src/pages/operation_edit'))
const AsyncFindingShow = makeAsyncPage(() => import('src/pages/operation_show/finding_show'))
const AsyncEvidenceList = makeAsyncPage(() => import('src/pages/operation_show/evidence_list'))
const AsyncFindingList = makeAsyncPage(() => import('src/pages/operation_show/finding_list'))
const AsyncAdminSettings = makeAsyncPage(() => import('src/pages/admin'))
const AsyncAccountSettings = makeAsyncPage(() => import('src/pages/account_settings'))
const AsyncNotFound = makeAsyncPage(() => import('src/pages/not_found'))

/**
 * Redirect provides a mechanism to redirect a user to the indicated URL
 * @param props.to [Required] The base url to go to
 * @param props.queryBuilder [Optional] A function to create the query string from the params
 * 
 * @returns A React-Router Navigate element to redirect the user to the indicated path
 */
function Redirect(props: {
  to: string,
  queryBuilder?: (params: Readonly<Params<string>>) => string
}) {
  const params = useParams()
  const query = props.queryBuilder?.(params)

  return <Navigate to={`${props.to}${query ? `?${query}` : ""}`} replace />
}

export default () => {
  const user = React.useContext(AuthContext).user
  const isSuperAdmin = useUserIsSuperAdmin()

  if (user == null) return (
    <Routes>
      <Route path="/login" element={<AsyncLogin />} />
      <Route path="/login/:schemeCode" element={<AsyncLogin />} />
      <Route path="/autherror/*" >
        <Route index element={<Redirect to="/login" />} />
        <Route path="recoveryfailed" element={<AuthRecoveryFailed />} />
        <Route path="noaccess" element={<AuthNoAccess />} />
        <Route path="noverify" element={<AuthNoVerify />} />
        <Route path="incomplete" element={<AuthIncomplete />} />
        <Route path="disabled" element={<AuthDisabled />} />
        <Route path="registrationdisabled" element={<AuthNoRegistration />} />
      </Route>
      <Route path="*" element={<Redirect to="/login" />} />
    </Routes>
  )

  return (
    <Routes>
      <Route path="/login" element={<Redirect to="/operations" />} />
      <Route path="/" element={<Redirect to="/operations" />} />

      <Route path="/operations/*">
        <Route index element={<AsyncOperationList />} />
        <Route path=":slug/*" >
          <Route index element={<Redirect to={`evidence`} />} />
          <Route path="evidence" element={<AsyncEvidenceList />} />
          <Route
            path="evidence/:uuid"
            element={
              <Redirect to={`../evidence`} queryBuilder={(params) => `q=uuid%3A${params.uuid}`}/>
            }
          />
          {/* ^^^ we need to do ../evidence because .. points to :slug, while . points to evidence/:uuid */}
          <Route path="findings" element={<AsyncFindingList />} />
          <Route path="findings/:uuid" element={<AsyncFindingShow />} />
          <Route path="edit/*">
            <Route index element={<Redirect to={`settings`} />} />
            <Route path="*" element={<AsyncOperationEdit />} />
          </Route>
        </Route>
      </Route>

      {/* Account Settings */}
      <Route path="/account/*" >
        <Route index element={<Redirect to="profile" />} />
        <Route path="*" element={<AsyncAccountSettings />} />
      </Route>

      {/* Admin Settings */}
      {isSuperAdmin && (
        <Route path="/admin/*" >
          <Route index element={<Redirect to="users" />} />
          <Route path="*" element={<AsyncAdminSettings />} />
        </Route>
      )}

      {/* AuthError routes that an admin might reach if testing */}
      <Route path="/autherror/recoveryfailed" element={<NoAccess />} />
      <Route path="*" element={<AsyncNotFound />} />
    </Routes >
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
