import { AShirtService } from "src/ashirt"
import { ProcessRequest } from "src/helpers/request_validation"
import { default as axios, AxiosError } from 'axios'

export type ProcessResultDTO =
  | ProcessResultNormal
  | ProcessResultComplete
  | ProcessResultDeferred

type ProcessResultNormal = {
  action: "rejected" | "error"
  content?: string
}

type ProcessResultComplete = {
  action: "processed"
  content: string
}

type ProcessResultDeferred = {
  action: "deferred"
}

export const handleActionProcess = async (
  body: ProcessRequest,
  svc: AShirtService
): Promise<ProcessResultDTO> => {

  if (body.contentType !== 'image') {
    return {
      action: 'rejected',
    }
  }

  try {
    const resp = await svc.getEvidenceContent(body.operationSlug, body.evidenceUuid)
    const content = convertToPrettyHex(resp.data)

    return {
      action: 'processed',
      content
    }
  }
  catch(err: unknown) {
    const content = axios.isAxiosError(err)
      ? (err as AxiosError).message
      : null

    const rtn: ProcessResultNormal = {
      action: 'error',
      ...(content ? { content } : {}),
    }
    return rtn
  }
}

const convertToPrettyHex = (data: Buffer): string => {
  const hexData = data.toString('hex')

  const lines = chunkSubstr(hexData, 32)

  // convert lines from abc123 to AB C1 23, combine lines with a newline seperator

  const toPrettyLine = (line: string) => {
    const chars = chunkSubstr(line, 2).map(chunk => chunk.toUpperCase())
    const interpretedChars = chars.map(v => {
      const val = parseInt(v, 16)
      if (val > 0x32 && val < 0x7F) { // exclude space, non-printable, del, and extended ascii
        return String.fromCharCode(val)
      }
      return '.'
    })

    return `${chars.join(' ')}    ${interpretedChars.join('')}`
  }

  const content = lines
    .map(toPrettyLine)
    .join('\n')
  return content
}

// from https://stackoverflow.com/a/29202760
function chunkSubstr(str: string, size: number): Array<string> {
  const numChunks = Math.ceil(str.length / size)
  const chunks = new Array(numChunks)

  for (let i = 0, o = 0; i < numChunks; ++i, o += size) {
    chunks[i] = str.substr(o, size)
  }

  return chunks
}
