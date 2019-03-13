// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import {useAsyncComponent} from 'src/helpers'
const importRendererAsync = () => import('./sync_renderer').then(module => module.default)

// Renders a markdown string passed to `children` after asynchronously downloading the markdown renderer javascript bundle
export default (props: {
  className?: string,
  children: string,
}) => {
  const AsyncMarkdownRenderer = useAsyncComponent(importRendererAsync)

  return (
    <div className={props.className}>
      <AsyncMarkdownRenderer children={props.children} />
    </div>
  )
}
