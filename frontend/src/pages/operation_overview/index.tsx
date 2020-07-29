// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import Button from 'src/components/button'
import { default as Tag, tagColorStyle } from 'src/components/tag'

import { RouteComponentProps } from 'react-router-dom'
import { getTagsByEvidenceUsage } from 'src/services'
import { useWiredData } from 'src/helpers'
import { differenceInCalendarDays, setHours, setMinutes, setSeconds } from 'date-fns'

import { default as Timeline, TimelineHeaders, DateHeader, SidebarHeader } from 'react-calendar-timeline'
import { ReactCalendarItemRendererProps, ReactCalendarGroupRendererProps } from 'react-calendar-timeline'
// make sure you include the timeline stylesheet or the timeline will not be styled
import './timeline.css'

const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{ slug: string }>) => {
  const { slug } = props.match.params
  const wiredTags = useWiredData(React.useCallback(() => getTagsByEvidenceUsage({ operationSlug: slug }), [slug]))

  return (
    <>
      <Button className={cx('back-button')} icon={require('./back.svg')} onClick={() => props.history.goBack()}>Back</Button>

      {wiredTags.render(tags => {
        const [firstDate, lastDate] = maxRange(tags.map(tag => tag.usages))
        const groups = tags.map(tag => ({
          id: tag.id,
          title: tag.name,
        }))

        let rangeCount = 0
        const items = tags.map((tag) => {
          const ranges = datesToRanges(tag.usages)
          const tagColors = tagColorStyle(tag.colorName)
          const rtn = ranges.map(([start, end], i) => ({
            id: rangeCount + i,
            group: tag.id,
            title: "",
            start_time: toStartOfDay(start),
            end_time: toEndOfDay(end),
            canChangeGroup: false,
            bgColor: tagColors.backgroundColor,
          }))
          rangeCount += rtn.length

          return rtn
        }).flat(1)

        const itemRender = (props: ReactCalendarItemRendererProps<any>) => {
          const borderColor = props.itemContext.selected ? "#880" : "#000"
          const borderWidth = props.itemContext.selected ? 3 : 1
          return (
            <div  {...props.getItemProps({
              style: {
                background: props.item.bgColor,
                borderColor,
                borderWidth
              },
            })}>
              <div>{props.itemContext.title} </div>
            </div>
          )
        }

        const timeChangeHandler = (visibleTimeStart: number, visibleTimeEnd: number, updateScrollCanvas: (s: number, e: number) => void) => {
          const [minTime, maxTime] = [firstDate, lastDate].map(date => date.getTime())
          if (visibleTimeStart < minTime && visibleTimeEnd > maxTime) {
            updateScrollCanvas(minTime, maxTime)
          }
          else if (visibleTimeStart < minTime) {
            updateScrollCanvas(minTime, minTime + (visibleTimeEnd - visibleTimeStart))
          }
          else if (visibleTimeEnd > maxTime) {
            updateScrollCanvas(maxTime - (visibleTimeEnd - visibleTimeStart), maxTime)
          }
          else {
            updateScrollCanvas(visibleTimeStart, visibleTimeEnd)
          }
        }

        const groupRenderer = (props: ReactCalendarGroupRendererProps<any>) => {
          const tag = tags.filter(someTag => someTag.id == props.group.id)[0]
          return <div style={{textAlign: "center"}}>
            <Tag name={tag.name} color={tag.colorName} />
            </div>
        }

        return (
          <Timeline
            groups={groups}
            items={items}
            defaultTimeStart={firstDate}
            defaultTimeEnd={lastDate}
            canMove={false}
            canResize={false}
            itemRenderer={itemRender}
            onTimeChange={timeChangeHandler}
            groupRenderer={groupRenderer}
            
          >
            <TimelineHeaders  >
              <DateHeader unit="primaryHeader" />
              <DateHeader />
            </TimelineHeaders>
          </Timeline>
        )
      })}
    </>
  )
}

const maxRange = (dates: Array<Array<Date>>) => {
  const sortedDates = dates.flat(1).sort()
  return [sortedDates[0], sortedDates[sortedDates.length - 1]]
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


const toStartOfDay = (day: Date) => setTime(day, 0, 0, 1)
const toEndOfDay = (day: Date) => setTime(day, 23, 59, 59)

const setTime = (day: Date, hour: number, minute: number, second: number) => setHours(setMinutes(setSeconds(day, second), minute), hour)
