// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { dropRight } from 'lodash'
import classnames from 'classnames/bind'

import WithLabel from 'src/components/with_label'
import PopoverMenu, { KeyboardModifiers } from 'src/components/popover_menu'
import Tag from 'src/components/tag'
import { TagColor } from 'src/helpers'

export * from './creator_chooser'
export * from './evidence_type_chooser'
export * from './tag_chooser'

const cx = classnames.bind(require('./stylesheet'))

export default function BulletChooser<T extends BulletProps>(props: {
  label: string
  options: Array<T> | ((inputVal: string) => Array<T>)
  value: Array<T>
  onNoValueSelected?: (inputValue: string) => Promise<T>
  valueRenderer?: BulletRenderer<T>
  noValueRenderer?: (inputValue: string) => React.ReactNode
  onChange: (tags: Array<T>) => void
  className?: string
  disabled?: boolean
  enableNot?: boolean
}) {
  const [inputValue, setInputValue] = React.useState("")
  const [dropdownVisible, setDropdownVisible] = React.useState(false)
  const [selectedTag, setSelectedTag] = React.useState<number>(-1)
  const [modifierHeld, setHeld] = React.useState(false)

  const getOptions = (): Array<T> => {
    return (
      typeof props.options == 'function'
        ? props.options(inputValue)
        : filterBullets(props.options, inputValue)
    )
  }
  const setModifierHeld = (e: KeyboardModifiers | boolean) => {
    if (props.enableNot) {
      setHeld(typeof e === 'boolean'
        ? e
        : e.ctrlKey || e.altKey
      )
    }
  }

  const { valueRenderer, noValueRenderer } = props
  const renderValFn = valueRenderer ?? StandardBulletRenderer
  const filteredValues = getOptions()

  // selectedTagsAsSet provide quick lookups into selected values
  const selectedTagsAsSet = getBulletIdAsSet(props.value)

  const toggleValue = async (selectedValue: T | null) => {
    const val = selectedValue
      ?? await props.onNoValueSelected?.(inputValue)
      ?? null

    if (val) {
      const modifiedVal = modifierHeld
        ? { ...val, modifier: "not" }
        : val

      let newValues: Array<T>

      const foundItem = selectedTagsAsSet[val.id]
      if (foundItem) {
        const setLessValue = props.value.filter(v => v.id !== val.id)
        newValues = (foundItem.modifier == modifiedVal.modifier)
          ? setLessValue
          : [...setLessValue, modifiedVal]
      }
      else {
        newValues = [...props.value, modifiedVal]
      }

      props.onChange(newValues)
      setInputValue("")
    }
  }

  const onInputKeyDown = (e: React.KeyboardEvent) => {
    setModifierHeld(e)

    if (inputValue === "") {
      if (['Backspace', 'Delete'].includes(e.key)) {
        if (selectedTag != -1) {
          props.onChange([...props.value.slice(0, selectedTag), ...props.value.slice(selectedTag + 1)])
          setSelectedTag(selectedTag - 1)
        }
        else if (e.key != 'Delete') {
          props.onChange(dropRight(props.value))
        }
      }
      else if (e.key === 'ArrowLeft') {
        let index = selectedTag - 1
        setSelectedTag(index > -1 ? index : props.value.length - 1)
      }
      else if (e.key === 'ArrowRight') {
        setSelectedTag((selectedTag + 1) % props.value.length)
      }
      else {
        setSelectedTag(-1)
      }
    }

    setDropdownVisible(true)
  }

  const renderer = (bullet: T | null): React.ReactNode => {
    if (bullet == null) {
      return noValueRenderer?.(inputValue) ?? <StandardNoValRenderer />
    }
    return renderValFn({ bullet })
  }

  // If there are no filtered tags we add a null option to render the "Create new tag" option
  const options: Array<T | null> = filteredValues.length > 0 ? filteredValues : [null]

  return (
    <WithLabel label={props.label}>
      <PopoverMenu
        onRequestClose={() => setDropdownVisible(false)}
        isOpen={dropdownVisible}
        options={modifierHeld
          ? options.map(o => o == null ? null : ({ ...o, modifier: "not" }))
          : options
        }
        renderer={renderer}
        iconRenderer={t => t && selectedTagsAsSet[t.id] && require('./check.svg')}
        onSelect={toggleValue}
        onKeyModifierChanged={setModifierHeld}
      >
        <div className={cx('input', props.className, { focus: dropdownVisible })}>
          {props.value.map(
            (val, idx) => renderValFn({ bullet: val, key: idx, selected: (idx == selectedTag) })
          )}
          <input
            onChange={e => setInputValue(e.target.value)}
            value={inputValue}
            onKeyDown={onInputKeyDown}
            onKeyUp={setModifierHeld}
            onFocus={() => setDropdownVisible(true)}
          />
        </div>
      </PopoverMenu>
    </WithLabel>
  )
}

export type BulletProps = {
  id: string | number
  name: string
  modifier?: "not"
  color?: TagColor
}

type BulletIdSet = Record<BulletProps['id'], BulletProps>
const getBulletIdAsSet = (bullets: Array<BulletProps>): BulletIdSet => {
  const tagIdSet: BulletIdSet = {}
  bullets.forEach(b => tagIdSet[b.id] = b)
  return tagIdSet
}

export type BulletRendererProps<T extends BulletProps> = {
  bullet: T
  key?: React.Key
  selected?: boolean
}

export type BulletRenderer<T extends BulletProps> = (props: BulletRendererProps<T>) => React.ReactNode

function StandardBulletRenderer<T extends BulletProps>(props: BulletRendererProps<T>) {
  return (
    <Tag
      name={`${props.bullet.modifier == "not" ? "NOT " : ""}${props.bullet.name}`}
      color={props.bullet.color ?? "blue"}
      key={props.key}
      selected={props.selected}
    />
  )
}

const StandardNoValRenderer = () => <em>No Matches</em>

function filterBullets<T extends BulletProps>(values: Array<T>, filter: string): Array<T> {
  if (filter === "") {
    return values
  }
  filter = filter.toUpperCase()
  return values.filter(val =>
    val.name.toUpperCase().indexOf(filter) > -1
  )
}
