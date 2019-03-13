// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Evidence, Finding } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'
import { computeDelta } from 'src/helpers'
import { findingFromDto, evidenceFromDto } from './data_sources/converters'

export async function getFindings(i: {
  operationSlug: string,
  query: string,
}): Promise<Array<Finding>> {
  const findings = await ds.listFindings({ operationSlug: i.operationSlug }, i.query)
  return findings.map(findingFromDto)
}

export async function getFindingsOfEvidence(i: {
  operationSlug: string,
  evidenceUuid: string,
}): Promise<Array<Finding>> {
  return getFindings({
    operationSlug: i.operationSlug,
    query: `with-evidence:${JSON.stringify(i.evidenceUuid)}`,
  })
}

export async function getFinding(i: {
  operationSlug: string,
  findingUuid: string,
}): Promise<{ finding: Finding, evidence: Array<Evidence> }> {
  const [finding, evidence] = await Promise.all([
    ds.readFinding(i),
    ds.readFindingEvidence(i),
  ])
  return {
    finding: findingFromDto(finding),
    evidence: evidence.map(evidenceFromDto),
  }
}

export async function getFindingCategories(): Promise<Array<string>> {
  return [
    'Product',
    'Network',
    'Enterprise',
    'Vendor',
    'Behavioral',
    'Detection Gap',
  ]
}

export async function createFinding(i: {
  operationSlug: string,
  category: string,
  title: string,
  description: string,
}): Promise<Finding> {
  const finding = await ds.createFinding({ operationSlug: i.operationSlug }, {
    category: i.category,
    title: i.title,
    description: i.description,
  })
  return findingFromDto(finding)
}

export async function changeEvidenceOfFinding(i: {
  operationSlug: string,
  findingUuid: string,
  oldEvidence: Array<Evidence>,
  newEvidence: Array<Evidence>,
}): Promise<void> {
  const [adds, subs] = computeDelta(i.oldEvidence.map(evi => evi.uuid), i.newEvidence.map(evi => evi.uuid))
  await ds.updateFindingEvidence(
    { operationSlug: i.operationSlug, findingUuid: i.findingUuid },
    { evidenceToAdd: adds, evidenceToRemove: subs },
  )
}

export async function removeEvidenceFromFinding(i: {
  operationSlug: string,
  findingUuid: string,
  evidenceUuid: string,
}): Promise<void> {
  await ds.updateFindingEvidence(
    { operationSlug: i.operationSlug, findingUuid: i.findingUuid },
    { evidenceToAdd: [], evidenceToRemove: [i.evidenceUuid] },
  )
}

export async function updateFinding(i: {
  findingUuid: string,
  operationSlug: string,
  category: string,
  title: string,
  description: string,
  readyToReport: boolean,
  ticketLink: string | null,
}): Promise<void> {
  await ds.updateFinding({ operationSlug: i.operationSlug, findingUuid: i.findingUuid }, {
    category: i.category,
    title: i.title,
    description: i.description,
    readyToReport: i.readyToReport,
    ticketLink: i.ticketLink,
  })
}

export async function deleteFinding(i: {
  findingUuid: string,
  operationSlug: string,
}): Promise<void> {
  await ds.deleteFinding(i)
}
