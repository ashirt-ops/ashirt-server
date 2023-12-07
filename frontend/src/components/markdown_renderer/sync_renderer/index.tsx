import * as React from 'react'
import classnames from 'classnames/bind'
// @ts-ignore - module react-markdown does not have associated types (gets imported as any type)
import ReactMarkdown from 'react-markdown'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children: string
}) => (
  <ReactMarkdown className={cx('markdown')}>{props.children}</ReactMarkdown>
)
