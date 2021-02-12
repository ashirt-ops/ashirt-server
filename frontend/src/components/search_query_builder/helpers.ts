// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dateFns from 'date-fns'

import { parseQuery, parseDateRangeString } from 'src/helpers'
import { Tag } from 'src/global_types'

export type SearchOptions = {
  text: string,
  sortAsc: boolean,
  uuid?: string,
  tags?: Array<Tag>,
  operator?: string,
  dateRange?: [Date, Date],
  hasLink?: boolean,
  withEvidenceUuid?: string,
}

const quoteText = (tagName: string) => tagName.includes(' ') ? `"${tagName}"` : tagName

const dateToRange = (dates: [Date, Date]) => {
  const fmt = (d: Date) => dateFns.format(d, 'yyyy-MM-dd')

  return `${fmt(dates[0])},${fmt(dates[1])}`
}

export const stringifySearch = (searchOpts: SearchOptions) => {
  return ([
    searchOpts.text,
    searchOpts.tags ? searchOpts.tags.map(tag => `tag:${quoteText(tag.name)}`).join(' ') : '',
    searchOpts.operator ? `operator:${searchOpts.operator}` : '',
    searchOpts.dateRange ? `range:${dateToRange(searchOpts.dateRange)}` : '',
    (searchOpts.hasLink != undefined) ? `linked:${searchOpts.hasLink}` : '',
    searchOpts.sortAsc ? 'sort:asc' : '',
    searchOpts.withEvidenceUuid ? `with-evidence:${searchOpts.withEvidenceUuid}` : '',
    searchOpts.uuid ? `uuid:${searchOpts.uuid}` : '',
  ])
    .filter(item => item != '') // remove the entries that aren't actually present
    .join(' ')
}

export const stringToSearch = (searchText: string, allTags: Array<Tag> = []) => {
  const tokens: { [key: string]: Array<string> } = parseQuery(searchText)

  const opts: SearchOptions = {
    text: '',
    sortAsc: false,
  }

  Object.entries(tokens).forEach(([key, values]) => {
    if (key == '') {
      opts.text = values.map(item => quoteText(item)).join(' ')
    }
    else if (key == 'tag') {
      opts.tags = values.
        map(tagName => allTags.find(tag => tag.name == tagName)).
        filter(isNotUndefined)
    }
    else if (key == 'operator') {
      opts.operator = values[0]
    }
    else if (key == 'range') {
      const range = parseDateRangeString(values[0])
      if (range) {
        opts.dateRange = range
      }
    }
    else if (key == 'linked') {
      const interpretedVal = values[0].toLowerCase().trim()
      if (interpretedVal == 'true' || interpretedVal == 'false') {
        opts.hasLink = (interpretedVal == 'true')
      }
    }
    else if (key == 'uuid') {
      opts.uuid = values[0]
    }
    else if (key == 'with-evidence') {
      opts.withEvidenceUuid = values[0]
    }
    else if (key == 'sort') {
      opts.sortAsc = ['asc', 'ascending', 'chronological'].includes(values[0])
    }
  })

  return opts
}

function isNotUndefined<T>(t: T | undefined): t is T {
  return t != undefined
}
