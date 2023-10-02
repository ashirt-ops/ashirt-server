// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dtos from './dtos/dtos'
import * as types from 'src/global_types'

type EvidenceUuid = { evidenceUuid: string }
type FindingUuid = { findingUuid: string }
type OpSlug = { operationSlug: string }
type UserSlug = { userSlug: string }
type UserGroupSlug = { userGroupSlug: string }
type QueryId = { queryId: number }
type TagId = { tagId: number }
type FindingCategoryId = { findingCategoryId: number }
type ServiceWorkerId = { serviceWorkerId: number }
type Name = { name: string }
type OpAndVarSlugs = { operationSlug: string, varSlug: string }

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

type TagPayload = {
  name: string,
  colorName: string
}

type ServiceWorkerPayload = {
  name: string
  config: string
}

export interface DataSource {
  flags(): Promise<dtos.Flags>

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
  moveEvidence(ids: OpSlug & EvidenceUuid, fromOperationSlug: string): Promise<void>
  createEvidenceMetadata(ids: OpSlug & EvidenceUuid, payload: { source: string, body: string }): Promise<void>
  updateEvidenceMetadata(ids: OpSlug & EvidenceUuid, payload: { source: string, body: string }): Promise<void>
  readEvidenceMetadata(ids: OpSlug & EvidenceUuid): Promise<Array<dtos.EvidenceMetadata>>
  runServiceWorkerForEvidence(ids: OpSlug & EvidenceUuid & { source: string }): Promise<void>
  runServiceWorkerBatch(ids: OpSlug, payload: { workers: Array<string>, evidenceUuids: Array<string> } ): Promise<void>
  runAllServiceWorkersForEvidence(ids: OpSlug & EvidenceUuid): Promise<void>

  listFindingCategories(includeDeleted: boolean): Promise<Array<dtos.FindingCategory>>
  createFindingCategory(payload: { category: string }): Promise<dtos.FindingCategory>
  deleteFindingCategory(ids: FindingCategoryId, payload: { delete: boolean }): Promise<void>
  updateFindingCategory(ids: FindingCategoryId, payload: { category: string }): Promise<void>
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
  updateOperation(ids: OpSlug, payload: { name: string }): Promise<void>
  listUserPermissions(ids: OpSlug, query: { name?: string }): Promise<Array<dtos.UserOperationRole>>
  listUserGroupPermissions(ids: OpSlug, query: { name?: string }): Promise<Array<dtos.UserGroupOperationRole>>
  updateUserPermissions(ids: OpSlug, payload: { userSlug: string, role: types.UserRole }): Promise<void>
  updateUserGroupPermissions(ids: OpSlug, payload: { userGroupSlug: string, role: types.UserRole }): Promise<void>
  deleteOperation(ids: OpSlug): Promise<void>
  setFavorite(ids: OpSlug, payload: { favorite: boolean }): Promise<void>

  listUsers(query: string, includeDeleted: boolean): Promise<Array<dtos.User>>
  readUser(ids: UserSlug): Promise<dtos.UserOwnView>
  listEvidenceCreators(ids: OpSlug): Promise<Array<dtos.User>>,
  updateUser(ids: UserSlug, payload: UserPayload): Promise<void>
  deleteUserAuthScheme(ids: UserSlug & { authSchemeName: string }): Promise<void>
  adminListUsers(query: { deleted: boolean, name?: string }): Promise<types.PaginationResult<dtos.UserAdminView>>
  adminCreateHeadlessUser(payload: UserPayload): Promise<dtos.CreateUserOutput>

  listUserGroups(query: string, includeDeleted: boolean, operationSlug: string): Promise<Array<dtos.UserGroupAdminView>>
  adminListUserGroups(query: { deleted: boolean }): Promise<dtos.UserGroupAdminView[]>
  adminCreateUserGroup(payload: { slug: string, name: string, userSlugs: string[] }): Promise<void>
  adminDeleteUserGroup(ids: UserGroupSlug): Promise<void>
  adminModifyUserGroup(ids: UserGroupSlug, payload: { newName: string | null, userSlugsToAdd: string[], userSlugsToRemove: string[], }): Promise<void>

