import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const ErrorDisplay = (props: {
  title?: string,
  err?: Error,
  children?: React.ReactNode,
}) => (
  <div className={cx('root')}>
    <img src={require('./icon.svg')} />
    <div className={cx('message')}>
      {props.title || 'An Error Occurred'}:
    </div>
    <div className={cx('message')}>
      {props.err && props.err.message}
      {props.children}
    </div>
  </div>
)

export default ErrorDisplay
