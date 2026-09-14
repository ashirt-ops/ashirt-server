import { describe, it, expect, vi, beforeEach } from 'vitest'
import xhr, { xhrText, reqMultipart } from './request_helper'

function mockFetch(status: number, body: unknown) {
  const isText = typeof body === 'string'
  const response = {
    status,
    json: () => Promise.resolve(isText ? {} : body),
    text: () => Promise.resolve(isText ? body : JSON.stringify(body)),
  } as unknown as Response
  vi.mocked(globalThis.fetch).mockResolvedValue(response)
}

beforeEach(() => {
  vi.spyOn(globalThis, 'fetch')
})

describe('xhr', () => {
  it('makes a GET request to /web + path', async () => {
    mockFetch(200, { ok: true })
    await xhr('GET', '/api/data')
    expect(fetch).toHaveBeenCalledWith('/web/api/data', { method: 'GET' })
  })

  it('returns parsed JSON on success', async () => {
    mockFetch(200, { value: 42 })
    const result = await xhr('GET', '/api/data')
    expect(result).toEqual({ value: 42 })
  })

  it('makes a POST request with JSON body', async () => {
    mockFetch(201, {})
    await xhr('POST', '/api/data', { name: 'Alice' })
    expect(fetch).toHaveBeenCalledWith('/web/api/data', {
      method: 'POST',
      body: JSON.stringify({ name: 'Alice' }),
      headers: { 'Content-Type': 'application/json' },
    })
  })

  it('appends query string parameters', async () => {
    mockFetch(200, {})
    await xhr('GET', '/api/data', null, { page: '2' })
    const [url] = vi.mocked(fetch).mock.calls[0]
    expect(url).toContain('page=2')
  })

  it('throws HttpError on 4xx responses', async () => {
    mockFetch(404, {})
    await expect(xhr('GET', '/api/missing')).rejects.toThrow()
  })

  it('throws HttpError on 5xx responses', async () => {
    mockFetch(500, {})
    await expect(xhr('GET', '/api/error')).rejects.toThrow()
  })

  it('throws HttpError with body.error message when present', async () => {
    mockFetch(400, { error: 'bad request' })
    await expect(xhr('GET', '/api/bad')).rejects.toThrow('bad request')
  })

  it('makes a PUT request', async () => {
    mockFetch(200, {})
    await xhr('PUT', '/api/data', { updated: true })
    expect(fetch).toHaveBeenCalledWith('/web/api/data', expect.objectContaining({ method: 'PUT' }))
  })
})

describe('xhrText', () => {
  it('returns text response', async () => {
    mockFetch(200, 'hello text')
    const result = await xhrText('GET', '/api/text')
    expect(result).toBe('hello text')
  })
})

describe('reqMultipart', () => {
  it('makes a multipart request and returns JSON', async () => {
    mockFetch(200, { id: 1 })
    const formData = new FormData()
    const result = await reqMultipart('POST', '/api/upload', formData)
    expect(result).toEqual({ id: 1 })
    expect(fetch).toHaveBeenCalledWith('/web/api/upload', {
      method: 'POST',
      body: formData,
    })
  })

  it('throws on error status in multipart response', async () => {
    mockFetch(500, { error: 'upload failed' })
    await expect(reqMultipart('POST', '/api/upload', new FormData())).rejects.toThrow()
  })
})
