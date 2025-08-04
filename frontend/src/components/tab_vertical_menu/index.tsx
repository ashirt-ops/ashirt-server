import * as React from 'react'
import classnames from 'classnames/bind'
import { Outlet } from 'react-router'

import { default as ListMenu, NavListItem, ListItem } from 'src/components/list_menu'

const cx = classnames.bind(require('./stylesheet'))

export type Tab = {
  id: string,
  disabled?: boolean
  label: string
  content?: React.ReactNode
}

export type NavTab = {
  id: string,
  label: string
  query?: Record<string, string>
}

export const NavVerticalTabMenu = (props: {
  title: string
  tabs: Array<NavTab>
  children: React.ReactNode
}) => {
  return (
    <nav className={cx('root')}>
      <div className={cx("tabmenu")}>
        <header>{props.title}</header>
        <ListMenu>
          {props.tabs.map(tab => (
            <NavListItem key={tab.id} name={tab.label} to={tab.id} query={tab.query} />
          ))}
        </ListMenu>
      </div>
      <div className={cx("content")}>
        {props.children}
        <Outlet />
      </div>
    </nav>
  )
}

export const VerticalTabMenu = (props: {
  title: string,
  tabs: Array<Tab>
  selectedTabIndex: number
  onTabChanged?: (tab: Tab, tabIndex: number) => void,
}) => (
  <div className={cx('root')}>
    <div className={cx("tabmenu")}>
      <header>{props.title}</header>
      <ListMenu>
        {props.tabs.map((tab, index) => <ListItem
          key={tab.id}
          name={tab.label}
          selected={index == props.selectedTabIndex}
          onSelect={() => props.onTabChanged?.(tab, index)}
        />)}
      </ListMenu>
    </div>
    <div className={cx("content")}>
      {props.tabs[props.selectedTabIndex].content}
    </div>
  </div>
)
