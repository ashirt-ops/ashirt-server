// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Tag from 'src/components/tag'
import WithLabel from 'src/components/with_label'
import classnames from 'classnames/bind'
import {Tag as TagType} from 'src/global_types'
import {dropRight} from 'lodash'
import {getTags, createTag} from 'src/services'
import {randomTagColorName} from 'src/helpers/tag_colors'
import PopoverMenu from 'src/components/popover_menu'
const cx = classnames.bind(require('./stylesheet'))

function filterTags(allTags: Array<TagType>, filter: string): Array<TagType> {
  if (filter === "") return allTags
  filter = filter.toUpperCase()
  return allTags.filter(tag =>
    tag.name.toUpperCase().indexOf(filter) > -1
   )
}

function getTagIdSet(tags: Array<TagType>): {[id: number]: true} {
  const tagIdSet: {[id: number]: true} = {}
  for (let t of tags) tagIdSet[t.id] = true
  return tagIdSet
}

export default (props: {
  operationSlug: string,
  className?: string,
  disabled?: boolean,
  label: string,
  onChange: (tags: Array<TagType>) => void,
  value: Array<TagType>,
}) => {
  const [allTags, setAllTags] = React.useState<Array<TagType>>([])
  const [inputValue, setInputValue] = React.useState("")
  const [dropdownVisible, setDropdownVisible] = React.useState(false)
  const [selectedTag, setSelectedTag] = React.useState<number>(-1)

  let filteredTags = filterTags(allTags, inputValue)
  const activeTagIdSet = getTagIdSet(props.value)

  const reloadTags = () => { getTags({operationSlug: props.operationSlug}).then(setAllTags) }
  React.useEffect(reloadTags, [props.operationSlug])

  const toggleTag = async (maybeTag: TagType|null) => {
    let tag: TagType
    if (maybeTag == null) {
      tag = await createTag({operationSlug: props.operationSlug, name: inputValue, colorName: randomTagColorName()})
      reloadTags()
    } else {
      tag = maybeTag
    }
    if (activeTagIdSet[tag.id]) props.onChange(props.value.filter(t => t.id !== tag.id))
    else props.onChange([...props.value, tag])
    setInputValue("")
  }

  const onInputKeyDown = (e: React.KeyboardEvent) => {
    if(inputValue === "") {
      if( ['Backspace', 'Delete'].includes(e.key)) {
        if (selectedTag != -1) {
          props.onChange([...props.value.slice(0, selectedTag), ...props.value.slice(selectedTag + 1)])
          setSelectedTag(selectedTag - 1)
        }
        else if (e.key != 'Delete'){
          props.onChange(dropRight(props.value))
        }
      }
      else if (e.key === 'ArrowLeft' && inputValue === "") {
        let index = selectedTag - 1
        setSelectedTag(index > -1 ? index : props.value.length - 1)
      }
      else if (e.key === 'ArrowRight' && inputValue === "") {
        setSelectedTag((selectedTag + 1) % props.value.length)
      }
      else {
        setSelectedTag(-1)
      }
    }
    setDropdownVisible(true)
  }

  const renderer = (maybeTag: TagType|null) => {
    if (maybeTag == null) {
      return <>Create new tag: <Tag name={inputValue} color="" /></>
    }
    return <Tag name={maybeTag.name} color={maybeTag.colorName} />
  }

  // If there are no filtered tags we add a null option to render the "Create new tag" option
  const options = filteredTags.length > 0 ? filteredTags : [null]

  return (
    <WithLabel label={props.label}>
      <PopoverMenu
        onRequestClose={() => setDropdownVisible(false)}
        isOpen={dropdownVisible}
        options={options}
        renderer={renderer}
        iconRenderer={t => t && activeTagIdSet[t.id] && require('./check.svg')}
        onSelect={toggleTag}
      >
        <div className={cx('input', props.className, {focus: dropdownVisible})}>
          {props.value.map((tag, i) => <Tag key={tag.id} selected={i == selectedTag} name={tag.name} color={tag.colorName} />)}
          <input
            onChange={e => setInputValue(e.target.value)}
            value={inputValue}
            onKeyDown={onInputKeyDown}
            onFocus={() => setDropdownVisible(true)}
          />
        </div>
      </PopoverMenu>
    </WithLabel>
  )
}
