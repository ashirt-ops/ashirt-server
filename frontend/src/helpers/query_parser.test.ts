import {addTagToQuery, addOperatorToQuery, addOrUpdateDateRangeInQuery, getDateRangeFromQuery} from './query_parser'
import {expect} from 'chai'

describe('Query Parser Helpers', function() {
  describe('#addTagToQuery', function() {
    it('adds tags to empty queries', function() {
      expect(addTagToQuery('', 'MyTag')).to.equal('tag:MyTag')
    })
    it('adds tags to queries with text search', function() {
      expect(addTagToQuery('some "full text" search', 'TagToAdd')).to.equal('some "full text" search tag:TagToAdd')
    })
    it('appends tags if there are already tags in the query', function() {
      expect(addTagToQuery('some text tag:SomeTag', 'AnotherTag')).to.equal('some text tag:SomeTag tag:AnotherTag')
    })
    it('does not duplicate existing tags', function() {
      expect(addTagToQuery('tag:AlreadyAdded', 'AlreadyAdded')).to.equal('tag:AlreadyAdded')
    })
    it('handles tags with spaces', function() {
      expect(addTagToQuery('tag:"Existing Tag"', 'New Tag')).to.equal('tag:"Existing Tag" tag:"New Tag"')
    })
  })

  describe('#addOperatorToQuery', function() {
    it('adds an operator to queries', function() {
      expect(addOperatorToQuery('some query', 'alice')).to.equal('some query operator:alice')
    })

    it('replaces existing operator', function() {
      expect(addOperatorToQuery('query with operator operator:alice', 'bob')).to.equal('query with operator operator:bob')
    })
  })

  describe('#addOrUpdateDateRangeInQuery', function() {
    describe('removing ranges', function() {
      it('does nothing to queries without range', function() {
        expect(addOrUpdateDateRangeInQuery('some query "without range" tag:NoRange', null)).to.equal('some query "without range" tag:NoRange')
      })
      it('removes range from queries with a range', function() {
        expect(addOrUpdateDateRangeInQuery('query with range range:2019-01-02,2019-05-09 tag:HasRange', null)).to.equal('query with range tag:HasRange')
      })
    })

    describe('adding or updating ranges', function() {
      it('adds the range to queries without range', function() {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 23, 59, 59, 999)]
        expect(addOrUpdateDateRangeInQuery('query tag:NoRange', range)).to.equal('query tag:NoRange range:2019-06-05,2019-08-09')
      })
      it('updates the range for queryies with a range', function() {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 23, 59, 59, 999)]
        expect(addOrUpdateDateRangeInQuery('query range:2019-01-01,2019-12-31 tag:HasRange', range)).to.equal('query range:2019-06-05,2019-08-09 tag:HasRange')
      })
      it('adds time to range if the start of the range is not the start of day', function() {
        const range: [Date, Date] = [new Date(2019, 5, 5, 15, 30), new Date(2019, 7, 9, 23, 59, 59, 999)]
        expect(addOrUpdateDateRangeInQuery('', range)).to.equal(`range:${new Date(2019, 5, 5, 15, 30).toISOString()},2019-08-09`)
      })
      it('adds time to range if the end of the range is not the end of day', function() {
        const range: [Date, Date] = [new Date(2019, 5, 5), new Date(2019, 7, 9, 15, 30)]
        expect(addOrUpdateDateRangeInQuery('', range)).to.equal(`range:2019-06-05,${new Date(2019, 7, 9, 15, 30).toISOString()}`)
      })
    })
  })

  describe('#getDateRangeFromQuery', function() {
    it('returns null if there is no date range', function() {
      expect(getDateRangeFromQuery('some query "without range" tag:NoRange')).to.be.null
    })
    it('returns a date range from a query with a range', function() {
      expect(getDateRangeFromQuery('my query range:2019-05-22,2019-07-26 tag:WithRange')).to.eql([new Date(2019, 4, 22), new Date(2019, 6, 26, 23, 59, 59, 999)])
    })
    it('returns null if range is invalid', function() {
      expect(getDateRangeFromQuery('my query range:notarange')).to.be.null
      expect(getDateRangeFromQuery('my query range:2019-05-05,NaN')).to.be.null
    })
  })
})
