import { describe, it, expect } from 'vitest'
import { tagColorNames, tagColorNameToColor, shiftColor, randomTagColorName } from './tag_colors'

describe('tagColorNames', () => {
  it('contains expected colors', () => {
    expect(tagColorNames).toContain('blue')
    expect(tagColorNames).toContain('red')
    expect(tagColorNames).toContain('disabledGray')
  })

  it('contains light variants', () => {
    expect(tagColorNames).toContain('lightBlue')
    expect(tagColorNames).toContain('lightRed')
  })

  it('has 21 colors', () => {
    expect(tagColorNames).toHaveLength(21)
  })
})

describe('tagColorNameToColor', () => {
  it('returns the correct hex value for blue', () => {
    expect(tagColorNameToColor('blue')).toBe(0x0e5a8a)
  })

  it('returns the correct hex value for red', () => {
    expect(tagColorNameToColor('red')).toBe(0xa82a2a)
  })

  it('returns 0x000000 for an unknown color', () => {
    expect(tagColorNameToColor('notacolor')).toBe(0x000000)
  })
})

describe('shiftColor', () => {
  it('shifts blue to lightBlue', () => {
    expect(shiftColor('blue')).toBe('lightBlue')
  })

  it('shifts lightBlue back to blue', () => {
    expect(shiftColor('lightBlue')).toBe('blue')
  })

  it('shifts red to lightRed', () => {
    expect(shiftColor('red')).toBe('lightRed')
  })

  it('disabledGray shifts to itself', () => {
    expect(shiftColor('disabledGray')).toBe('disabledGray')
  })
})

describe('randomTagColorName', () => {
  it('returns a string in tagColorNames', () => {
    const result = randomTagColorName()
    expect(tagColorNames).toContain(result)
  })

  it('returns different values over multiple calls', () => {
    const results = new Set(Array.from({ length: 50 }, () => randomTagColorName()))
    expect(results.size).toBeGreaterThan(1)
  })
})
