import { describe, it, expect } from 'vitest'
import {
  addTagToQuery,
  addOperatorToQuery,
  addOrUpdateDateRangeInQuery,
  getDateRangeFromQuery,
} from './query_parser'

describe('Query Parser Helpers', () => {
  describe('#addTagToQuery', () => {
    it('adds tags to empty queries', () => {
      expect(addTagToQuery('', 'MyTag')).toEqual('tag:MyTag')
    })
    it('adds tags to queries with text search', () => {
      expect(addTagToQuery('some "full text" search', 'TagToAdd')).toEqual(
        'some "full text" search tag:TagToAdd',
      )
    })
    it('appends tags if there are already tags in the query', () => {
      expect(addTagToQuery('some text tag:SomeTag', 'AnotherTag')).toEqual(
        'some text tag:SomeTag tag:AnotherTag',
      )
    })
    it('does not duplicate existing tags', () => {
      expect(addTagToQuery('tag:AlreadyAdded', 'AlreadyAdded')).toEqual('tag:AlreadyAdded')
    })
    it('handles tags with spaces', () => {
      expect(addTagToQuery('tag:"Existing Tag"', 'New Tag')).toEqual(
        'tag:"Existing Tag" tag:"New Tag"',
      )
    })
  })

  describe('#addOperatorToQuery', () => {
    it('adds an operator to queries', () => {
      expect(addOperatorToQuery('some query', 'alice')).toEqual('some query operator:alice')
    })

    it('replaces existing operator', () => {
      expect(addOperatorToQuery('query with operator operator:alice', 'bob')).toEqual(
        'query with operator operator:bob',
      )
    })
  })

  describe('#addOrUpdateDateRangeInQuery', () => {
    describe('removing ranges', () => {
      it('does nothing to queries without range', () => {
        expect(addOrUpdateDateRangeInQuery('some query "without range" tag:NoRange', null)).toEqual(
          'some query "without range" tag:NoRange',
        )
      })
      it('removes range from queries with a range', () => {
        expect(
          addOrUpdateDateRangeInQuery(
            'query with range range:2019-01-02,2019-05-09 tag:HasRange',
            null,
          ),
        ).toEqual('query with range tag:HasRange')
      })
    })

    describe('adding or updating ranges', () => {
      it('adds the range to queries without range', () => {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 23, 59, 59, 999)]
        expect(addOrUpdateDateRangeInQuery('query tag:NoRange', range)).toEqual(
          'query tag:NoRange range:2019-06-05,2019-08-09',
        )
      })
      it('updates the range for queries with a range', () => {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 23, 59, 59, 999)]
        expect(
          addOrUpdateDateRangeInQuery('query range:2019-01-01,2019-12-31 tag:HasRange', range),
        ).toEqual('query range:2019-06-05,2019-08-09 tag:HasRange')
      })
      it('adds time to range if the start of the range is not the start of day', () => {
        const range: [Date, Date] = [
          new Date(2019, 5, 5, 15, 30),
          new Date(2019, 7, 9, 23, 59, 59, 999),
        ]
        expect(addOrUpdateDateRangeInQuery('', range)).toEqual(
          `range:${new Date(2019, 5, 5, 15, 30).toISOString()},2019-08-09`,
        )
      })
      it('adds time to range if the end of the range is not the end of day', () => {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 15, 30)]
        expect(addOrUpdateDateRangeInQuery('', range)).toEqual(
          `range:2019-06-05,${new Date(2019, 7, 9, 15, 30).toISOString()}`,
        )
      })
    })
  })

  describe('#getDateRangeFromQuery', () => {
    it('returns null if there is no date range', () => {
      expect(getDateRangeFromQuery('some query "without range" tag:NoRange')).toBeNull()
    })
    it('returns a date range from a query with a range', () => {
      expect(getDateRangeFromQuery('my query range:2019-05-22,2019-07-26 tag:WithRange')).toEqual([
        new Date(2019, 4, 22),
        new Date(2019, 6, 26, 23, 59, 59, 999),
      ])
    })
    it('returns null if range is invalid', () => {
      expect(getDateRangeFromQuery('my query range:notarange')).toBeNull()
      expect(getDateRangeFromQuery('my query range:2019-05-05,NaN')).toBeNull()
    })
  })
})
