import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  small?: boolean
}) => (
  <div className={cx('root', props.className, {small: props.small})}>
    <div className={cx('spinner')} />
  </div>
)
