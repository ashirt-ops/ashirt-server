// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import classnames from 'classnames/bind'
import {ClickPopover} from 'src/components/popover'
import {default as Menu, MenuItem, MenuSeparator} from 'src/components/menu'
import {logout} from 'src/services'
import {useUserIsSuperAdmin} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

// @ts-ignore - npm package @types/react-router-dom needs to be updated (https://github.com/DefinitelyTyped/DefinitelyTyped/issues/40131)
import {useHistory} from 'react-router-dom'

const Avatar = (props: {
  url: string,
}) => (
  <div className={cx('avatar')} style={props.url ? {backgroundImage: `url(${props.url})`} : {}} />
)

const logoutAndRedirect = () => (
  logout().then(() => {
    window.location.pathname = '/'
  })
)

const UserMenuDropdown = (props: {
  name: string,
  avatar: string,
}) => {
  const history = useHistory()
  const isSuperAdmin = useUserIsSuperAdmin()

  return (
    <Menu>
      <div className={cx('account-menu-top')} onClick={e => e.stopPropagation()}>
        <Avatar url={props.avatar} />
        {props.name}
      </div>
      <MenuSeparator />
      <MenuItem icon={require('./settings.svg')} onClick={() => history.push('/account/profile')}>Account Settings</MenuItem>
      {isSuperAdmin && <MenuItem icon={require('./admin.svg')} onClick={() => history.push('/admin')}>Admin</MenuItem>}
      <MenuItem icon={require('./power.svg')} onClick={logoutAndRedirect}>Sign Out</MenuItem>
    </Menu>
  )
}

export default (props: {
}) => {
  const {user} = React.useContext(AuthContext)
  if (user == null) return null
  const name = `${user.firstName} ${user.lastName}`
  const avatar = ''

  return (
    <ClickPopover closeOnContentClick content={<UserMenuDropdown name={name} avatar={avatar} />}>
      <a className={cx('user-menu')}>
        <Avatar url={avatar} />
      </a>
    </ClickPopover>
  )
}
