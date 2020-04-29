// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Evidence, Finding, Tag, SubmittableEvidence, CodeBlock } from 'src/global_types'
import { DataSource } from './data_sources/data_source'
import { computeDelta } from 'src/helpers'
import { evidenceFromDto } from './data_sources/converters'

export async function getEvidenceList(ds: DataSource, i: {
  operationSlug: string,
  query: string,
}): Promise<Array<Evidence>> {
  const evidence = await ds.listEvidence({ operationSlug: i.operationSlug }, i.query)
  return evidence.map(evidenceFromDto)
}

export async function getEvidenceAsCodeblock(ds: DataSource, i: {
  operationSlug: string,
  evidenceUuid: string,
}): Promise<CodeBlock> {
  const evi = JSON.parse(await ds.readEvidenceContent(i))
  return {
    type: 'codeblock',
    language: evi.contentSubtype,
    code: evi.content,
    source: evi.metadata ? evi.metadata['source'] : null,
  }
}

export async function getEvidenceAsString(ds: DataSource, i: {
  operationSlug: string,
  evidenceUuid: string,
}): Promise<string> {
  return await ds.readEvidenceContent(i)
}

export async function createEvidence(ds: DataSource, i: {
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

  await ds.createEvidence({ operationSlug: i.operationSlug }, formData)
}

export async function updateEvidence(ds: DataSource, i: {
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

  await ds.updateEvidence(
    { operationSlug: i.operationSlug, evidenceUuid: i.evidenceUuid },
    formData,
  )
}

export async function changeFindingsOfEvidence(ds: DataSource, i: {
  operationSlug: string,
  evidenceUuid: string,
  oldFindings: Array<Finding>,
  newFindings: Array<Finding>,
}): Promise<void> {
  const [adds, subs] = computeDelta(i.oldFindings.map(f => f.uuid), i.newFindings.map(f => f.uuid))

  await Promise.all(adds.map(findingUuid => ds.updateFindingEvidence(
    { operationSlug: i.operationSlug, findingUuid },
    { evidenceToAdd: [i.evidenceUuid], evidenceToRemove: [] },
  )))
  await Promise.all(subs.map(findingUuid => ds.updateFindingEvidence(
    { operationSlug: i.operationSlug, findingUuid },
    { evidenceToAdd: [], evidenceToRemove: [i.evidenceUuid] },
  )))
}

export async function deleteEvidence(ds: DataSource, i: {
  operationSlug: string,
  evidenceUuid: string,
  deleteAssociatedFindings: boolean,
}): Promise<void> {
  await ds.deleteEvidence(
    { operationSlug: i.operationSlug, evidenceUuid: i.evidenceUuid },
    { deleteAssociatedFindings: i.deleteAssociatedFindings },
  )
}
