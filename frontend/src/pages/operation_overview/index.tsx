// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import Button from 'src/components/button'
import { tagColorStyle } from 'src/components/tag'

import { RouteComponentProps } from 'react-router-dom'
import { getTagsByEvidenceUsage } from 'src/services'
import { useWiredData } from 'src/helpers'
import { differenceInCalendarDays, setHours, setMinutes, setSeconds, addDays } from 'date-fns'

import Timeline from 'react-calendar-timeline'
// make sure you include the timeline stylesheet or the timeline will not be styled
import 'react-calendar-timeline/lib/Timeline.css'

const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{ slug: string }>) => {
  const { slug } = props.match.params
  const wiredTags = useWiredData(React.useCallback(() => getTagsByEvidenceUsage({ operationSlug: slug }), [slug]))

  return (
    <>
      <Button className={cx('back-button')} icon={require('./back.svg')} onClick={() => props.history.goBack()}>Back</Button>

      {wiredTags.render(tags => {
        const groups = tags.map(tag => ({
          id: tag.id,
          title: tag.name,
        }))
        let rangeCount = 0
        const items = tags.map((tag) => {
          const ranges = datesToRanges(tag.usages)
          const rtn = ranges.map(([start, end], i) => ({
            id: rangeCount + i,
            group: tag.id,
            title: "",
            // start_time: toStartOfDay(start),
            // end_time: toEndOfDay(end),
            start_time: start,
            end_time: end,
            canChangeGroup: false,
            // className: cx("something")
            style: tagColorStyle(tag.colorName),
          }))
          rangeCount += rtn.length

          return rtn
        }).flat(1)

        return (<Timeline
          groups={groups}
          items={items}
          defaultTimeStart={addDays(new Date(), -14)}
          defaultTimeEnd={addDays(new Date(), 14)}
          canMove={false}
          canResize={false}
        // itemRenderer={ }
        />)
      })}
    </>
  )
}


const datesToRanges = (dates: Array<Date>) => {
  const ranges = []

  let start = null
  let nextEndDate = new Date()
  for (const date of dates) {
    if (start == null) {
      start = date
      nextEndDate = start
      continue
    }
    const diff = differenceInCalendarDays(date, nextEndDate)
    if (diff == 1) {
      nextEndDate = date
    }
    else {
      ranges.push([start, nextEndDate])
      start = date
      nextEndDate = start
    }
  }
  if (start != null) {
    ranges.push([start, nextEndDate])
  }
  return ranges
}


const toStartOfDay = (day: Date) => setHours(setMinutes(setSeconds(day, 0), 0), 0)
const toEndOfDay = (day: Date) => setHours(setMinutes(setSeconds(day, 59), 59), 23)
