// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import classnames from 'classnames/bind'
import {ClickPopover} from 'src/components/popover'
import { NavLink, useLocation } from 'react-router-dom'
const cx = classnames.bind(require('./stylesheet'))

// Usage:
// <ListMenu>
//   <ListItem name="One" selected onSelect={...} />
//   <ListItem name="Two" onSelect={...} />
//   <ListItem name="Three" onSelect={...} />
// </ListMenu>
export default (props: {
  className?: string,
  children: React.ReactNode,
}) => (
    <ul className={cx('root', props.className)}>
      {props.children}
    </ul>
  )

// NavListItem is a version of ListItem that sets and relies on the react-router location.pathname
// to determine selectedness
export const NavListItem = (props: {
  name: string,
  exact?: boolean,
  to: string
}) => {
  const location = useLocation()
  return (
    <li className={cx({ selected: (location.pathname.endsWith(`/${props.to}`)) })}>
      <NavLink end={props.exact} to={props.to} >
        {props.name}
      </NavLink>
    </li>
  )
}

export const ListItem = (props: {
  name: string,
  onSelect: () => void,
  selected: boolean,
}) => (
    <BaseListItem {...props} />
  )

export const ListItemWithSaveButton = (props: {
  name: string,
  onSave: () => void,
  onSelect: () => void,
  selected: boolean,
}) => (
    <BaseListItem {...props} className="withSaveButton">
      <Button
        small
        className={cx('saveButton')}
        children="Save"
        onClick={e => {e.stopPropagation(); props.onSave()}}
      />
    </BaseListItem>
  )

export const ListItemWithMenu = (props: {
  menu: React.ReactElement,
  name: string,
  onSelect: () => void,
  selected: boolean,
}) => (
    <BaseListItem {...props} className="withMenu">
      <div onClick={e => e.stopPropagation()}>
        <ClickPopover content={props.menu} closeOnContentClick>
          <button className={cx('menuButton')}>⋮</button>
        </ClickPopover>
      </div>
    </BaseListItem>
  )

const BaseListItem = (props: {
  className?: string,
  name: string,
  onSelect: () => void,
  selected: boolean,
  children?: React.ReactNode,
}) => (
    <li className={cx({selected: props.selected}, props.className)}>
      <a
        href="#"
        onClick={e => {e.preventDefault(); props.onSelect()}}
        children={props.name}
      />
      {props.children}
    </li>
  )
