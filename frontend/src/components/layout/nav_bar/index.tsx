// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import UserMenu from './user_menu'
import classnames from 'classnames/bind'
import {Link} from 'react-router-dom'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
}) => (
  <div className={cx('root')}>
    <Link to="/operations" className={cx('logo')}>
      <img src={require('./t-shirt.svg')} />
      ASHIRT
    </Link>

    <div className={cx('right')}>
      <UserMenu />
    </div>
  </div>
)
