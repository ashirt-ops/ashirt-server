import {computeDelta} from './compute_delta'
import {expect} from 'chai'

describe('#computeDelta', function() {
  it('marks elements found in after but not before as additions', function() {
    const [additions, _subtractions] = computeDelta([1, 2, 3, 5, 8], [2, 4, 6, 8, 10])
    expect(additions).to.eql([4, 6, 10])
  })

  it('marks elements found in before but not after as subtractions', function() {
    const [_additions, subtractions] = computeDelta([1, 2, 3, 5, 8], [2, 4, 6, 8, 10])
    expect(subtractions).to.eql([1, 3, 5])
  })

  it('handles empty case', function() {
    expect(computeDelta([], [])).to.eql([[], []])
  })
})
