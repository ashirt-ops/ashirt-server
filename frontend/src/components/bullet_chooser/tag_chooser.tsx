import * as React from 'react'

import Tag from 'src/components/tag'
import { Tag as TagType } from 'src/global_types'
import { getTags, createTag } from 'src/services'
import { randomTagColorName, TagColor, shiftColor } from 'src/helpers/tag_colors'
import BulletChooser, { BulletProps } from 'src/components/bullet_chooser'
import { FilterModified } from 'src/helpers'
import { isNotUndefined } from 'src/helpers/is_not_undefined'
import classNames from 'classnames/bind'

const cx = classNames.bind(require('./stylesheet.styl'))

export const TagChooser = (props: {
  operationSlug: string,
  className?: string,
  disabled?: boolean,
  label: string,
  onChange: (tags: Array<TagType>) => void,
  value: Array<TagType>,
}) => {
  const [allTags, setAllTags] = React.useState<Array<TagType>>([])

  const reloadTags = () => { getTags({ operationSlug: props.operationSlug }).then(setAllTags) }
  React.useEffect(reloadTags, [props.operationSlug])

  const CreateNewTagElem = (props: { name: string }) => (
    <>
      Create new tag: <Tag name={props.name} color="" />
    </>
  )

  return (
    <BulletChooser
      label={props.label}
      options={allTags}
      value={props.value}
      onChange={props.onChange}
      valueRenderer={(b) => (
        <Tag name={b.bullet.name} color={b.bullet.colorName} {...b} />
      )}
      rowRenderer={(b) => (
        <>
          <Tag name={b.bullet.name} color={b.bullet.colorName} {...b} />
          {b.bullet.description && <p className={cx('tag-description')}>{b.bullet.description}</p>}
        </>
      )}
      noValueRenderer={(inputValue: string) => <CreateNewTagElem name={inputValue} />}
      onNoValueSelected={
        async (inputValue: string) => {
          const newVal = createTag({
            operationSlug: props.operationSlug,
            name: inputValue,
            colorName: randomTagColorName()
          })
          reloadTags()
          return newVal
        }}
    />
  )
}

/**
 * TagPicker is a specialized version of TagChooser. This allows for two features: first, it uses
 * a more-standard bulletProps. Second, it allows for "not" versions of tags. Additionally, this
 * does not allow for creating tags. This is primarily designed for the query builder. When used in
 * other scenarios, you should consider opting for TagChooser
 *
 * @param props
 * @returns
 */
export const TagPicker = (props: {
  operationSlug: string,
  className?: string,
  disabled?: boolean,
  label: string,
  onChange: (tags: Array<BulletProps>) => void,
  value: Array<BulletProps>,
  enableNot?: boolean
}) => {
  const [allTags, setAllTags] = React.useState<Array<TagType>>([])

  const reloadTags = () => { getTags({ operationSlug: props.operationSlug }).then(setAllTags) }
  React.useEffect(reloadTags, [props.operationSlug])

  return (
    <BulletChooser
      label={props.label}
      value={props.value}
      onChange={props.onChange}
      options={allTags.map(tagToBulletProps).filter(isNotUndefined)}
      enableNot={props.enableNot}
      valueRenderer={(b) => {
        const props = b.bullet.modifier == "not"
          ? { name: `NOT ${b.bullet.name}`, color: shiftColor(b.bullet.color ?? "disabledGray") }
          : { name: b.bullet.name, color: b.bullet.color ?? "disabledGray" }
        return (
          <Tag
            {...props}
            {...b}
          />
        )
      }}
    />
  )
}

export const tagToBulletProps = (tag: FilterModified<TagType> | undefined): BulletProps | undefined => {
  if (!tag) {
    return undefined
  }
  return {
    id: tag.id,
    name: tag.name,
    color: tag.colorName as TagColor,
    modifier: tag.modifier == 'not' ? "not" : undefined
  }
}

export const bulletPropsToTag = (val: BulletProps): TagType => {
  return {
    id: parseInt(`${val.id}`),
    name: val.name,
    colorName: val.color ?? "",
  }
}

export default TagChooser
