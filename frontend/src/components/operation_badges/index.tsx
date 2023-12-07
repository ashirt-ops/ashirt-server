import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  numUsers: number,
  showDetailsModal: () => void,
}) => (
  <button className={cx('root', props.className)} onClick={() => props.showDetailsModal()}>
    <div
      className={cx('icon', 'users')}
      title={`${props.numUsers} user${props.numUsers === 1 ? ' belongs' : 's belong'} to this operation`}
      children={props.numUsers}
    />
  </button>
)
