// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { User } from 'src/global_types'
import { useDataSource, listUsers } from 'src/services'

import Input from 'src/components/input'
import PopoverMenu from 'src/components/popover_menu'

const userToName = (u: User) => `${u.firstName} ${u.lastName}`

// TODO - REMOVE THIS COMPONENT
// Right now this component is only being used on the operation edit page as a `user search` field.
// However the user edit page should probably combine this component and the user filter into  a single
// component, thus removing the need for the hacks here.
export default (props: {
  value: User|null,
  onChange: (user: User|null) => void,
}) => {
  const ds = useDataSource()
  const [inputValue, setInputValue] = React.useState('')
  const [dropdownVisible, setDropdownVisible] = React.useState(false)
  const [searchResults, setSearchResults] = React.useState<Array<User>>([])
  const [loading, setLoading] = React.useState(false)

  React.useEffect(() => {
    setInputValue(props.value ? userToName(props.value) : '')
  }, [props.value])

  React.useEffect(() => {
    if (inputValue === '') return
    const reload = () => {
      listUsers(ds, { query: inputValue })
        .then(setSearchResults)
        .then(() => setLoading(false))
    }

    // Manually debounce for now since this component is going away
    const timeout = setTimeout(reload, 250)
    return () => { clearTimeout(timeout) }
  }, [ds, inputValue])

  const onRequestClose = () => {
    setDropdownVisible(false)
  }

  const onChange = (v: string) => {
    setLoading(v !== '')
    setInputValue(v)
    if (props.value != null) props.onChange(null)
  }

  const onSelect = (u: User) => {
    props.onChange(u)
    setDropdownVisible(false)
  }

  return (
    <PopoverMenu
      onRequestClose={onRequestClose}
      isOpen={dropdownVisible && !loading && inputValue != ''}
      options={searchResults}
      renderer={userToName}
      onSelect={onSelect}
      noOptionsMessage="No users found"
    >
      <Input
        label="User Search"
        value={inputValue}
        onChange={onChange}
        onFocus={() => setDropdownVisible(true)}
        onClick={() => setDropdownVisible(true)}
        loading={loading}
      />
    </PopoverMenu>
  )
}
