import { type MouseEvent } from 'react'
import Tag from 'src/components/tag'
import classnames from 'classnames/bind'
import { type Tag as TagType } from 'src/global_types'
const cx = classnames.bind(require('./stylesheet'))

const TagList = (props: {
  tags: Array<TagType>
  onTagClick?: (t: TagType, e: MouseEvent) => void
}) => (
  <div className={cx('root')}>
    {props.tags.map((tag) => (
      <Tag
        key={tag.id}
        name={tag.name}
        onClick={(e) => props.onTagClick?.(tag, e) ?? null}
        color={tag.colorName}
      />
    ))}
  </div>
)
export default TagList
