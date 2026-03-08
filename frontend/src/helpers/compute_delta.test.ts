import { describe, it, expect } from 'vitest'
import { computeDelta } from './compute_delta'

describe('#computeDelta', () => {
  it('marks elements found in after but not before as additions', () => {
    const [additions, _subtractions] = computeDelta([1, 2, 3, 5, 8], [2, 4, 6, 8, 10])
    expect(additions).toEqual([4, 6, 10])
  })

  it('marks elements found in before but not after as subtractions', () => {
    const [_additions, subtractions] = computeDelta([1, 2, 3, 5, 8], [2, 4, 6, 8, 10])
    expect(subtractions).toEqual([1, 3, 5])
  })

  it('handles empty case', () => {
    expect(computeDelta([], [])).toEqual([[], []])
  })
})
