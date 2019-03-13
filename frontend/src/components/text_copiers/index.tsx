// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { default as Button, ButtonStyle} from 'src/components/button'
import Input from 'src/components/input'
import Tooltip from 'src/components/tooltip'

const cx = classnames.bind(require('./stylesheet'))

export const CopyTextButton = (props: ButtonStyle & {
  icon?: string,
  title?: string,
  textToCopy: string,
  children?: React.ReactNode,
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
