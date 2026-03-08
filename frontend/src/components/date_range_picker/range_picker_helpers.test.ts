import { describe, it, expect, beforeEach } from 'vitest'
import { type DateRange, addDateToRange, stringifyRange } from './range_picker_helpers'
import { getMonth, format } from 'date-fns'

describe('Range Picker Helpers', () => {
  describe('#addDateToRange', () => {
    describe('no range specified', () => {
      it('sets the range to a single day', () => {
        expect(addDateToRange(new Date(2019, 12, 21), null)).toEqual([
          new Date(2019, 12, 21, 0, 0, 0, 0),
          new Date(2019, 12, 21, 23, 59, 59, 999),
        ])
      })
    })
    describe('single day specified', () => {
      let range: DateRange
      beforeEach(() => {
        range = [new Date(2019, 11, 11, 0, 0, 0, 0), new Date(2019, 11, 11, 23, 59, 59, 999)]
      })
      it('does nothing if the date is the same as the single day', () => {
        expect(addDateToRange(new Date(2019, 11, 11), range)).toEqual(range)
      })
      it('sets the end date if new date is after the single day', () => {
        expect(addDateToRange(new Date(2019, 12, 21), range)).toEqual([
          new Date(2019, 11, 11, 0, 0, 0, 0),
          new Date(2019, 12, 21, 23, 59, 59, 999),
        ])
      })
      it('sets the start date if new date is before the single day', () => {
        expect(addDateToRange(new Date(2019, 10, 15), range)).toEqual([
          new Date(2019, 10, 15, 0, 0, 0, 0),
          new Date(2019, 11, 11, 23, 59, 59, 999),
        ])
      })
    })
    describe('range specified', () => {
      it('replaces the range', () => {
        // placeholder — see addDateToRange implementation
      })
    })
  })

  describe('#stringifyRange', () => {
    it('stringifies null ranges', () => {
      expect(stringifyRange(null)).toEqual('Any Date')
    })
    it('stringifies ranges in the current year', () => {
      const currentYear = new Date().getFullYear()
      const currentMonth = getMonth(new Date())

      const range: DateRange = [
        new Date(currentYear, (currentMonth + 1) % 12, 5),
        new Date(currentYear, (currentMonth + 2) % 12, 25),
      ]
      const monthNames = range.map((d) => format(d, 'MMM'))

      expect(stringifyRange(range)).toEqual(`${monthNames[0]} 5th to ${monthNames[1]} 25th`)
    })
    it('stringifies single day ranges', () => {
      const range: DateRange = [new Date(2014, 6, 26), new Date(2014, 6, 26)]
      expect(stringifyRange(range)).toEqual('Jul 26th 2014')
    })
    it('stringifies ranges from previous years', () => {
      const range: DateRange = [new Date(2012, 1, 8), new Date(2012, 5, 19)]
      expect(stringifyRange(range)).toEqual('Feb 8th 2012 to Jun 19th 2012')
    })
    it('stringifies ranges that include this year', () => {
      const currentYear = new Date().getFullYear()
      const nextMonth = (getMonth(new Date()) + 1) % 12
      const fromDate = new Date(currentYear - 1, 11, 15)
      const toDate = new Date(currentYear, nextMonth, 2)
      const range: DateRange = [fromDate, toDate]
      expect(stringifyRange(range)).toEqual(
        `Dec 15th ${currentYear - 1} to ${format(toDate, 'MMM')} 2nd ${currentYear}`,
      )
    })
  })
})
