// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import { Evidence, Finding } from 'src/global_types'
import { computeDelta } from 'src/helpers'

export async function getFindings(i: {
  operationSlug: string,
  query: string,
}): Promise<Array<Finding>> {
  const findings = await req('GET', `/operations/${i.operationSlug}/findings`, null, { query: i.query })

  return findings.map((finding: any) => ({
    ...finding,
    occurredFrom: finding.occurredFrom ? new Date(finding.occurredFrom) : null,
    occurredTo: finding.occurredTo ? new Date(finding.occurredTo) : null,
  }))
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
    req('GET', `/operations/${i.operationSlug}/findings/${i.findingUuid}`),
    req('GET', `/operations/${i.operationSlug}/findings/${i.findingUuid}/evidence`),
  ])

  return {
    finding: {
      ...finding,
      occurredFrom: finding.occurredFrom ? new Date(finding.occurredFrom) : null,
      occurredTo: finding.occurredTo ? new Date(finding.occurredTo) : null,
    },
    evidence: evidence.map((evi: any) => ({
      ...evi,
      occurredAt: new Date(evi.occurredAt),
    })),
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
  occurredAt: Date,
}): Promise<Finding> {
  return await req('POST', `/operations/${i.operationSlug}/findings`, {
    category: i.category,
    title: i.title,
    description: i.description,
  })
}

export async function changeEvidenceOfFinding(i: {
  operationSlug: string,
  findingUuid: string,
  oldEvidence: Array<Evidence>,
  newEvidence: Array<Evidence>,
}): Promise<void> {
  const [adds, subs] = computeDelta(i.oldEvidence.map(evi => evi.uuid), i.newEvidence.map(evi => evi.uuid))
  await req('PUT', `/operations/${i.operationSlug}/findings/${i.findingUuid}/evidence`, {
    evidenceToAdd: adds,
    evidenceToRemove: subs,
  })
}

export async function removeEvidenceFromFinding(i: {
  operationSlug: string,
  findingUuid: string,
  evidenceUuid: string,
}): Promise<void> {
  await req('PUT', `/operations/${i.operationSlug}/findings/${i.findingUuid}/evidence`, {
    evidenceToAdd: [],
    evidenceToRemove: [i.evidenceUuid],
  })
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
  await req('PUT', `/operations/${i.operationSlug}/findings/${i.findingUuid}`, {
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
  await req('DELETE', `/operations/${i.operationSlug}/findings/${i.findingUuid}`, {})
}
