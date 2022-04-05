// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { Route, Routes, useLocation } from 'react-router-dom'
import { UserOwnView } from 'src/global_types'
import { getUser } from 'src/services'
import { useUserIsSuperAdmin, useWiredData } from 'src/helpers'

import ApiKeys from './api_keys'
import AuthMethods from './auth_methods'
import { NavVerticalTabMenu } from 'src/components/tab_vertical_menu'
import Profile from './profile'
import Security from './security'

const cx = classnames.bind(require('./stylesheet'))

export const AccountSettings = () => {
  const location = useLocation()
  const isSuperAdmin = useUserIsSuperAdmin()
  const user = new URLSearchParams(location.search).get('user')

  // if a non-admin user tries to access someone's profile (not themselves), show the caller their
  // own profile. The backend prevents the user's data from displaying, but this provides a more
  // sane way of locking out a user
  const userSlug = (isSuperAdmin && user != null)
    ? user
    : ""

  const wiredProfile = useWiredData<UserOwnView>(
    React.useCallback(() => getUser({ userSlug }), [userSlug])
  )

  const bus = BuildReloadBus()

  React.useEffect(() => {
    bus.onReload(wiredProfile.reload)
    return () => { bus.offReload(wiredProfile.reload) }
  })

  return <>
    <div className={cx('root')}>
      {wiredProfile.render(p => {
        const query = user !== null ? { user } : undefined
        const tabs = [
          { id: "profile", label: "Profile", query },
          { id: "authmethods", label: "Authentication Methods", query },
          ...(userSlug
            ? []
            : [{ id: "security", label: "Security" }]
          ),
          { id: "apikeys", label: "API Keys", query },
        ]

        return (<>
          <NavVerticalTabMenu title="Account Settings" tabs={tabs}>
            {userSlug &&
              <em className={cx('notice')}>Editing the profile of: <em className={cx('editing-user-name')}>{`${p.firstName} ${p.lastName}`}</em></em>
            }
            <Routes>
              <Route path="profile" element={<Profile {...bus} profile={p} />} />
              <Route path="authmethods" element={<AuthMethods {...bus} profile={p} allowLinking={!userSlug} />} />
              {
                !userSlug && (<Route path="security" element={<Security />} />)
              }
              <Route path="apikeys" element={<ApiKeys profile={p} />} />
            </Routes>
          </NavVerticalTabMenu>
        </>)
      })}
    </div>
  </>
}
export default AccountSettings
