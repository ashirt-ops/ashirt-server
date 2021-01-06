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

export default function ComboBox<T>(props: {
  options: Array<ComboBoxItem<T>>,
  onChange: (newValue: T) => void,
  label: string,
  value: T,
  className?: string,
  disabled?: boolean,
}) {
  const [inputValue, setInputValue] = React.useState('')
  const [dropdownVisible, setDropdownVisible] = React.useState(false)

  const filteredOptions = filterOptions(props.options, inputValue)

  React.useEffect(() => {
    setInputValue(valueToName(props.value, props.options))
  }, [props.value, props.options])

  const onSelect = (item: ComboBoxItem<T>) => {
    props.onChange(item.value)
    setDropdownVisible(false)
    setInputValue(item.name)
  }

  const onInputFocus = () => {
    setDropdownVisible(true)
    setInputValue('')
  }

  const onRequestClose = () => {
    setDropdownVisible(false)
    setInputValue(valueToName(props.value, props.options))
  }

  return (
    <PopoverMenu
      isOpen={dropdownVisible}
      onRequestClose={onRequestClose}
      onSelect={onSelect}
      options={filteredOptions}
      renderer={item => item.name}
      noOptionsMessage="No Matches"
    >
      <Input
        label={props.label}
        className={cx('arrow', props.className)}
        onChange={setInputValue}
        onFocus={onInputFocus}
        onKeyDown={() => setDropdownVisible(true)}
        value={inputValue}
      />
    </PopoverMenu>
  )
}
