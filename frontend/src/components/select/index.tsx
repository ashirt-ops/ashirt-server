// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import WithLabel from 'src/components/with_label'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children: React.ReactNode,
  className?: string,
  disabled?: boolean,
  label?: string,
  onChange: (v: string) => void,
  value: string,
}) => (
  <WithLabel className={cx(props.className)} label={props.label}>
    <div className={cx('wrapper')}>
      <select value={props.value} onChange={e => props.onChange(e.target.value)} disabled={props.disabled}>
        {props.children}
      </select>
      <div className={cx('visual', {disabled: props.disabled})}>
        {getNameForValue(props.value, props.children)}
      </div>
    </div>
  </WithLabel>
)

function getNameForValue(curValue: string, children: React.ReactNode): string {
  let name = curValue
  React.Children.map(children, child => {
    if (!React.isValidElement(child)) return
    let {value, children} = child.props
    if (value == null) value = children
    if (value === curValue) name = children
  })
  return name
}
