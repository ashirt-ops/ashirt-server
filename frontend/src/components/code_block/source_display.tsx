// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { CodeBlock } from 'src/global_types'
import { trimURL } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  value: CodeBlock,
}) => {
  const source = props.value ? props.value.source : null
  if (!source) {
    return <span className={cx('source-none')}>No source available</span>
  }
  const { isAUrl, trimmedValue } = trimURL(source)

  const content = isAUrl
    ? <a onClick={(e) => e.stopPropagation()} href={source} target="_blank">{trimmedValue}</a>
    : <em>{trimmedValue}</em>
  return <span className={cx('source-some')}>From: {content}</span>
}
