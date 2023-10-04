// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import SourceDisplay from './source_display'
import classnames from 'classnames/bind'
import { C2Event } from 'src/global_types'
import ComboBox from 'src/components/combobox'
import { useAsyncComponent } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))
const importAceEditorAsync = () => import('./ace_editor').then(module => module.default)

export const C2EventViewer = (props: {
  value: C2Event,
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
          mode={''}
          value={props.value.command}
        />
      </div>
    </div>
  )
}

// WIP..
export const C2EventEditor = (props: {
  disabled?: boolean,
  onChange: (newValue: C2Event) => void,
  value: C2Event,
}) => {
  const AceEditor = useAsyncComponent(importAceEditorAsync)

  return (
    <div className={cx('code-editor')}>
      <div className={cx('controls')}>
        <Input
          label="C2"
          className={cx('source')}
          value={''}
          onChange={c2 => props.onChange({...props.value, c2})}
          disabled={props.disabled}
        />
      </div>
      <div className={cx('ace')}>
        <AceEditor
          mode={''}
          value={props.value.c2}
          onChange={c2 => props.onChange({...props.value, c2})}
          readOnly={props.disabled}
        />
      </div>
    </div>
  )
}
