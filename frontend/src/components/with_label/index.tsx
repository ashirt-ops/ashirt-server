import { type ReactNode } from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default function WithLabel(props: { children?: ReactNode; className?: string; label?: string }) {
  if (!props.label) {
    if (!props.children) return null
    return <div className={props.className}>{props.children}</div>
  }

  return (
    <label className={cx('root', props.className)}>
      <div className={cx('label-text')}>{props.label}</div>
      {props.children}
    </label>
  )
}
