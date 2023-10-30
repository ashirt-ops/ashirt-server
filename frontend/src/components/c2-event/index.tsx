// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { default as Input, SharedProps} from 'src/components/input'
import { C2Event } from 'src/global_types'
import ComboBox from 'src/components/combobox'
import { useAsyncComponent } from 'src/helpers'
import LoadingSpinner from 'src/components/loading_spinner'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))


export const C2EventTextArea = React.forwardRef((props: SharedProps & {
}, ref: React.RefObject<HTMLTextAreaElement>) => (
  <WithLabel className={cx('root', props.className)} label={props.label}>
    <textarea
      ref={ref}
      className={cx('c2-event-textarea')}
      disabled={props.disabled}
      name={props.name}
      onBlur={props.onBlur}
      onChange={e => { if (props.onChange) props.onChange(e.target.value) }}
      onClick={props.onClick}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      placeholder={props.placeholder}
      value={props.value}
    />
  </WithLabel>
))


export const C2EventViewer = (props: {
  disabled?: boolean,
  //onChange: (newValue: C2Event) => void,
  value: C2Event,
}) => {
  return (
    <div className={cx('c2-event-grid')}>
      <div className={cx('c2framework')}>
        <Input 
                label="C2"
                className={cx('c2-event-input')}
                value={props.value.c2 || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('operator')}>
        <Input
                label="Operator"
                className={cx('c2-event-input')}
                value={props.value.c2Operator || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('user-context')}>
        <Input
                label="User Context"
                className={cx('c2-event-input')}
                value={props.value.userContext || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('beacon')}>
        <Input
                label="Beacon"
                className={cx('c2-event-input')}
                value={props.value.beacon || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('hostname')}>
        <Input
                label="Hostname"
                className={cx('c2-event-input')}
                value={props.value.hostname || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('intIP')}>
        <Input
                label="Internal IP"
                className={cx('c2-event-input')}
                value={props.value.internalIP || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('extIP')}>
        <Input
                label="Exeternal IP"
                className={cx('c2-event-input')}
                value={props.value.externalIP || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('process')}>
        <Input
                label="Process"
                className={cx('c2-event-input')}
                value={props.value.processName || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('procID')}>
        <Input
            label="Process ID"
            type="number"
            className={cx('c2-event-input')}
            value={props.value.processID !== undefined ? props.value.processID.toString() : ''}
            disabled={props.disabled}
        />
      </div>
      <div className={cx('command')}>
        <Input
                label="Command"
                className={cx('c2-event-input')}
                value={props.value.command || ''}
                disabled={props.disabled}
                readOnly
            />
      </div>
      <div className={cx('result')}>
        <C2EventTextArea
                  label="Result" className={cx('c2-event-input', 'resizeable')}
                  value={props.value.result || ''}
                  disabled={props.disabled}
                  readOnly
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
  return (
    <div className={cx('c2-event-grid')}>
      <div className={cx('c2framework')}>
        <Input 
                label="C2"
                className={cx('c2-event-input')}
                value={props.value.c2 || ''}
                disabled={props.disabled}
                onChange={c2 => props.onChange({...props.value, c2})}
            />
      </div>
      <div className={cx('operator')}>
        <Input
                label="Operator"
                className={cx('c2-event-input')}
                value={props.value.c2Operator || ''}
                disabled={props.disabled}
                onChange={c2Operator => props.onChange({...props.value, c2Operator})}
            />
      </div>
      <div className={cx('user-context')}>
        <Input
                label="User Context"
                className={cx('c2-event-input')}
                value={props.value.userContext || ''}
                disabled={props.disabled}
                onChange={userContext => props.onChange({...props.value, userContext})}
            />
      </div>
      <div className={cx('beacon')}>
        <Input
                label="Beacon"
                className={cx('c2-event-input')}
                value={props.value.beacon || ''}
                disabled={props.disabled}
                onChange={beacon => props.onChange({...props.value, beacon})}
            />
      </div>
      <div className={cx('hostname')}>
        <Input
                label="Hostname"
                className={cx('c2-event-input')}
                value={props.value.hostname || ''}
                disabled={props.disabled}
                onChange={hostname => props.onChange({...props.value, hostname})}
            />
      </div>
      <div className={cx('intIP')}>
        <Input
                label="Internal IP"
                className={cx('c2-event-input')}
                value={props.value.internalIP || ''}
                disabled={props.disabled}
                onChange={internalIP => props.onChange({...props.value, internalIP})}
            />
      </div>
      <div className={cx('extIP')}>
        <Input
                label="Exeternal IP"
                className={cx('c2-event-input')}
                value={props.value.externalIP || ''}
                disabled={props.disabled}
                onChange={externalIP => props.onChange({...props.value, externalIP})}
            />
      </div>
      <div className={cx('process')}>
        <Input
                label="Process"
                className={cx('c2-event-input')}
                value={props.value.processName || ''}
                disabled={props.disabled}
                onChange={processName => props.onChange({...props.value, processName})}
            />
      </div>
      <div className={cx('procID')}>
        <Input
            label="Process ID"
            type="number"
            className={cx('c2-event-input')}
            value={props.value.processID !== undefined ? props.value.processID.toString() : ''}
            onChange={processID =>{
              const pidNum = parseFloat(processID);
              if (!isNaN(pidNum)) {
                props.onChange({ ...props.value, processID: pidNum });
              } else {
                console.error('Invalid input. Please enter a number.'); // cowboy error handling... whats the ASHIRT way to do this?
              }
            }}
            disabled={props.disabled}
        />
      </div>
      <div className={cx('command')}>
        <Input
                label="Command"
                className={cx('c2-event-input')}
                value={props.value.command || ''}
                disabled={props.disabled}
                onChange={command => props.onChange({...props.value, command})}
            />
      </div>
      <div className={cx('result')}>
        <C2EventTextArea
                  label="Result" className={cx('c2-event-input', 'resizeable')}
                  value={props.value.result || ''}
                  disabled={props.disabled}
                  onChange={result => props.onChange({...props.value, result})}
              />
        </div>
    </div>
  )
}