  listQueries(ids: OpSlug): Promise<Array<dtos.Query>>
  createQuery(ids: OpSlug, payload: { name: string, query: string, type: 'evidence' | 'findings' }): Promise<void>
  upsertQuery(ids: OpSlug, payload: { name: string, query: string, type: 'evidence' | 'findings', replaceName?: boolean }): Promise<void>
  updateQuery(ids: OpSlug & QueryId, payload: { name: string, query: string }): Promise<void>
  deleteQuery(ids: OpSlug & QueryId): Promise<void>

  listTags(ids: OpSlug): Promise<Array<dtos.TagWithUsage>>
  createTag(ids: OpSlug, payload: TagPayload): Promise<dtos.Tag>
  updateTag(ids: OpSlug & TagId, payload: TagPayload): Promise<void>
  deleteTag(ids: OpSlug & TagId): Promise<void>

  listDefaultTags(): Promise<Array<dtos.DefaultTag>>
  createDefaultTag(payload: TagPayload): Promise<dtos.DefaultTag>
  updateDefaultTag(ids: TagId, payload: TagPayload): Promise<void>
  deleteDefaultTag(ids: TagId): Promise<void>
  mergeDefaultTags(payload: Array<TagPayload>): Promise<void>

  adminListServiceWorkers(): Promise<Array<dtos.ServiceWorker>>
  adminCreateServiceWorker(payload: ServiceWorkerPayload): Promise<void>
  adminUpdateServiceWorker(ids: ServiceWorkerId, payload: ServiceWorkerPayload): Promise<void>
  adminDeleteServiceWorker(ids: ServiceWorkerId): Promise<void>
  adminUnDeleteServiceWorker(ids: ServiceWorkerId): Promise<void>
  adminTestServiceWorker(ids: ServiceWorkerId): Promise<dtos.ServiceWorkerTestOutput>
  listActiveServiceWorkers(): Promise<Array<dtos.ActiveServiceWorker>>

  // TODO these should go into their respective authschemes:
  createRecoveryCode(ids: UserSlug): Promise<{ code: string }>
  deleteExpiredRecoveryCodes(): Promise<void>
  getRecoveryMetrics(): Promise<any>
  adminChangePassword(i: { userSlug: string, newPassword: string }): Promise<void>
  adminCreateLocalUser(i: { firstName: string, lastName?: string, email: string }): Promise<dtos.NewUserCreatedByAdmin>,
  adminInviteUser(i: { firstName: string, lastName?: string, email: string }): Promise<{ code: string }>,
  getTotpForUser(ids: UserSlug): Promise<boolean>
  deleteTotpForUser(ids: UserSlug): Promise<void>

  listGlobalVars(): Promise<Array<dtos.GlobalVar>>
  createGlobalVar(payload: { name: string, value: string | null }): Promise<dtos.GlobalVar>
  updateGlobalVar(ids: Name, payload: { value: string | null, newName: string | null }): Promise<void>
  deleteGlobalVar(ids: Name): Promise<void>

  listOperationVars(ids: OpSlug): Promise<Array<dtos.OperationVar>>
  createOperationVar(ids: OpSlug, payload: { varSlug: string, name: string, value: string | null }): Promise<dtos.OperationVar>
  updateOperationVar(ids: OpAndVarSlugs, payload: {  value: string | null, name: string | null }): Promise<void>
  deleteOperationVar(ids: OpAndVarSlugs): Promise<void>
}

// Since both dtos & this file only contains typescript types, webpack doesn't pick up the
// changes unless there is some actual executable javascript reverenced from
// the app itself. By exporting an empty function and calling it in the app
// https://github.com/TypeStrong/ts-loader/issues/808
dtos.cacheBust()
export function cacheBust() { }
