import { describe, it, expect } from 'vitest'
import { clamp } from './clamp'

describe('clamp', () => {
  it('returns min when value is below min', () => {
    expect(clamp(-1, 0, 4)).toBe(0)
  })

  it('returns max when value is above max', () => {
    expect(clamp(5, 0, 4)).toBe(4)
  })

  it('returns value when within range', () => {
    expect(clamp(2, 0, 4)).toBe(2)
  })

  it('returns value when equal to min', () => {
    expect(clamp(0, 0, 4)).toBe(0)
  })

  it('returns value when equal to max', () => {
    expect(clamp(4, 0, 4)).toBe(4)
  })

  it('handles floating point values', () => {
    expect(clamp(2.5, 0, 4)).toBe(2.5)
  })

  it('handles negative ranges', () => {
    expect(clamp(-5, -10, -1)).toBe(-5)
  })

  it('clamps to min in a negative range', () => {
    expect(clamp(-15, -10, -1)).toBe(-10)
  })
})
