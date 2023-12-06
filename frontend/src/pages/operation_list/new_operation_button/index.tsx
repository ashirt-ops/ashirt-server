import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onClick: () => void,
}) => (
  <button
    className={cx('root')}
    onClick={props.onClick}
  >
    <div className={cx('circle', 'plus')} />
    New Operation
  </button>
)
