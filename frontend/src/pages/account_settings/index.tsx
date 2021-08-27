// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { RouteComponentProps } from 'react-router-dom'
import { UserOwnView } from 'src/global_types'
import { getUser } from 'src/services'
import { useWiredData } from 'src/helpers'

import ApiKeys from './api_keys'
import AuthMethods from './auth_methods'
import { NavVerticalTabMenu } from 'src/components/tab_vertical_menu'
import Profile from './profile'
import Security from './security'

const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{ slug: string }>) => {
  const forUser = props.match.params.slug
  const wiredProfile = useWiredData<UserOwnView>(React.useCallback(() => getUser({ userSlug: forUser }), [forUser]))

  const bus = BuildReloadBus()

  React.useEffect(() => {
    bus.onReload(wiredProfile.reload)
    return () => { bus.offReload(wiredProfile.reload) }
  })

  return <>
    <div className={cx('root')}>
      {wiredProfile.render(p => {
        const tabs = [
          { id: "profile", label: "Profile", content: <Profile {...bus} profile={p} /> },
          { id: "authmethods", label: "Authentication Methods", content: <AuthMethods {...bus} profile={p} allowLinking={!forUser}/> }
        ]
        if (!forUser) {
          tabs.push({ id: "security", label: "Security", content: <Security /> })
        }
        tabs.push({ id: "apikeys", label: "API Keys", content: <ApiKeys profile={p} /> })

        return (<>
          {forUser &&
            <em className={cx('notice')}>Editing the profile of: <em className={cx('editing-user-name')}>{`${p.firstName} ${p.lastName}`}</em></em>
          }
          <NavVerticalTabMenu title="Account Settings" {...props} tabs={tabs} />
        </>)
      })}
    </div>
  </>
}
