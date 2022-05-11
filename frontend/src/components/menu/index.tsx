// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  maxHeight?: number,
  children: React.ReactNode,
}) => (
  <div className={cx('root')} style={{ maxHeight: props.maxHeight }}>
    {props.children}
  </div>
)

export const MenuItem = (props: {
  children: React.ReactNode,
  icon?: string,
  onClick?: (e: React.MouseEvent<Element, MouseEvent>) => void,
  selected?: boolean,
  disabled?: boolean,
  danger?: boolean,
  onKeyUp?: (e: React.KeyboardEvent) => void,
  onKeyDown?: (e: React.KeyboardEvent) => void,
}) => {
  const ref = React.useRef<HTMLButtonElement | null>(null)

  React.useEffect(() => {
    if (props.selected && ref.current != null) {
      ref.current.scrollIntoView({ block: 'nearest' })
    }
  }, [props.selected])

  return (
    <button disabled={props.disabled}
      className={cx('menu-item', {
        selected: props.selected,
        clickable: props.onClick && !props.disabled,
        disabled: props.disabled,
        danger: props.danger,
      })}
      onClick={props.onClick}
      onKeyUp={props.onKeyUp}
      onKeyDown={props.onKeyDown}
      ref={ref}
    >
      {props.icon && (
        <div className={cx('icon')} style={{ backgroundImage: `url(${props.icon})` }} />
      )}
      {props.children}
    </button>
  )
}

export const MenuSeparator = () => <hr className={cx('separator')} />
