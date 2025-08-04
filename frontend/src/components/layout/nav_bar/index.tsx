import * as React from 'react'
import UserMenu from './user_menu'
import classnames from 'classnames/bind'
import {Link} from 'react-router'
const cx = classnames.bind(require('./stylesheet'))

export default () => (
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
