import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import { highlightSubstring } from './highlight_substring'

describe('highlightSubstring', () => {
  it('returns an array with all text when no match is found', () => {
    const result = highlightSubstring('hello world', 'xyz', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.textContent).toBe('hello world')
  })

  it('returns a single span for an empty string', () => {
    const result = highlightSubstring('', 'abc', 'hl')
    expect(result).toHaveLength(1)
  })

  it('respects minLength option and skips highlighting when query is too short', () => {
    const result = highlightSubstring('hello', 'h', 'hl', { minLength: 2 })
    const { container } = render(<div>{result}</div>)
    expect(container.textContent).toBe('hello')
    expect(container.querySelectorAll('.hl')).toHaveLength(0)
  })

  it('highlights a single match in the middle', () => {
    const result = highlightSubstring('hello world', 'world', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.querySelector('.hl')?.textContent).toBe('world')
  })

  it('highlights multiple matches', () => {
    const result = highlightSubstring('the cat and the dog', 'the', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.querySelectorAll('.hl')).toHaveLength(2)
  })

  it('highlights a match at the start of the string', () => {
    const result = highlightSubstring('hello world', 'hello', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.querySelector('.hl')?.textContent).toBe('hello')
  })

  it('highlights a match at the end of the string', () => {
    const result = highlightSubstring('hello world', 'world', 'hl')
    const { container } = render(<div>{result}</div>)
    const spans = container.querySelectorAll('.hl')
    expect(spans[spans.length - 1].textContent).toBe('world')
  })

  it('is case-sensitive by default', () => {
    const result = highlightSubstring('Hello World', 'hello', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.querySelectorAll('.hl')).toHaveLength(0)
  })

  it('supports case-insensitive matching via regexFlags', () => {
    const result = highlightSubstring('Hello World', 'hello', 'hl', { regexFlags: 'i' })
    const { container } = render(<div>{result}</div>)
    expect(container.querySelectorAll('.hl')).toHaveLength(1)
  })

  it('escapes regex special characters in the query', () => {
    const result = highlightSubstring('1+1=2', '1+1', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.querySelector('.hl')?.textContent).toBe('1+1')
  })

  it('preserves non-matching text around matches', () => {
    const result = highlightSubstring('abc def ghi', 'def', 'hl')
    const { container } = render(<div>{result}</div>)
    expect(container.textContent).toBe('abc def ghi')
  })
})
