// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dtos from '../dtos'
import * as types from 'src/global_types'
import { DataSource } from '../data_source'
import { fetchJsonp } from 'src/helpers/fetch_jsonp'
import { minBy, maxBy } from 'lodash'

export type OperationArchiveData = {
  id: number,
  slug: string
  name: string,
  status: number,
  tags: Array<{
    id: number,
    operationId: number,
    name: string
    colorName: string,
  }>,
  evidence: Array<{
    id: number,
    uuid: string,
    operationId: number,
    operatorId: number,
    description: string,
    contentType: string,
    fullImageKey: string,
    thumbImageKey: string,
    occurredAt: string,
  }>,
  findings: Array<{
    id: number,
    uuid: string,
    operationId: number,
    readyToReport: boolean,
    ticketLink: string | undefined,
    category: string,
    title: string,
    description: string,
  }>,
  evidence_finding_map: Array<{
    evidenceId: number,
    findingId: number,
  }>,
  tag_evidence_map: Array<{
    tagId: number,
    evidenceId: number,
  }>,
  users: Array<{
    id: number,
    slug: string,
    firstName: string,
    lastName: string,
    email: string,
    admin: boolean,
    disabled: boolean,
    headless: boolean,
  }>,
}

const archiveUser = {
  firstName: 'Archive',
  lastName: 'Viewer',
  slug: '',
  email: '',
  admin: false,
  authSchemes: [],
  headless: false,
}

export function makeArchiveDataSource(data: OperationArchiveData): DataSource {
  const operation = {
    id: data.id,
    slug: data.slug,
    name: data.name,
    numUsers: data.users.length,
    status: data.status,
  }

  return {
    listApiKeys: unimplementedInArchive,
    createApiKey: unimplementedInArchive,
    deleteApiKey: unimplementedInArchive,

    readCurrentUser: async () => archiveUser,
    logout: unimplementedInArchive,
    adminSetUserFlags: unimplementedInArchive,
    listSupportedAuths: unimplementedInArchive,
    listAuthDetails: unimplementedInArchive,
    adminDeleteUser: unimplementedInArchive,
    deleteGlobalAuthScheme: unimplementedInArchive,

    listEvidence: async () => data.evidence.map(e => mapEvidence(data, e)),
    createEvidence: unimplementedInArchive,
    readEvidenceContent: async (ids) => await readEvidenceContentFromArchive(data, ids.evidenceUuid),
    readEvidenceImage: (ids) => getEvidencePath(data, ids.evidenceUuid),
    updateEvidence: unimplementedInArchive,
    deleteEvidence: unimplementedInArchive,

    listFindings: async () => data.findings.map(f => mapFinding(data, f)),
    createFinding: unimplementedInArchive,
    readFinding: async (ids) => mapFinding(data, find(data.findings, { uuid: ids.findingUuid })),
    updateFinding: unimplementedInArchive,
    deleteFinding: unimplementedInArchive,
    readFindingEvidence: async (ids) => lookupEvidenceForFinding(data, find(data.findings, { uuid: ids.findingUuid }).id),
    updateFindingEvidence: unimplementedInArchive,

    listOperations: async () => [operation],
    adminListOperations: unimplementedInArchive,
    createOperation: unimplementedInArchive,
    readOperation: async () => operation,
    updateOperation: unimplementedInArchive,
    listUserPermissions: async () => data.users.map(user => ({ user, role: 'read' })),
    updateUserPermissions: unimplementedInArchive,

    listUsers: async (query) => data.users.filter(u => (u.firstName+u.lastName).toLowerCase().indexOf(query) > -1),
    readUser: unimplementedInArchive,
    updateUser: unimplementedInArchive,
    deleteUserAuthScheme: unimplementedInArchive,
    adminListUsers: unimplementedInArchive,
    adminCreateHeadlessUser: unimplementedInArchive,

    listQueries: async () => [], // TODO - queries is null?
    createQuery: unimplementedInArchive,
    updateQuery: unimplementedInArchive,
    deleteQuery: unimplementedInArchive,

    listTags: async () => data.tags,
    createTag: unimplementedInArchive,
    updateTag: unimplementedInArchive,
    deleteTag: unimplementedInArchive,

    // TODO these should go into their respective authschemes:
    createRecoveryCode: unimplementedInArchive,
    deleteExpiredRecoveryCodes: unimplementedInArchive,
    getRecoveryMetrics: unimplementedInArchive,
    adminChangePassword: unimplementedInArchive,
  }
}

