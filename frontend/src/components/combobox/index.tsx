// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import PopoverMenu from 'src/components/popover_menu'

import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export type ComboBoxItem<T> = {
  name: string,
  value: T,
}

function valueToName<T>(value: T, options: Array<ComboBoxItem<T>>): string {
  for (let option of options) {
    if (option.value === value) return option.name
  }
  throw Error(`Bad value: ${value}`)
}

function filterOptions<T>(allOptions: Array<ComboBoxItem<T>>, filterValue: string): Array<ComboBoxItem<T>> {
  filterValue = filterValue.trim().toLowerCase()
  if (filterValue === '') return allOptions
  return allOptions.filter(v => v.name.toLowerCase().indexOf(filterValue) > -1)
}

/**
 * SmartComboBox (default export) is modeled after a type-able select-style component. This component
 * manages its own state, and conveys selections with the `onChange` prop. If you need more control
 * of the combobox, see {@link DumbComboBox}. There is also a true Select component available which
 * may suit your needs better. 
 */
export default function SmartComboBox<T>(props: {
  options: Array<ComboBoxItem<T>>,
  onChange: (newValue: T) => void,
  label: string,
  value: T,
  className?: string,
  disabled?: boolean,
}) {
  const [cbState, setCbState] = React.useState(initialComboBoxState(props.value, props.options))

  return ( <DumbComboBox 
    options={props.options}
    label={props.label}
    inputValue={cbState.inputValue}
    popoverOpen={cbState.popoverOpen}
    value={cbState.value}
    disabled={props.disabled}
    className={props.className}
    onChange={(stateChange: ComboBoxStateEvent<T>) => {
        const updatedState = handleTypicalComboBoxStateChangeEvent(cbState, stateChange)
        setCbState(updatedState)
        props.onChange(updatedState.value)
      }
    }
  />)
}

/**
 * DumbComboBox is a component that replicates a select while allowing typing. This component offers
 * complete control over the combobox experience. Typically, you will want to pair this with 
 * {@link initialComboBoxState} and/or do state management with
 * {@link handleTypicalComboBoxStateChangeEvent}
 * In fact, if that's all you want to do, you should probably use the default export 
 * ({@link SmartComboBox}).
 */
export function DumbComboBox<T>(props: {
  // display
  options: Array<ComboBoxItem<T>>,
  label: string,
  // state
  value: T,
  popoverOpen: boolean,
  inputValue: string,
  onChange: (newState: ComboBoxStateEvent<T>) => void,

  // optional values
  className?: string,
  disabled?: boolean,
  dontFilterOptions?: boolean
}) {
  const { onChange } = props
  const filteredOptions = props.dontFilterOptions
  ? props.options
  : filterOptions(props.options, props.inputValue)

  return (
    <StatelessComboBox
      label={props.label}
      options={filteredOptions}
      value={props.value}
      popoverOpen={props.popoverOpen}
      inputValue={props.inputValue}
      disabled={props.disabled}
      className={props.className}
      onRequestPopoverClose={() => onChange({ event: 'popover-closed' })}

      onPopoverSelect={(item: ComboBoxItem<T>) => { onChange({ event: 'popover-selected', item }) }}

      onInputChange={(newVal) => {
        onChange({
          event: 'input-changed',
          inputValue: newVal
        })
      }}
      onInputClick={() => onChange({ event: 'input-click' })}
      onInputFocus={() => onChange({ event: 'input-focused' })}
      onInputKeyDown={() => onChange({ event: 'input-keydown' })}
    />
  )
}

