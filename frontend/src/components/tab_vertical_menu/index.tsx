// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { Switch, Route, RouteComponentProps } from 'react-router-dom'
import { subUrl } from 'src/helpers'

import { default as ListMenu, NavListItem, ListItem } from 'src/components/list_menu'

const cx = classnames.bind(require('./stylesheet'))

export type Tab = {
  id: string,
  disabled?: boolean
  label: string
  content?: React.ReactNode
}

export const NavVerticalTabMenu = (props: RouteComponentProps & {
  title: string,
  tabs: Array<Tab>
}) => {
  return (
    <div className={cx('root')}>
      <div className={cx("tabmenu")}>
        <header>{props.title}</header>
        <ListMenu>
          {props.tabs.map((tab) => <NavListItem
            key={tab.id}
            name={tab.label}
            to={subUrl(props, { view: tab.id })} />)}
        </ListMenu>
      </div>
      <div className={cx("content")}>
        <Switch>
          {props.tabs.map((tab) => {
            return <Route key={tab.id} exact path={subUrl(props, { view: tab.id })} render={() => tab.content} />
          })}
        </Switch>
      </div>
    </div>
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
