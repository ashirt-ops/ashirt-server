import { computeDelta } from 'src/helpers'
import { default as req, reqMultipart, xhrText as reqText } from './request_helper'
import {Evidence, Finding, Tag, SubmittableEvidence, CodeBlock} from 'src/global_types'

export async function getEvidenceList(i: {
  operationSlug: string,
  query: string,
}): Promise<Array<Evidence>> {
  const evidence: Array<any> = await req('GET', `/operations/${i.operationSlug}/evidence`, null, {query: i.query})
  return evidence.map(evi => ({
    uuid: evi.uuid,
    description: evi.description,
    operator: evi.operator,
    occurredAt: new Date(evi.occurredAt),
    tags: evi.tags,
    contentType: evi.contentType
  }))
}

export async function getEvidenceAsCodeblock(i: {
  operationSlug: string,
  evidenceUuid: string,
}): Promise<CodeBlock> {
  const evi = await req('GET', `/operations/${i.operationSlug}/evidence/${i.evidenceUuid}/media`)
  return {
    type: 'codeblock',
    language: evi.contentSubtype,
    code: evi.content,
    source: evi.metadata ? evi.metadata['source'] : null,
  }
}

export async function getEvidenceAsString(i: {
  operationSlug: string,
  evidenceUuid: string,
}): Promise<string> {
  return await reqText('GET', `/operations/${i.operationSlug}/evidence/${i.evidenceUuid}/media`)
}

export async function createEvidence(i: {
  operationSlug: string,
  description: string,
  tagIds?: Array<number>,
  evidence: SubmittableEvidence,
}): Promise<void> {
  const formData = new FormData()
  formData.append('description', i.description)
  if (i.tagIds && i.tagIds.length > 0) {
    formData.append('tagIds', JSON.stringify(i.tagIds))
  }

  formData.append('contentType', i.evidence.type)

  if (i.evidence.type !== 'none') {
    formData.append('content', i.evidence.file)
  }

  await reqMultipart('POST', `/operations/${i.operationSlug}/evidence`, formData)
}

export async function updateEvidence(i: {
  operationSlug: string,
  evidenceUuid: string,
  description?: string,
  oldTags?: Array<Tag>,
  newTags?: Array<Tag>,
  updatedContent: Blob | null,
}): Promise<void> {
  const formData = new FormData()
  if (i.description !== undefined) {
    formData.append('description', i.description)
  }
  if (i.oldTags && i.newTags) {
    const [adds, subs] = computeDelta(i.oldTags.map(tag => tag.id), i.newTags.map(tag => tag.id))
    formData.append('tagsToAdd', JSON.stringify(adds))
    formData.append('tagsToRemove', JSON.stringify(subs))
  }

  if (i.updatedContent != null) {
    formData.append('content', i.updatedContent)
  }

  await reqMultipart('PUT', `/operations/${i.operationSlug}/evidence/${i.evidenceUuid}`, formData)
}

export async function changeFindingsOfEvidence(i: {
  operationSlug: string,
  evidenceUuid: string,
  oldFindings: Array<Finding>,
  newFindings: Array<Finding>,
}): Promise<void> {
  const [adds, subs] = computeDelta(i.oldFindings.map(f => f.uuid), i.newFindings.map(f => f.uuid))
  const updateFindings = (evidenceToAdd: Array<string>, evidenceToRemove: Array<string>) => (findingUuid: string) => (
    req('PUT', `/operations/${i.operationSlug}/findings/${findingUuid}/evidence`, {evidenceToAdd, evidenceToRemove})
  )
  await Promise.all(adds.map(updateFindings([i.evidenceUuid], [])))
  await Promise.all(subs.map(updateFindings([], [i.evidenceUuid])))
}

export async function deleteEvidence(i: {
  operationSlug: string,
  evidenceUuid: string,
  deleteAssociatedFindings: boolean,
}): Promise<void> {
  await req('DELETE', `/operations/${i.operationSlug}/evidence/${i.evidenceUuid}`, {
    deleteAssociatedFindings: i.deleteAssociatedFindings,
  })
}
