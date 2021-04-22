// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Popover from 'src/components/popover'
import { default as Menu, MenuItem } from 'src/components/menu'

// input key down helpers:
const isDown = (e: React.KeyboardEvent) => e.key === 'ArrowDown' || (e.ctrlKey && e.key === 'n') // down arrow or ctrl-n
const isUp = (e: React.KeyboardEvent) => e.key === 'ArrowUp' || (e.ctrlKey && e.key === 'p') // up arrow or ctrl-p

export default function PopoverMenu<T>(props: {
  children: React.ReactNode,
  isOpen: boolean,
  onRequestClose: () => void,
  onSelect: (item: T) => void,
  options: Array<T>,
  renderer: (item: T) => React.ReactNode,
  iconRenderer?: (item: T) => string,
  noOptionsMessage?: React.ReactNode,
}) {
  const [selectedIndex, setSelectedIndex] = React.useState(0)

  const selectItem = () => {
    if (selectedIndex < 0 || selectedIndex >= props.options.length) return
    props.onSelect(props.options[selectedIndex])
  }

  const changeSelectedIndex = (delta: number) => {
    if (!props.isOpen) return
    const maxIndex = props.options.length - 1
    const boundedIndex = Math.max(0, Math.min(maxIndex, selectedIndex + delta))
    setSelectedIndex(boundedIndex)
  }

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') selectItem()
    else if (e.key === 'Escape') props.onRequestClose()
    else if (isDown(e)) changeSelectedIndex(1)
    else if (isUp(e)) changeSelectedIndex(-1)
    else if (!e.ctrlKey) return
    e.preventDefault()
  }

  return (
    <Popover
      content={<PopoverMenuContent selectedIndex={selectedIndex} {...props} />}
      onRequestClose={props.onRequestClose}
      isOpen={props.isOpen}
    >
      <div onKeyDown={onKeyDown}>
        {props.children}
      </div>
    </Popover>
  )
}

function PopoverMenuContent<T>(props: {
  onSelect: (item: T) => void,
  options: Array<T>,
  renderer: (item: T) => React.ReactNode,
  selectedIndex: number,
  iconRenderer?: (item: T) => string,
  noOptionsMessage?: React.ReactNode,
}) {
  return (
    <Menu maxHeight={200}>
      {props.options.map((v, i) => (
        <MenuItem
          key={i}
          selected={i === props.selectedIndex}
          icon={props.iconRenderer && props.iconRenderer(v)}
          children={props.renderer(v)}
          onClick={() => props.onSelect(v)}
        />
      ))}
      {props.options.length === 0 && (
        <MenuItem children={props.noOptionsMessage} />
      )}
    </Menu>
  )
}
