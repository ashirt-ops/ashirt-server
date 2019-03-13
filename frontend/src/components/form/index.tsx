// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import classnames from 'classnames/bind'
import {Result} from 'src/global_types'
const cx = classnames.bind(require('./stylesheet'))

const DisplayResult = (props: {
  result: Result<string>,
}) => {
  const result = props.result
  if ('err' in result) return <div className={cx('result', 'error')}>{result.err.message}</div>
  return <div className={cx('result', 'success')}>{result.success}</div>
}

export default (props: {
  children?: React.ReactNode,
  result: Result<string> | null,
  loading: boolean,
  onSubmit: (e: React.FormEvent) => void,
  onCancel?: () => void,
  submitText?: string,
  cancelText?: string,
  submitDanger?: boolean,
}) => {
  const onCancel = (e: React.MouseEvent) => {
    e.preventDefault()
    if (props.onCancel) props.onCancel()
  }

  return (
    <form className={cx('root')} onSubmit={props.onSubmit}>
      {props.result && <DisplayResult result={props.result} />}
      <div className={cx('children')}>
        {props.children}
      </div>
      {props.submitText && (
        <Button
          primary={!props.submitDanger}
          danger={props.submitDanger}
          className={cx('button')}
          loading={props.loading}
          children={props.submitText}
        />
      )}
      {props.onCancel && <Button
        onClick={onCancel}
        className={cx('button')}
        disabled={props.loading}
        children={props.cancelText}
      />}
    </form>
  )
}
