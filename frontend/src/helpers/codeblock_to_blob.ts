import { CodeBlock } from 'src/global_types'

type JsonEvidence = {
  contentSubtype: string,
  content: string,
  metadata?: { [key: string]: string }
}

export const codeblockToBlob = (cb: CodeBlock): Blob => {
  const evidence: JsonEvidence = {
    content: cb.code,
    contentSubtype: cb.language,
  }
  if (cb.source != null) {
    evidence.metadata = { 'source': cb.source }
  }
  return new Blob([JSON.stringify(evidence)])
}
