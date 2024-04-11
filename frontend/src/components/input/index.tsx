// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import LoadingSpinner from 'src/components/loading_spinner'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export type SharedProps = {
  className?: string,
  disabled?: boolean,
  label?: string,
  name?: string,
  onBlur?: () => void,
  onChange?: (newValue: string) => void,
  onClick?: () => void,
  onFocus?: () => void,
  onKeyDown?: (e: React.KeyboardEvent) => void,
  placeholder?: string,
  readOnly?: boolean,
  value?: string,
  adjustHeight?: boolean,
}

export default React.forwardRef((props: SharedProps & {
  icon?: string,
  loading?: boolean,
  type?: string,
  inputClassName?: string,
  autoFocus?: true
}, ref: React.RefObject<HTMLInputElement>) => (
  <WithLabel className={cx('root', props.className)} label={props.label}>
    {props.loading && (
      <LoadingSpinner className={cx('spinner')} small />
    )}
    <input
      ref={ref}
      autoFocus={props.autoFocus}
      className={cx('input', {'has-icon': props.icon != null, loading: props.loading != null}, props.inputClassName)}
      disabled={props.disabled}
      name={props.name}
      onBlur={props.onBlur}
      onChange={e => { if (props.onChange) props.onChange(e.target.value) }}
      onClick={props.onClick}
      onFocus={props.onFocus}
      onKeyDown={props.onKeyDown}
      placeholder={props.placeholder}
      readOnly={props.readOnly}
      style={props.icon != null ? {backgroundImage: `url(${props.icon})`} : {}}
      type={props.type}
      value={props.value}
    />
  </WithLabel>
))

export const TextArea = React.forwardRef((props: SharedProps & {
}, ref: React.RefObject<HTMLTextAreaElement>) => {
  const adjustTextareaHeight = () => {
    const element = document.getElementById("autoResizeTextarea");
    element!.style.height = "auto";
    element!.style.height = element!.scrollHeight < 400 ? element!.scrollHeight + "px": "400px";
  };

  props.adjustHeight && React.useEffect(() => {
    adjustTextareaHeight();
  }, []);

  return (
  <WithLabel className={cx('root', props.className)} label={props.label}>
    <textarea
      ref={ref}
      id="autoResizeTextarea"
      className={cx('input', 'textarea')}
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
)})
