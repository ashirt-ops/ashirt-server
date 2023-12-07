import * as React from 'react'
import classnames from 'classnames/bind'

import { SourcelessCodeblock } from '../code_block'

const cx = classnames.bind(require('./components_ss'))

export const RawContent = (props: {
  content: string
  language?: string
}) => {

  return props.content == ''
    ? <em>No Content</em>
    : <SourcelessCodeblock className={cx('code-viewer')}
      code={props.content}
      language={props.language || null} />
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
