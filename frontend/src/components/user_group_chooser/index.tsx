import { useState, useEffect } from 'react'
import Input from 'src/components/input'
import PopoverMenu from 'src/components/popover_menu'
import { type UserGroup } from 'src/global_types'
import { listUserGroups } from 'src/services'

const userGroupToName = (u: UserGroup) => `${u.name}`

export default function UserGroupChooser(props: {
  value: UserGroup | null
  onChange: (userGroup: UserGroup | null) => void
  operationSlug: string
}) {
  const [inputValue, setInputValue] = useState('')
  const [dropdownVisible, setDropdownVisible] = useState(false)
  const [searchResults, setSearchResults] = useState<Array<UserGroup>>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    setInputValue(props.value ? userGroupToName(props.value) : '')
  }, [props.value])

  useEffect(() => {
    if (inputValue === '') return
    const reload = () => {
      listUserGroups({ query: inputValue, operationSlug: props.operationSlug })
        .then(setSearchResults)
        .then(() => setLoading(false))
    }

    const timeout = setTimeout(reload, 250)
    return () => {
      clearTimeout(timeout)
    }
  }, [inputValue, props.operationSlug])

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
