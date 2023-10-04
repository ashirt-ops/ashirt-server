// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import SourceDisplay from './source_display'
import classnames from 'classnames/bind'
import supportedLanguages from './supported_languages'
import { CodeBlock } from 'src/global_types'
import ComboBox from 'src/components/combobox'
import { useAsyncComponent } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))
const importAceEditorAsync = () => import('./ace_editor').then(module => module.default)

export const CodeBlockEditor = (props: {
  disabled?: boolean,
  onChange: (newValue: CodeBlock) => void,
  value: CodeBlock,
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('code-editor')}>
      <div className={cx('controls')}>
        <ComboBox
          label="Language"
          className={cx('language')}
          options={supportedLanguages}
          value={props.value.language}
          onChange={language => props.onChange({...props.value, language})}
          disabled={props.disabled}
          nonValueDefault=""
        />
        <Input
          label="Source"
          className={cx('source')}
          value={props.value.source || ''}
          onChange={source => props.onChange({...props.value, source})}
          disabled={props.disabled}
        />
      </div>
      <div className={cx('ace')}>
        <AceEditor
          mode={props.value.language}
          value={props.value.code}
          onChange={code => props.onChange({...props.value, code})}
          readOnly={props.disabled}
        />
      </div>
    </div>
  )
}

export const CodeBlockViewer = (props: {
  value: CodeBlock,
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('code-viewer')}>
      <div className={cx('source')}>
        <SourceDisplay value={props.value} />
      </div>
      <div className={cx('ace')}>
        <AceEditor
          readOnly
          mode={props.value.language}
          value={props.value.code}
        />
      </div>
    </div>
  )
}

export const SourcelessCodeblock = (props: {
  code: string,
  language: string | null
  className?: string
  editable?: boolean
  onChange?: (newValue: string) => void,
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('code-viewer', props.className)}>
      <div className={cx('ace', 'no-source')}>
        <AceEditor
          readOnly={props.editable ? false : true}
          mode={props.language || ''}
          value={props.code}
          onChange={props.onChange}
        />
      </div>
    </div>
  )
}

