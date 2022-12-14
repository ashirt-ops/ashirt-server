// Copyright 2020, Yahoo Inc.
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
  onChange?: (...args: any[]) => void,
}) => (
  <label title={props.title} className={cx('root', props.className)} onClick={e => e.stopPropagation()}>
    <input
      type="checkbox"
      checked={props.value}
      disabled={props.disabled}
      title={props.title}
      onChange={props.onChange}
    />
    <div className={cx('indicator')} />
    {props.label}
  </label>
)
