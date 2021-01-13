import * as React from 'react'
import classnames from 'classnames/bind'

import { useAsyncComponent } from 'src/helpers'
import { Header } from 'har-format'

const cx = classnames.bind(require('./stylesheet'))
const importAceEditorAsync = () => import('../code_block/ace_editor').then(module => module.default)


export const PrettyHeaders = (props: {
  headers: Array<Header>
}) => {
  let content
  if (props.headers.length == 0) {
    content = [<em className={cx('pretty-headers-no-content')}>No Captured Headers</em>]
  }
  else {
    content = props.headers
      .sort((a, b) => a.name.toLowerCase().localeCompare(b.name.toLowerCase()))
      .map((h, i) => (
        <div key={i} className={cx('pretty-headers-entry')}>
          <em className={cx('pretty-headers-key')}>{h.name}:</em>
          <span className={cx('pretty-headers-value')}>{h.value}</span>
        </div>
      ))
  }

  return (
      <div className={cx('pretty-headers-container')}>{...content}</div>
  )
}

export const RawContent = (props: {
  content: string
  language?: string
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('ace-container')}>
      <div className={cx('ace')}>
        <AceEditor
          readOnly
          mode={props.language ? props.language : ''}
          value={props.content}
        />
      </div>
    </div>
  )
}

export const Section = (props: {

}) => {
  
}
