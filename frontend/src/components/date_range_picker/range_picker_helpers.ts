// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dateFns from 'date-fns'

export type DateRange = [Date, Date]
export type MaybeDateRange = DateRange | null

function dateToString(d: Date, includeYear: boolean): string {
  if (dateFns.isToday(d)) return 'Today'
  if (dateFns.isYesterday(d)) return 'Yesterday'
  return dateFns.format(d, includeYear ? 'MMM do yyyy' : 'MMM do')
}

function singleDaySelected(r: DateRange): boolean {
  return dateFns.isSameDay(r[0], r[1])
}

export function addDateToRange(d: Date, r: MaybeDateRange): DateRange {
  if (r == null || !singleDaySelected(r)) {
    return lengthenRangeToDayBoundaries([d, d])
  }
  if (dateFns.isBefore(r[0], d)) {
    return lengthenRangeToDayBoundaries([r[0], d])
  } else {
    return lengthenRangeToDayBoundaries([d, r[0]])
  }
}

export function lengthenRangeToDayBoundaries(r: DateRange): DateRange {
  return [dateFns.startOfDay(r[0]), dateFns.endOfDay(r[1])]
}

export function stringifyRange(r: MaybeDateRange): string {
  if (r == null) return 'Any Date'
  if (singleDaySelected(r)) return dateToString(r[0], !dateFns.isThisYear(r[0]))
  const includeYear = !dateFns.isSameYear(r[0], r[1]) || !dateFns.isThisYear(r[0])
  return `${dateToString(r[0], includeYear)} to ${dateToString(r[1], includeYear)}`
}
