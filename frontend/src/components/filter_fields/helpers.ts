// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dateFns from 'date-fns'

import { parseQuery, parseDateRangeString, ParsedQuery, FilterModifier, FilterModified } from 'src/helpers'
import { Tag, User } from 'src/global_types'
import { BulletProps, creatorToBulletProps, supportedEvidenceCount, tagToBulletProps, textToBulletProps } from 'src/components/bullet_chooser'
import { isNotUndefined } from 'src/helpers/is_not_undefined'

export type SearchOptions = {
  text: string,
  meta?: Array<BulletProps>,
  sortAsc: boolean,
  uuid?: string,
  tags?: Array<BulletProps>,
  operator?: Array<BulletProps>,
  type?: Array<BulletProps>,
  dateRange?: [Date, Date],
  hasLink?: boolean,
  withEvidenceUuid?: Array<string>,
}

const quoteText = (tagName: string) => tagName.includes(' ') ? `"${tagName}"` : tagName

const dateToRange = (dates: [Date, Date]) => {
  const fmt = (d: Date) => dateFns.format(d, 'yyyy-MM-dd')

  return `${fmt(dates[0])},${fmt(dates[1])}`
}

const itemize = <T>(vals: Array<T> | undefined, tf: (v: T) => string): string => (
  vals ? vals.map(tf).join(' ') : ''
)

const itemizeBulletProps = (
  data: Array<BulletProps> | undefined, label: string, tf?: (v: BulletProps) => string
): string => {
  const orNot = (v: BulletProps): string => {
    const trueVal = tf?.(v) ?? v.id
    return (
      v.modifier == 'not'
        ? `!${trueVal}`
        : `${trueVal}`
    )
  }

  return itemize(data, (t) => `${label}:${orNot(t)}`)
}

export const stringifySearch = (searchOpts: SearchOptions) => {
  return ([
    searchOpts.text,
    itemizeBulletProps(searchOpts.meta, "meta", meta => quoteText(meta.name)),
    itemizeBulletProps(searchOpts.tags, "tag", tag => quoteText(tag.name)),
    itemizeBulletProps(searchOpts.operator, "operator"),
    searchOpts.dateRange ? `range:${dateToRange(searchOpts.dateRange)}` : '',
    (searchOpts.hasLink != undefined) ? `linked:${searchOpts.hasLink}` : '',
    searchOpts.sortAsc ? 'sort:asc' : '',
    itemizeBulletProps(searchOpts.type, "type"),
    itemize(searchOpts.withEvidenceUuid, (evi => `with-evidence:${evi}`)),
    searchOpts.uuid ? `uuid:${searchOpts.uuid}` : '',
  ])
    .filter(item => item != '') // remove the entries that aren't actually present
    .join(' ')
}

const findAndModify = <T>(
  values: Array<T>,
  finder: (v: T) => boolean,
  modifier?: FilterModifier
): (FilterModified<T>) | undefined => {
  const result = values.find(finder)
  return result ? { ...result, modifier } : undefined
}

export const stringToSearch = (
  searchText: string,
  allTags: Array<Tag> = [],
  allCreators: Array<User> = []
) => {
  const tokens: ParsedQuery = parseQuery(searchText)

  const opts: SearchOptions = {
    text: '',
    sortAsc: false,
  }

  Object.entries(tokens).forEach(([key, filterValues]) => {
    if (key == '') {
      opts.text = filterValues.map(item => quoteText(item.value)).join(' ')
    }
    else if (key == 'tag') {
      opts.tags = filterValues
        .map(fVal => findAndModify(allTags, (tag => tag.name == fVal.value), fVal.modifier))
        .map(tagToBulletProps)
        .filter(isNotUndefined)
    }
    else if (key == 'meta') {
      opts.meta = filterValues
        .map(fv => fv.value)
        .map(textToBulletProps)
        .filter(isNotUndefined)

    }
    else if (key == 'operator') {
      opts.operator = filterValues
        .map(fVal => findAndModify(allCreators, (c => c.slug == fVal.value), fVal.modifier))
        .map(creatorToBulletProps)
        .filter(isNotUndefined)
    }
    else if (key == 'range') {
      const range = parseDateRangeString(filterValues[0].value)
      if (range) {
        opts.dateRange = range
      }
    }
    else if (key == 'type') {
      opts.type = filterValues
        .map(fVal => findAndModify(supportedEvidenceCount, (t => t.id == fVal.value), fVal.modifier))
        .filter(isNotUndefined)
    }
    else if (key == 'linked') {
      const interpretedVal = filterValues[0].value.toLowerCase().trim()
      if (interpretedVal == 'true' || interpretedVal == 'false') {
        opts.hasLink = (interpretedVal == 'true')
      }
    }
    else if (key == 'uuid') {
      opts.uuid = filterValues[0].value
    }
    else if (key == 'with-evidence') {
      opts.withEvidenceUuid = filterValues.map(v => v.value)
    }
    else if (key == 'sort') {
      opts.sortAsc = ['asc', 'ascending', 'chronological'].includes(filterValues[0].value)
    }
  })

  return opts
}
