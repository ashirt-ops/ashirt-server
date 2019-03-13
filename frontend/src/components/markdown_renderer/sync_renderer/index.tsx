// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
// @ts-ignore - module react-markdown does not have associated types (gets imported as any type)
import ReactMarkdown from 'react-markdown'
// @ts-ignore - module react-syntax-highlighter does not have associated types (gets imported as any type)
import SyntaxHighlighter from 'react-syntax-highlighter'
const cx = classnames.bind(require('./stylesheet'))

const CodeBlock = (props: {
  language: string,
  value: string,
}) => (
  <SyntaxHighlighter
    className={cx('code-block')}
    language={props.language}
    children={props.value}
    style={{}}
  />
)

export default (props: {
  children: string
}) => (
  <ReactMarkdown
    className={cx('markdown')}
    source={props.children}
    renderers={{code: CodeBlock}}
  />
)
