import { type KeyboardEvent, type ReactNode, useState } from 'react'
import Popover from 'src/components/popover'
import { default as Menu, MenuItem } from 'src/components/menu'

// input key down helpers:
const isDown = (e: KeyboardEvent) => e.key === 'ArrowDown' || (e.ctrlKey && e.key === 'n') // down arrow or ctrl-n
const isUp = (e: KeyboardEvent) => e.key === 'ArrowUp' || (e.ctrlKey && e.key === 'p') // up arrow or ctrl-p

export default function PopoverMenu<T>(props: {
  children: ReactNode
  isOpen: boolean
  onRequestClose: () => void
  onSelect: (item: T, e?: KeyboardEvent) => void
  options: Array<T> | (() => Array<T>)
  renderer: (item: T) => ReactNode
  iconRenderer?: (item: T) => string
  noOptionsMessage?: ReactNode
  onKeyModifierChanged?: (modifiers: KeyboardModifiers) => void
}) {
  const [selectedIndex, setSelectedIndex] = useState(0)

  const getOptions = (): Array<T> => {
    return typeof props.options == 'function' ? props.options() : props.options
  }

  const selectItem = (e: KeyboardEvent) => {
    if (selectedIndex < 0 || selectedIndex >= props.options.length) return
    props.onSelect(getOptions()[selectedIndex], e)
  }

  const changeSelectedIndex = (delta: number) => {
    if (!props.isOpen) return
    const maxIndex = props.options.length - 1
    const boundedIndex = Math.max(0, Math.min(maxIndex, selectedIndex + delta))
    setSelectedIndex(boundedIndex)
  }

  const onKeyUp = (e: KeyboardEvent) => {
    props.onKeyModifierChanged?.(reactKeyToModifiers(e))
  }

  const onKeyDown = (e: KeyboardEvent) => {
    props.onKeyModifierChanged?.(reactKeyToModifiers(e))

    if (e.key === 'Enter') selectItem(e)
    else if (e.key === 'Escape') props.onRequestClose()
    else if (e.key === 'Tab') {
      props.onRequestClose()
      return
    } else if (isDown(e)) changeSelectedIndex(1)
    else if (isUp(e)) changeSelectedIndex(-1)
    else if (!e.ctrlKey) return
    e.preventDefault()
  }

  return (
    <Popover
      content={
        <PopoverMenuContent selectedIndex={selectedIndex} {...props} options={getOptions()} />
      }
      onRequestClose={props.onRequestClose}
      isOpen={props.isOpen}
    >
      <div onKeyDown={onKeyDown} onKeyUp={onKeyUp}>
        {props.children}
      </div>
    </Popover>
  )
}

function PopoverMenuContent<T>(props: {
  onSelect: (item: T) => void
  options: Array<T>
  renderer: (item: T) => ReactNode
  selectedIndex: number
  iconRenderer?: (item: T) => string
  noOptionsMessage?: ReactNode
  onKeyModifierChanged?: (modifiers: KeyboardModifiers) => void
}) {
  const onKey = (e: KeyboardEvent) => {
    props.onKeyModifierChanged?.(reactKeyToModifiers(e))
  }

  return (
    <Menu maxHeight={200}>
      {props.options.map((v, i) => (
        <MenuItem
          key={i}
          selected={i === props.selectedIndex}
          icon={props.iconRenderer && props.iconRenderer(v)}
          children={props.renderer(v)}
          onKeyDown={onKey}
          onKeyUp={onKey}
          onClick={(e) => {
            props.onKeyModifierChanged?.(reactKeyToModifiers(e))
            props.onSelect(v)
          }}
        />
      ))}
      {props.options.length === 0 && <MenuItem children={props.noOptionsMessage} />}
    </Menu>
  )
}

export type KeyboardModifiers = {
  altKey: boolean
  ctrlKey: boolean
  metaKey: boolean
  shiftKey: boolean
}

const reactKeyToModifiers = (e: KeyboardModifiers) => {
  return {
    altKey: e.altKey,
    ctrlKey: e.ctrlKey,
    metaKey: e.metaKey,
    shiftKey: e.shiftKey,
  }
}
