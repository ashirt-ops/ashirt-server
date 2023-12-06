import * as React from 'react'
import classnames from 'classnames/bind'
import { NavLink } from 'react-router-dom'

import LoadingSpinner from 'src/components/loading_spinner'

const cx = classnames.bind(require('./stylesheet'))

export type ButtonStyle = {
  active?: boolean
  danger?: boolean,
  primary?: boolean,
  small?: boolean,
  size?: "small" | "medium" | "normal",
}

function styleToClassname(style: ButtonStyle): string {
  return cx({
    active: style.active,
    danger: style.danger,
    primary: style.primary,
    small: style.small === true || style.size === 'small',
    medium: style.size === 'medium'
  });
}

const Button = (props: ButtonStyle & {
  children?: React.ReactNode,
  className?: string,
  disabled?: boolean,
  icon?: string,
  afterIcon?: string,
  loading?: boolean,
  onClick?: (e: React.MouseEvent) => void,
  title?: string,
  doNotSubmit?: boolean,
  autoFocus?: boolean,
}) => {
  return (
    <button
      type={props.doNotSubmit ? 'button' : undefined}
      className={cx('root', props.className, styleToClassname(props), {
        disabled: props.disabled || props.loading,
        loading: props.loading,
      })}
      disabled={props.disabled || props.loading}
      onClick={props.onClick}
      title={props.title}
      autoFocus={props?.autoFocus}
    >
      {props.loading && <LoadingSpinner small className={cx('spinner')} />}
      <div className={cx('children')}>
        {props.icon && <img src={props.icon} />}
        <span>{props.children}</span>
        {props.afterIcon && <img src={props.afterIcon} />}
      </div>
    </button>
  )
}

export const NavLinkButton = (props: ButtonStyle & {
  children: React.ReactNode,
  className?: string,
  exact?: boolean,
  icon?: string,
  to: string,
}) => (
  <NavLink
    end={props.exact}
    className={({isActive}) => cx('root', props.className, styleToClassname(props), isActive ? "active" : null)}
    to={props.to}
  >
    <div className={cx('children')}>
      {props.icon && <img src={props.icon} />}
      <span>{props.children}</span>
    </div>
  </NavLink>
)

export const ButtonGroup = (props: {
  children: React.ReactNode,
  className?: string,
}) => (
  <div className={cx('button-group', props.className)}>
    {props.children}
  </div>
)

export default Button
