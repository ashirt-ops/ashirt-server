import * as React from 'react'
import classnames from 'classnames/bind'

import { useAsyncComponent } from 'src/helpers'

const cx = classnames.bind(require('./components_ss'))
const importAceEditorAsync = () => import('../code_block/ace_editor').then(module => module.default)

export const RawContent = (props: {
  content: string
  language?: string
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return props.content == ''
      ? <em>No Content</em>
      : <div className={cx('code-viewer')}>
        <div className={cx('ace')}>
          <AceEditor
            readOnly
            mode={props.language ? props.language : ''}
            value={props.content}
          />
        </div>
    </div>
}

export const EvidenceHeader = (props: {
  creator: string,
  version: string
}) => (
  <div className={cx('header')}>
    From:
    <em className={cx('header-creator')}>
      {props.creator}
    </em>
    @
    <em className={cx('header-version')}>
      {props.version}
    </em>
  </div>
)
