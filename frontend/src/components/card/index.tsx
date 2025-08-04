import * as React from 'react'
import classnames from 'classnames/bind'
import {Link} from 'react-router'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  linkTo?: string,
  children: React.ReactNode,
}) => (
  props.linkTo == null ? (
    <div className={cx('root', props.className)}>
      {props.children}
    </div>
  ) : (
    <Link to={props.linkTo} className={cx('root', props.className)}>
      {props.children}
    </Link>
  )
)
