// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  value?: boolean,
  label?: string,
  title?: string,
  disabled?: boolean,
  onChange?: (value: boolean) => void,
}) => (
  <label title={props.title} className={cx('root', props.className)} onClick={e => e.stopPropagation()}>
    <input
      type="checkbox"
      checked={props.value}
      disabled={props.disabled}
      title={props.title}
      onChange={e => props.onChange && props.onChange(e.target.checked)}
    />
    <div className={cx('indicator')} />
    {props.label}
  </label>
)
