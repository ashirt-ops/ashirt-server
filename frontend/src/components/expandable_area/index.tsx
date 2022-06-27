// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export const ExpandableSection = (props: {
  children: React.ReactNode
  label: React.ReactNode
  className?: string
  labelClassName?: string
  initiallyExpanded?: boolean
  onExpanded?: (expanded: boolean)=>void
}) => {
  const [expanded, setExpanded] = React.useState<boolean>(props.initiallyExpanded || false)
  const expandedState = expanded ? 'expanded' : 'condensed'
  return (
    <section className={cx('expandable-root')}>
      <div className={cx(`expandable-arrow`, `${expandedState}`, props.className)} onClick={() => {
        const newValue = !expanded
        setExpanded(newValue)
        props.onExpanded && props.onExpanded(newValue)
      }}>
        <section className={cx('expandable-section-label', props.labelClassName)}>{props.label}</section>
      </div>
      <div className={cx('expandable-content')}>
        {expanded && props.children}
      </div>
    </section>
  )
}
