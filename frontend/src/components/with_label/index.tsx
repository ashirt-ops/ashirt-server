// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children?: React.ReactNode,
  className?: string,
  label?: string,
}) => {
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
