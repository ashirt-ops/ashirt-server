import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default function<T>(props: {
  disabled?: boolean,
  getLabel: (t: T) => string,
  groupLabel: string,
  onChange: (v: T) => void,
  options: Array<T>,
  value: T,
}) {
  return (
    <div className={cx({disabled: props.disabled})}>
      <div className={cx('group-label')}>{props.groupLabel}</div>
      {props.options.map((option, i) => (
        <label key={`${i}.${props.getLabel(option)}`} className={cx('option')}>
          {props.getLabel(option)}
          <input type="radio" disabled={props.disabled} checked={option === props.value} onChange={e => props.onChange(option)} />
          <span />
        </label>
      ))}
    </div>
  )
}
