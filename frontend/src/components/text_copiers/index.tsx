// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { default as Button, ButtonStyle } from 'src/components/button'
import Input from 'src/components/input'
import Tooltip from 'src/components/tooltip'

const cx = classnames.bind(require('./stylesheet'))

export const CopyTextButton = (props: ButtonStyle & {
  icon?: string,
  title?: string,
  textToCopy: string,
  children?: React.ReactNode,
  disabled?: boolean,
  className?: string,
}) => {
  const [tooltipOpen, setTooltipOpen] = React.useState(false)

  const onClick = async () => {
    try {
      await navigator.clipboard.writeText(props.textToCopy)
      setTooltipOpen(true)
      await new Promise(resolve => setTimeout(resolve, 2000))
      setTooltipOpen(false)
    } catch (err) {
      console.error(err)
      prompt('Failed to copy to clipboard. Please manually copy:', props.textToCopy)
    }
  }

  return (
    <Tooltip content="Copied!" isOpen={tooltipOpen}>
      <Button {...props} onClick={onClick} children={props.children} />
    </Tooltip>
  )
}

export const InputWithCopyButton = (props: {
  label: string,
  value: string,
}) => (
  <div className={cx('input-with-copy')}>
    <Input label={props.label} readOnly value={props.value} />
    <CopyTextButton textToCopy={props.value} title="Copy to clipboard" icon={require('./copy.svg')} />
  </div>
)

export type CopyBlob = {
  type: 'blob'
  content: () => Promise<Blob>
}

export type CopyText = {
  type: 'text'
  content: () => Promise<string>
}

export const CopyButton = (props: ButtonStyle &
  (CopyBlob | CopyText) &
{
  title?: string,
  children?: React.ReactNode,
  disabled?: boolean,
  className?: string,
}) => {
  type SuccessType = "normal" | "success" | "failure"
  const [showIcon, setShowIcon] = React.useState<SuccessType>("normal")

  const onClick = async (e: React.MouseEvent<Element, MouseEvent>) => {
    let txt = ''
    let iconType: SuccessType = "success"
    try {
      if (props.type === 'text') {
        txt = await props.content()
        await navigator.clipboard.writeText(txt)
      }
      else { // image
        // Jan, 2022: This feature is emerging, and poorly supported
        // Chrome only supports pngs
        // firefox has support behind a flag, and does not seem to support jpgs
        const data = await props.content()
        // @ts-ignore
        const item = new ClipboardItem({ [data.type]: data })
        // @ts-ignore
        await navigator.clipboard.write([item])
      }
    } catch (err) {
      console.error(err)
      if (props.type === 'text') {
        prompt('Failed to copy to clipboard. Please manually copy:', txt)
      }
      else {
        prompt('Failed to copy to clipboard.')
      }
      iconType = "failure"
    }

    setShowIcon(iconType)
    await new Promise(resolve => setTimeout(resolve, 2000))
    setShowIcon("normal")
  }
  const { className, ...renderProps } = props
  return (
    <Button
      className={cx('copy-success-button', props.className)}
      {...renderProps}
      icon={showIcon == 'success'
        ? require('src/res/success.svg')
        : showIcon == 'failure'
          ? require('src/res/error.svg')
          : require('src/res/copy.svg')
      }
      onClick={onClick}
      children={props.children}
    />
  )
}
