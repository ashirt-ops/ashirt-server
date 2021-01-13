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
    content = [<em className={cx('pretty-headers-key')}>No Captured Headers</em>]
  }
  else {
    content = props.headers
      .sort((a, b) => a.name.toLowerCase().localeCompare(b.name.toLowerCase()))
      .map((header) => (
        <>
          <em className={cx('pretty-headers-key')}>{header.name}:</em>
          <span className={cx('pretty-headers-value')}>{header.value}</span>
        </>
      ))
  }

  return (
    <div className={cx('pretty-headers-container')}>
      {content.map((el, i) => <div key={i} className={cx('pretty-headers-entry')}>{el}</div>)}
    </div>
  )
}

export const RawContent = (props: {
  content: string
  language?: string
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('code-viewer')}>
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
