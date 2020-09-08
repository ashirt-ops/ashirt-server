// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

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
