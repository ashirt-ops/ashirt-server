import * as React from 'react'

import SyncMarkdownRenderer from './sync_renderer'

export default (props: {
  className?: string,
  children: string,
}) => (
  <div className={props.className}>
    <SyncMarkdownRenderer children={props.children} />
  </div>
)
