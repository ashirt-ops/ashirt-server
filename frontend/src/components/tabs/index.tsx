import * as React from 'react'
import { default as Button, ButtonGroup } from 'src/components/button'

import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export type Tab = {
  id: string
  disabled?: boolean
  label?: string
  content?: React.ReactNode
  active?: boolean
}

export const findInitialIndex = (allTabs: Array<Tab>, tabIndexOrName?: number | string): number => {
  switch (typeof tabIndexOrName) {
    case 'undefined':
      const undefMatches = allTabs
        .map((t, i) => { return { 'tab': t, 'idx': i } })
        .filter(t => t.tab.active)
      return undefMatches.length > 0 ? undefMatches[0].idx : 0
    case 'number':
      return tabIndexOrName >= 0 && tabIndexOrName < allTabs.length ? tabIndexOrName : 0
    case 'string':
      const matches = allTabs.map((t, i) => { return { 'tab': t, 'idx': i } })
        .filter(t => t.tab.label === tabIndexOrName)
      return matches.length > 0 ? matches[0].idx : 0
    default:
      return 0
  }
}

const TabMenu = (props: {
  tabs: Array<Tab>
  initialActiveTab?: number | string,
  className?: string,
  tabClassName?: string,
  onTabChanged?: (tab: Tab, tabIndex: number) => void,
}) => {
  const defaultActiveTab = findInitialIndex(props.tabs, props.initialActiveTab)
  const [activeTabIndex, setActiveTabIndex] = React.useState<number>(defaultActiveTab)

  const activeTab = props.tabs[activeTabIndex]

  if (props.tabs.length == 0) {
    return null
  }

  return (
    <>
      <ButtonGroup className={cx(props.className)} >
        {props.tabs.map((t, idx) => {
          return (
            <Button
              key={t.id}
              title={t.label}
              className={cx(props.tabClassName)}
              onClick={(e) => {
                e.preventDefault()
                setActiveTabIndex(idx)
                if (props.onTabChanged) {
                  props.onTabChanged(props.tabs[idx], idx)
                }
              }}
              disabled={t.disabled}
              active={idx == activeTabIndex}>{t.label}</Button>
          )
        })}
      </ButtonGroup>
      <div>
        {activeTab.content}
      </div>
    </>
  )
}

export default TabMenu