function StatelessComboBox<T>(props: {
  options: Array<ComboBoxItem<T>>,
  label: string,
  value: T,
  inputValue: string,
  popoverOpen: boolean,
  onPopoverSelect: (item: ComboBoxItem<T>) => void,
  onRequestPopoverClose: () => void,
  onInputKeyDown?: () => void
  onInputFocus?: () => void
  onInputChange?: (newVal: string) => void
  onInputClick?: () => void,
  className?: string,
  disabled?: boolean,
}) {
  return (
    <PopoverMenu
      isOpen={props.popoverOpen}
      onRequestClose={props.onRequestPopoverClose}
      onSelect={props.onPopoverSelect}
      options={props.options}
      renderer={item => item.name}
      noOptionsMessage="No Matches"
    >
      <Input
        label={props.label}
        className={cx('arrow', props.className)}
        onChange={props.onInputChange}
        onFocus={props.onInputFocus}
        onKeyDown={props.onInputKeyDown}
        onClick={props.onInputClick}
        value={props.inputValue}
        disabled={props.disabled}
      />
    </PopoverMenu>
  )
}

export function initialComboBoxState<T>(val: T, options: Array<ComboBoxItem<T>>): ComboBoxManagedState<T>
export function initialComboBoxState<T>(val: ComboBoxItem<T>, options?: never): ComboBoxManagedState<T>

export function initialComboBoxState<T>(
  val: T | ComboBoxItem<T>, 
  options?: Array<ComboBoxItem<T>>
): ComboBoxManagedState<T> {
  const cbItem: ComboBoxItem<T> = Array.isArray(options)
    ? { name: valueToName(val, options), value: val as T}
    : val as ComboBoxItem<T>
  return {
    popoverOpen: false,
    inputValue: cbItem.name,
    value: cbItem.value,
    lastSelection: cbItem
  }
}

/**
 * handleTypicalComboBoxStateChangeEvent is responsible for taking an existing combobox state, and
 * generating a new state off of the given combobox event (for typical usecases). Specifically,
 * this is focused on making sure that the input value is correct, and that the popover is in the
 * correct (i.e. open or closed). This can be applied either in whole, or in part if you want to
 * update how a particular event is handled.
 * 
 * In many ways, this is like a mini redux-like reducer
 * 
 * @param currentState The current state of the combobox
 * @param e the combobox event
 * @returns A revised state 
 */
export function handleTypicalComboBoxStateChangeEvent<T>(
  currentState: ComboBoxManagedState<T>,
  e: ComboBoxStateEvent<T>,
): ComboBoxManagedState<T> {

  if (e.event == 'popover-selected') {
    return {
      ...currentState,
      popoverOpen: false,
      inputValue: e.item.name,
      value: e.item.value,
      lastSelection: e.item
    }
  }
  if (e.event == 'input-changed') {
    return {
      ...currentState,
      inputValue: e.inputValue,
    }
  }
  if (e.event == 'input-keydown') {
    return {
      ...currentState,
      popoverOpen: true,
    }
  }
  if (e.event == 'input-focused') {
    return {
      ...currentState,
      popoverOpen: true,
      inputValue: '',
    }
  }
  if (e.event == 'popover-closed') {
    return {
      ...currentState,
      popoverOpen: false,
      inputValue: currentState.lastSelection.name,
      value: currentState.lastSelection.value,
    }
  }
  // the rest of the events are unhandled
  return currentState
}

export type ComboBoxManagedState<T> = {
  popoverOpen: boolean
  inputValue: string
  value: T
  lastSelection: ComboBoxItem<T>,
}

export type ComboBoxStateEvent<T> =
  | ComboBoxStateFocused
  | ComboBoxStateInputChange
  | ComboBoxStateInputDown
  | ComboBoxStateItemSelected<T>
  | ComboBoxStatePopoverClosed
  | ComboBoxStateInputClicked

type ComboBoxStateFocused = {
  event: "input-focused"
}

type ComboBoxStateInputChange = {
  event: 'input-changed'
  inputValue: string
}

type ComboBoxStateInputDown = {
  event: 'input-keydown'
}

type ComboBoxStateItemSelected<T> = {
  event: 'popover-selected'
  item: ComboBoxItem<T>
}

type ComboBoxStatePopoverClosed = {
  event: 'popover-closed'
}

type ComboBoxStateInputClicked = {
  event: 'input-click'
}
