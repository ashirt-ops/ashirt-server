// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import PopoverMenu from 'src/components/popover_menu'
import {UserGroup} from 'src/global_types'
import {listUserGroups} from 'src/services'

const userGroupToName = (u: UserGroup) => `${u.name}`

// TODO - REMOVE THIS COMPONENT
// Right now this component is only being used on the operation edit page as a `user group search` field.
// However the user group edit page should probably combine this component and the user group filter into  a single
// component, thus removing the need for the hacks here.
export default (props: {
  value: UserGroup|null,
  onChange: (userGroup: UserGroup|null) => void,
}) => {
  const [inputValue, setInputValue] = React.useState('')
  const [dropdownVisible, setDropdownVisible] = React.useState(false)
  const [searchResults, setSearchResults] = React.useState<Array<UserGroup>>([])
  const [loading, setLoading] = React.useState(false)

  React.useEffect(() => {
    setInputValue(props.value ? userGroupToName(props.value) : '')
  }, [props.value])

  React.useEffect(() => {
    if (inputValue === '') return
    const reload = () => {
      listUserGroups({query: inputValue})
        .then(setSearchResults)
        .then(() => setLoading(false))
    }

    // Manually debounce for now since this component is going away
    const timeout = setTimeout(reload, 250)
    return () => { clearTimeout(timeout) }
  }, [inputValue])

  const onRequestClose = () => {
    setDropdownVisible(false)
  }

  const onChange = (v: string) => {
    setLoading(v !== '')
    setInputValue(v)
    if (props.value != null) props.onChange(null)
  }

  const onSelect = (u: UserGroup) => {
    props.onChange(u)
    setDropdownVisible(false)
  }

  return (
    <PopoverMenu
      onRequestClose={onRequestClose}
      isOpen={dropdownVisible && !loading && inputValue != ''}
      options={searchResults}
      renderer={userGroupToName}
      onSelect={onSelect}
      noOptionsMessage="No user groups found"
    >
      <Input
        label="User Group Search"
        value={inputValue}
        onChange={onChange}
        onFocus={() => setDropdownVisible(true)}
        onClick={() => setDropdownVisible(true)}
        loading={loading}
      />
    </PopoverMenu>
  )
}