function unimplementedInArchive() {
  return Promise.reject(Error('Not available in archive'))
}

// find finds an element by passed predicate, or throws if not found
// example:
// const users = [{id: 1, name: 'alice'}, {id: 2, name: 'bob'}]
// find(users, {id: 2}) -> {id: 2, name: 'bob'}
// This is similar to lodash's `find` except it throws on not found
function find<T extends P, P extends { [k: string]: unknown }>(arr: Array<T>, predicate: P): T {
  const found = arr.find(item => (
    Object.keys(predicate)
      .map(key => item[key] === predicate[key])
      .reduce((p, c) => p && c, true)
  ))
  if (found == null) throw Error(`Could not find ${JSON.stringify(predicate)} in ${JSON.stringify(arr)}`)
  return found
}

function getEvidencePath(data: OperationArchiveData, evidenceUuid: string): string {
  return `media/${find(data.evidence, { uuid: evidenceUuid}).fullImageKey}`
}

async function readEvidenceContentFromArchive(data: OperationArchiveData, evidenceUuid: string): Promise<string> {
  const path = getEvidencePath(data, evidenceUuid)
  const jsonpResult = await fetchJsonp('evidenceJsonp', path)
  if (typeof jsonpResult === 'string') return jsonpResult
  return JSON.stringify(jsonpResult)
}

function mapEvidence(data: OperationArchiveData, evi: OperationArchiveData["evidence"][0]): dtos.Evidence {
  return {
    ...evi,
    operator: find(data.users, { id: evi.operatorId }),
    tags: lookupTagsForEvidence(data, evi.id),
  }
}

function mapFinding(data: OperationArchiveData, finding: OperationArchiveData["findings"][0]): dtos.Finding {
  return {
    ...finding,
    ...calculateFindingRange(data, finding.id),
    tags: lookupTagsForFinding(data, finding.id),
    numEvidence: data.evidence_finding_map.filter(map => map.findingId === finding.id).length,
  }
}

function lookupEvidenceForFinding(data: OperationArchiveData, findingId: number): Array<dtos.Evidence> {
  return data.evidence_finding_map
    .filter(map => map.findingId === findingId)
    .map(map => mapEvidence(data, find(data.evidence, { id: map.evidenceId })))
}

function lookupTagsForEvidence(data: OperationArchiveData, evidenceId: number): Array<dtos.Tag> {
  return data.tag_evidence_map
    .filter(map => map.evidenceId === evidenceId)
    .map(map => find(data.tags, { id: map.tagId }))
}

function lookupTagsForFinding(data: OperationArchiveData, findingId: number): Array<dtos.Tag> {
  return data.evidence_finding_map
    .filter(map => map.findingId === findingId)
    .flatMap(map => lookupTagsForEvidence(data, map.evidenceId))
}

function calculateFindingRange(data: OperationArchiveData, findingId: number): { occurredFrom: string | undefined, occurredTo: string | undefined } {
  const evidenceIds = data.evidence_finding_map
    .filter(map => map.findingId === findingId)
    .map(map => map.evidenceId)
  const evidenceOccurrences = data.evidence
    .filter(evi => evidenceIds.includes(evi.id))
    .map(evi => evi.occurredAt)
  return {
    occurredFrom: minBy(evidenceOccurrences),
    occurredTo: maxBy(evidenceOccurrences),
  }
}
