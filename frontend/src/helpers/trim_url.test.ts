import { describe, it, expect } from 'vitest'
import { trimURL } from './trim_url'

describe('trimURL', () => {
  it('detects a URL', () => {
    expect(trimURL('https://github.com/owner/repo').isAUrl).toBe(true)
  })

  it('detects a non-URL string', () => {
    expect(trimURL('just some text').isAUrl).toBe(false)
  })

  it('returns the string unchanged when it is short enough', () => {
    expect(trimURL('short text').trimmedValue).toBe('short text')
  })

  it('truncates a long non-URL string', () => {
    const long = 'a'.repeat(60)
    const { trimmedValue } = trimURL(long)
    expect(trimmedValue.endsWith('...')).toBe(true)
    expect(trimmedValue.length).toBeLessThan(60)
  })

  it('elides middle path segments for deep URLs', () => {
    const { trimmedValue } = trimURL('https://github.com/microsoft/vscode/blob/master/README.md')
    expect(trimmedValue).toContain('...')
    expect(trimmedValue).toContain('README.md')
    expect(trimmedValue).toContain('github.com')
  })

  it('keeps short URL paths intact', () => {
    const { trimmedValue } = trimURL('https://github.com/owner/repo')
    expect(trimmedValue).not.toContain('...')
    expect(trimmedValue).toBe('github.com/owner/repo')
  })

  it('strips query string from URL', () => {
    const { trimmedValue } = trimURL('https://example.com/path?foo=bar')
    expect(trimmedValue).not.toContain('?')
    expect(trimmedValue).not.toContain('foo')
  })

  it('strips fragment from the last path segment', () => {
    const { trimmedValue } = trimURL('https://example.com/page#section')
    expect(trimmedValue).not.toContain('#')
  })
})
