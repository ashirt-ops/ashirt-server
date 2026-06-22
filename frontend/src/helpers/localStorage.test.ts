import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { getIncludeDeletedUsers, setIncludeDeletedUsers } from './localStorage'

function makeLocalStorageMock() {
  const store = new Map<string, string>()
  return {
    getItem: (key: string) => store.get(key) ?? null,
    setItem: (key: string, value: string) => store.set(key, value),
    removeItem: (key: string) => store.delete(key),
    clear: () => store.clear(),
  }
}

beforeEach(() => {
  vi.stubGlobal('localStorage', makeLocalStorageMock())
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('getIncludeDeletedUsers / setIncludeDeletedUsers', () => {
  it('returns false by default', () => {
    expect(getIncludeDeletedUsers()).toBe(false)
  })

  it('returns true after setting to true', () => {
    setIncludeDeletedUsers(true)
    expect(getIncludeDeletedUsers()).toBe(true)
  })

  it('returns false after setting to false', () => {
    setIncludeDeletedUsers(true)
    setIncludeDeletedUsers(false)
    expect(getIncludeDeletedUsers()).toBe(false)
  })

  it('persists across multiple get calls', () => {
    setIncludeDeletedUsers(true)
    expect(getIncludeDeletedUsers()).toBe(true)
    expect(getIncludeDeletedUsers()).toBe(true)
  })

  it('can toggle back and forth', () => {
    setIncludeDeletedUsers(true)
    expect(getIncludeDeletedUsers()).toBe(true)
    setIncludeDeletedUsers(false)
    expect(getIncludeDeletedUsers()).toBe(false)
  })
})
