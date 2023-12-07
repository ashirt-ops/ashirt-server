import {DateRange, addDateToRange, stringifyRange} from './range_picker_helpers'
import {expect} from 'chai'
import {getMonth, format} from 'date-fns'

describe('Range Picker Helpers', function() {
  describe('#addDateToRange', function() {
    describe('no range specified', function() {
      it('sets the range to a single day', function() {
        expect(addDateToRange(new Date(2019, 12, 21), null)).to.eql([
          new Date(2019, 12, 21, 0, 0, 0, 0),
          new Date(2019, 12, 21, 23, 59, 59, 999),
        ])
      })
    })
    describe('single day specified', function() {
      beforeEach(function() {
        this.range = [new Date(2019, 11, 11, 0, 0, 0, 0), new Date(2019, 11, 11, 23, 59, 59, 999)]
      })
      it('does nothing if the date is the same as the single day', function() {
        expect(addDateToRange(new Date(2019, 11, 11), this.range)).to.eql(this.range)
      })
      it('sets the end date if new date is after the single day', function() {
        expect(addDateToRange(new Date(2019, 12, 21), this.range)).to.eql([
          new Date(2019, 11, 11, 0, 0, 0, 0),
          new Date(2019, 12, 21, 23, 59, 59, 999),
        ])
      })
      it('sets the start date if new date is before the single day', function() {
        expect(addDateToRange(new Date(2019, 10, 15), this.range)).to.eql([
          new Date(2019, 10, 15, 0, 0, 0, 0),
          new Date(2019, 11, 11, 23, 59, 59, 999),
        ])
      })
    })
    describe('range specified', function() {
      it('replaces the range', function() {
      })
    })
  })

  describe('#stringifyRange', function() {
    it('stringifies null ranges', function() {
      expect(stringifyRange(null)).to.equal('Any Date')
    })
    it('stringifies ranges in the current year', function() {
      const currentYear = new Date().getFullYear()
      const currentMonth = getMonth(new Date())

      const range: DateRange = [new Date(currentYear, (currentMonth + 1) % 12, 5), new Date(currentYear, (currentMonth + 2) % 12, 25)]
      const monthNames = range.map(d => format(d, 'MMM'))

      expect(stringifyRange(range)).to.equal(`${monthNames[0]} 5th to ${monthNames[1]} 25th`)
    })
    it('stringifies single day ranges', function() {
      const range: DateRange = [new Date(2014, 6, 26), new Date(2014, 6, 26)]
      expect(stringifyRange(range)).to.equal('Jul 26th 2014')
    })
    it('stringifies ranges from previous years', function() {
      const range: DateRange = [new Date(2012, 1, 8), new Date(2012, 5, 19)]
      expect(stringifyRange(range)).to.equal('Feb 8th 2012 to Jun 19th 2012')
    })
    it('stringifies ranges that include this year', function() {
      const currentYear = new Date().getFullYear()
      const nextMonth = (getMonth(new Date()) + 1) % 12
      const fromDate = new Date(currentYear-1, 11, 15)
      const toDate = new Date(currentYear, nextMonth, 2)
      const range: DateRange = [fromDate, toDate]
      expect(stringifyRange(range)).to.equal(`Dec 15th ${currentYear-1} to ${format(toDate, 'MMM')} 2nd ${currentYear}`)
    })
  })
})
