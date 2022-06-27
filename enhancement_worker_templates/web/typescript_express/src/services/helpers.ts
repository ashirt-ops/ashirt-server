import { randomFillSync } from "crypto"
import { FileData } from "./types"

function randomChars(length: number): string {
  const buff = Buffer.alloc(length)
  return randomFillSync(buff).toString('base64url')
}

export function encodeForm(fields: Record<string, string | undefined>, files: Record<string, FileData | undefined>) {
  const boundary = "----AShirtFormData-" + randomChars(30)
  const newline = "\r\n"
  const boundaryStart = "--" + boundary + newline
  const lastBoundary = "--" + boundary + "--" + newline

  let fieldBuffer = Buffer.from("")
  Object.entries(fields).forEach(([key, value]) => {
    if (value === undefined) {
      return
    }
    const text = boundaryStart +
      `Content-Disposition: form-data; name="${key}"` +
      newline + newline +
      value +
      newline
    fieldBuffer = Buffer.concat([fieldBuffer, Buffer.from(text)])
  })

  let fileBuffer = Buffer.from("")
  Object.entries(files).forEach(([key, fd]) => {
    if (fd === undefined) {
      return
    }
    const textPart = `${boundaryStart}` +
      `Content-Disposition: form-data; name="${key}"; filename="${fd.filename}"` +
      `${newline}Content-Type: ${fd.mimetype}` +
      `${newline}${newline}`

    fileBuffer = Buffer.concat([fileBuffer, Buffer.from(textPart), fd.content, Buffer.from(newline)])
  })

  return {
    boundary: boundary,
    data: Buffer.concat([
      fieldBuffer,
      fileBuffer,
      Buffer.from(lastBoundary),
    ])
  }
}
