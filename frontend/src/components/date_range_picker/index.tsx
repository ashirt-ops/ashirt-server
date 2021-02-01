// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import 'react-day-picker/lib/style.css'
import * as React from 'react'
import Button from 'src/components/button'
import DayPicker from 'react-day-picker'
import classnames from 'classnames/bind'
import {ClickPopover} from 'src/components/popover'
import {subDays, startOfMonth, endOfMonth} from 'date-fns'
import {DateRange, MaybeDateRange, stringifyRange, addDateToRange, lengthenRangeToDayBoundaries} from './range_picker_helpers'
const cx = classnames.bind(require('./stylesheet'))

const now = new Date
const presetRanges: {[name: string]: DateRange} = {
  'Today'        : lengthenRangeToDayBoundaries([now, now]),
  'This Month'   : lengthenRangeToDayBoundaries([startOfMonth(now), endOfMonth(now)]),
  'Past 90 Days' : lengthenRangeToDayBoundaries([subDays(now, 90), now]),
}

const DropDown = (props: {
  range: MaybeDateRange,
  onSelectRange: (newRange: MaybeDateRange) => void,
}) => (
  <div className={cx('popup')}>
    <div className={cx('shortcuts')}>
      <button onClick={() => props.onSelectRange(null)}>Clear Range</button>
      <br />
      {Object.keys(presetRanges).map(rangeName => (
        <button key={rangeName} onClick={() => props.onSelectRange(presetRanges[rangeName])}>{rangeName}</button>
      ))}
    </div>
    <DayPicker
      className={cx('day-picker')}
      numberOfMonths={2}
      selectedDays={props.range != null ? {from: props.range[0], to: props.range[1]} : undefined}
      onDayClick={d => props.onSelectRange(addDateToRange(d, props.range))}
      modifiers={props.range != null ? {start: props.range[0], end: props.range[1]} : undefined}
    />
  </div>
)

export default (props: {
  range: MaybeDateRange,
  onSelectRange: (r: MaybeDateRange) => void,
}) => (
  <ClickPopover content={<DropDown {...props} />}>
    <Button doNotSubmit className={cx('open-button')} icon={require('./icon.svg')} title={stringifyRange(props.range)} />
  </ClickPopover>
)
