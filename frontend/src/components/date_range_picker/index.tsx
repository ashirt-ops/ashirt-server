import 'react-day-picker/dist/style.css'
import * as React from 'react'
import { DayPicker } from 'react-day-picker'
import classnames from 'classnames/bind'
import { subDays, startOfMonth, endOfMonth } from 'date-fns'

import Button from 'src/components/button'
import Popover from 'src/components/popover'

import {
  DateRange,
  MaybeDateRange,
  stringifyRange,
  addDateToRange,
  lengthenRangeToDayBoundaries
} from './range_picker_helpers'

const cx = classnames.bind(require('./stylesheet'))

const now = new Date
const presetRanges: { [name: string]: DateRange } = {
  'Today': lengthenRangeToDayBoundaries([now, now]),
  'This Month': lengthenRangeToDayBoundaries([startOfMonth(now), endOfMonth(now)]),
  'Past 90 Days': lengthenRangeToDayBoundaries([subDays(now, 90), now]),
}

const DropDown = (props: {
  range: MaybeDateRange,
  onSelectRange: (newRange: MaybeDateRange) => void,
  onButtonClick: () => void
}) => (
  <div className={cx('popup')}>
    <div className={cx('shortcuts')}>
      <button onClick={() => props.onSelectRange(null)}>Clear Range</button>
      <br />
      {Object.keys(presetRanges).map(rangeName => (
        <button key={rangeName} onClick={() => props.onSelectRange(presetRanges[rangeName])}>{rangeName}</button>
      ))}
    </div>
    <div className={cx('day-picker-area')}>
      <DayPicker
        className={cx('day-picker')}
        mode='range'
        numberOfMonths={2}
        selected={props.range != null ? { from: props.range[0], to: props.range[1] } : undefined}
        onDayClick={d => props.onSelectRange(addDateToRange(d, props.range))}
        modifiers={props.range != null ? { start: props.range[0], end: props.range[1] } : undefined}
      />
      <Button primary className={cx('close-button')} onClick={props.onButtonClick} >Close</Button>
    </div>
  </div>
)

export default (props: {
  range: MaybeDateRange,
  onSelectRange: (r: MaybeDateRange) => void,
}) => {
  const [isOpen, setIsOpen] = React.useState(false)

  return (
    <Popover isOpen={isOpen} onClick={() => setIsOpen(true)} onRequestClose={() => setIsOpen(false)}
      content={<DropDown onButtonClick={() => setIsOpen(false)} {...props} />}>
      <Button doNotSubmit className={cx('open-button')} icon={require('./icon.svg')} title={stringifyRange(props.range)} />
    </Popover>
  )
}
