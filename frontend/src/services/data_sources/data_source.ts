// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dtos from './dtos'
import * as types from 'src/global_types'

type EvidenceUuid = { evidenceUuid: string }
type FindingUuid = { findingUuid: string }
type OpSlug = { operationSlug: string }
type UserSlug = { userSlug: string }
type QueryId = { queryId: number }
type TagId = { tagId: number }

type FindingPayload = {
  category: string,
  title: string,
  description: string,
}

type UserPayload = {
  firstName: string,
  lastName: string,
  email: string,
}

export interface DataSource {
  listApiKeys(ids?: UserSlug): Promise<Array<dtos.APIKey>>
  createApiKey(ids: UserSlug): Promise<dtos.APIKey>
  deleteApiKey(ids: UserSlug & { accessKey: string }): Promise<void>

  readCurrentUser(): Promise<dtos.UserOwnView>
  logout(): Promise<void>
  adminSetUserFlags(ids: UserSlug, flags: { disabled: boolean, admin: boolean }): Promise<void>
  listSupportedAuths(): Promise<Array<dtos.SupportedAuthScheme>>
  listAuthDetails(): Promise<Array<dtos.DetailedAuthenticationInfo>>
  adminDeleteUser(ids: UserSlug): Promise<void>
  deleteGlobalAuthScheme(ids: { schemeName: string }): Promise<void>

  listEvidence(ids: OpSlug, query: string): Promise<Array<dtos.Evidence>>
  createEvidence(ids: OpSlug, formData: FormData): Promise<void>
  readEvidenceContent(ids: OpSlug & EvidenceUuid): Promise<string>
  updateEvidence(ids: OpSlug & EvidenceUuid, formData: FormData): Promise<void>
  deleteEvidence(ids: OpSlug & EvidenceUuid, payload: { deleteAssociatedFindings: boolean }): Promise<void>
  getEvidenceMigrationDifference(ids: OpSlug & EvidenceUuid, fromOperationSlug: string): Promise<dtos.TagDifference>

  listFindings(ids: OpSlug, query: string): Promise<Array<dtos.Finding>>
  createFinding(ids: OpSlug, payload: FindingPayload): Promise<dtos.Finding>
  readFinding(ids: OpSlug & FindingUuid): Promise<dtos.Finding>
  updateFinding(ids: OpSlug & FindingUuid, payload: FindingPayload & { readyToReport: boolean, ticketLink: string | null }): Promise<void>
  deleteFinding(ids: OpSlug & FindingUuid): Promise<void>
  readFindingEvidence(ids: OpSlug & FindingUuid): Promise<Array<dtos.Evidence>>
  updateFindingEvidence(ids: OpSlug & FindingUuid, payload: { evidenceToAdd: Array<string>, evidenceToRemove: Array<string> }): Promise<void>

  listOperations(): Promise<Array<dtos.Operation>>
  adminListOperations(): Promise<Array<dtos.Operation>>
  createOperation(payload: { slug: string, name: string }): Promise<dtos.Operation>
  readOperation(ids: OpSlug): Promise<dtos.Operation>
  updateOperation(ids: OpSlug, payload: { name: string, status: types.OperationStatus }): Promise<void>
  listUserPermissions(ids: OpSlug, query: { name?: string }): Promise<Array<dtos.UserOperationRole>>
  updateUserPermissions(ids: OpSlug, payload: { userSlug: string, role: types.UserRole }): Promise<void>

  listUsers(query: string, includeDeleted: boolean): Promise<Array<dtos.User>>
  readUser(ids: UserSlug): Promise<dtos.UserOwnView>
  updateUser(ids: UserSlug, payload: UserPayload): Promise<void>
  deleteUserAuthScheme(ids: UserSlug & { authSchemeName: string }): Promise<void>
  adminListUsers(query: { deleted: boolean, name?: string }): Promise<types.PaginationResult<dtos.UserAdminView>>
  adminCreateHeadlessUser(payload: UserPayload): Promise<void>

  listQueries(ids: OpSlug): Promise<Array<dtos.Query>>
  createQuery(ids: OpSlug, payload: { name: string, query: string, type: 'evidence' | 'findings' }): Promise<void>
  updateQuery(ids: OpSlug & QueryId, payload: { name: string, query: string }): Promise<void>
  deleteQuery(ids: OpSlug & QueryId): Promise<void>

  listTags(ids: OpSlug): Promise<Array<dtos.Tag>>
  createTag(ids: OpSlug, payload: { name: string, colorName: string }): Promise<dtos.Tag>
  updateTag(ids: OpSlug & TagId, payload: { name: string, colorName: string }): Promise<void>
  deleteTag(ids: OpSlug & TagId): Promise<void>

  // TODO these should go into their respective authschemes:
  createRecoveryCode(ids: UserSlug): Promise<{ code: string }>
  deleteExpiredRecoveryCodes(): Promise<void>
  getRecoveryMetrics(): Promise<any>
  adminChangePassword(i: { userSlug: string, newPassword: string }): Promise<void>
}
