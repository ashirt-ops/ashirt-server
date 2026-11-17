import { describe, it, expect } from 'vitest'
import { codeblockToBlob } from './codeblock_to_blob'

function readBlobAsText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsText(blob)
  })
}

describe('codeblockToBlob', () => {
  it('returns a Blob', () => {
    const blob = codeblockToBlob({
      type: 'codeblock',
      code: 'hello',
      language: 'plaintext',
      source: null,
    })
    expect(blob).toBeInstanceOf(Blob)
  })

  it('encodes code and language in the blob', async () => {
    const blob = codeblockToBlob({
      type: 'codeblock',
      code: 'console.log(1)',
      language: 'javascript',
      source: null,
    })
    const text = await readBlobAsText(blob)
    const parsed = JSON.parse(text)
    expect(parsed.content).toBe('console.log(1)')
    expect(parsed.contentSubtype).toBe('javascript')
  })

  it('includes source in metadata when provided', async () => {
    const blob = codeblockToBlob({
      type: 'codeblock',
      code: 'x',
      language: 'text',
      source: 'myfile.txt',
    })
    const text = await readBlobAsText(blob)
    const parsed = JSON.parse(text)
    expect(parsed.metadata).toEqual({ source: 'myfile.txt' })
  })

  it('omits metadata when source is not provided', async () => {
    const blob = codeblockToBlob({ type: 'codeblock', code: 'x', language: 'text', source: null })
    const text = await readBlobAsText(blob)
    const parsed = JSON.parse(text)
    expect(parsed.metadata).toBeUndefined()
  })

  it('omits metadata when source is null', async () => {
    const blob = codeblockToBlob({ type: 'codeblock', code: 'x', language: 'text', source: null })
    const text = await readBlobAsText(blob)
    const parsed = JSON.parse(text)
    expect(parsed.metadata).toBeUndefined()
  })
})
