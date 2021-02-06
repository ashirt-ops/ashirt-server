import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./expandable_area_ss'))

export const ExpandableSection = (props: {
  children: string
  content: React.ReactNode
  className?: string
  initiallyExpanded?: boolean
}) => {
  const [expanded, setExpanded] = React.useState<boolean>(props.initiallyExpanded || false)
  const expandedState = expanded ? 'expanded' : 'condensed'
  return (
    <section className={cx('expandable-root')}>
      <div className={cx(`expandable-arrow`, `${expandedState}`, props.className)} onClick={() => setExpanded(!expanded)}>
        <section className={cx('expandable-section-label')}>{props.children}</section>
      </div>
      <div className={cx('expandable-content')}>
        {expanded && props.content}
      </div>
    </section>
  )
}

