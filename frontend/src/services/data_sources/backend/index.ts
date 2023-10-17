// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { DataSource, cacheBust } from '../data_source'
import { default as req, xhrText as reqText, reqMultipart } from './request_helper'
cacheBust()

export const backendDataSource: DataSource = {
  flags: () => req('GET', '/flags'),

  listApiKeys: ids => req('GET', '/user/apikeys', null, ids),
  createApiKey: ids => req('POST', `/user/${ids.userSlug}/apikeys`),
  deleteApiKey: ids => req('DELETE', `/user/${ids.userSlug}/apikeys/${ids.accessKey}`),

  readCurrentUser: () => req('GET', '/user'),
  logout: () => req('POST', '/logout'),
  adminSetUserFlags: (ids, flags) => req('POST', `/admin/${ids.userSlug}/flags`, flags),
  listSupportedAuths: () => req('GET', '/auths'),
  listAuthDetails: () => req('GET', '/auths/breakdown'),
  adminDeleteUser: ids => req('DELETE', `/admin/user/${ids.userSlug}`),
  deleteGlobalAuthScheme: ids => req('DELETE', `/auths/${ids.schemeName}`),

  listEvidence: (ids, query) => req('GET', `/operations/${ids.operationSlug}/evidence`, null, { query }),
  createEvidence: (ids, formData) => reqMultipart('POST', `/operations/${ids.operationSlug}/evidence`, formData),
  readEvidenceContent: ids => req('GET', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/media`),
  readEvidenceContentCodeblock: ids => req('GET', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/codeblock`),
  updateEvidence: (ids, formData) => reqMultipart('PUT', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, formData),
  deleteEvidence: (ids, payload) => req('DELETE', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, payload),
  getEvidenceMigrationDifference: (ids, fromOperationSlug) => req('GET', `/move/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, null, { sourceOperationSlug: fromOperationSlug }),
  moveEvidence: (ids, fromOperationSlug) => req('PUT', `/move/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, { sourceOperationSlug: fromOperationSlug }),
  createEvidenceMetadata: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/metadata`, payload),
  updateEvidenceMetadata: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/metadata`, payload),
  readEvidenceMetadata: (ids) => req('GET', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/metadata`),
  runServiceWorkerForEvidence: (ids) => req('PUT', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/metadata/${ids.source}/run`),
  runServiceWorkerBatch: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/metadata/run`, payload),

  runAllServiceWorkersForEvidence: (ids) => req('PUT', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/metadata/run`),

  listFindingCategories: (includeDeleted) => req('GET', `/findings/categories`, null, { includeDeleted }),
  createFindingCategory: (payload) => req('POST', `/findings/category`, payload),
  deleteFindingCategory: (ids, payload) => req('DELETE', `/findings/category/${ids.findingCategoryId}`, payload),
  updateFindingCategory: (ids, payload) => req('PUT', `/findings/category/${ids.findingCategoryId}`, payload),
  listFindings: (ids, query) => req('GET', `/operations/${ids.operationSlug}/findings`, null, { query }),
  createFinding: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/findings`, payload),
  readFinding: ids => req('GET', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`),
  updateFinding: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`, payload),
  deleteFinding: ids => req('DELETE', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`),
  readFindingEvidence: ids => req('GET', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}/evidence`),
  updateFindingEvidence: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}/evidence`, payload),

  listOperations: () => req('GET', '/operations'),
  adminListOperations: () => req('GET', '/admin/operations'),
  createOperation: payload => req('POST', '/operations', payload),
  readOperation: ids => req('GET', `/operations/${ids.operationSlug}`),
  updateOperation: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}`, payload),
  listUserPermissions: (ids, query) => req('GET', `/operations/${ids.operationSlug}/users`, null, query),
  listUserGroupPermissions: (ids, query) => req('GET', `/operations/${ids.operationSlug}/usergroups`, null, query),
  updateUserPermissions: (ids, payload) => req('PATCH', `/operations/${ids.operationSlug}/users`, payload),
  updateUserGroupPermissions: (ids, payload) => req('PATCH', `/operations/${ids.operationSlug}/usergroups`, payload),
  deleteOperation: (ids) => req('DELETE', `/operations/${ids.operationSlug}`),
  setFavorite: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/favorite`, payload),

  listUsers: (query, includeDeleted) => req('GET', '/users', null, { query, includeDeleted }),
  readUser: ids => req('GET', `/user`, null, ids),
  listEvidenceCreators: (ids) => req('GET', `/operations/${ids.operationSlug}/evidence/creators`),
  updateUser: (ids, payload) => req('POST', `/user/profile/${ids.userSlug}`, payload),
  deleteUserAuthScheme: ids => req('DELETE', `/user/${ids.userSlug}/scheme/${ids.authSchemeName}`),
  adminListUsers: query => req('GET', '/admin/users', null, query),
  adminCreateHeadlessUser: payload => req('POST', "/admin/user/headless", payload),

  listUserGroups: (query, includeDeleted, operationSlug) => req('GET', '/usergroups', null, { query, includeDeleted, operationSlug }),
  adminCreateUserGroup: payload => req('POST', '/admin/usergroups', payload),
  adminListUserGroups: query => req('GET', '/admin/usergroups', null, query),
  adminDeleteUserGroup: ids => req('DELETE', `/admin/usergroups/${ids.userGroupSlug}`),
  adminModifyUserGroup: (ids, payload) => req('PUT', `/admin/usergroups/${ids.userGroupSlug}`, payload),

  listQueries: ids => req('GET', `/operations/${ids.operationSlug}/queries`),
  createQuery: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/queries`, payload),
  upsertQuery: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/queries`, payload),
  updateQuery: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/queries/${ids.queryId}`, payload),
  deleteQuery: ids => req('DELETE', `/operations/${ids.operationSlug}/queries/${ids.queryId}`),

  listTags: ids => req('GET', `/operations/${ids.operationSlug}/tags`),
  createTag: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/tags`, payload),
  updateTag: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/tags/${ids.tagId}`, payload),
  deleteTag: ids => req('DELETE', `/operations/${ids.operationSlug}/tags/${ids.tagId}`),

  listDefaultTags: () => req('GET', `/admin/tags`),
  createDefaultTag: (payload) => req('POST', `/admin/tags`, payload),
  updateDefaultTag: (ids, payload) => req('PUT', `/admin/tags/${ids.tagId}`, payload),
  deleteDefaultTag: (ids) => req('DELETE', `/admin/tags/${ids.tagId}`),
  mergeDefaultTags: (payload) => req('POST', `/admin/merge/tags`, payload),

  adminListServiceWorkers: () => req('GET', `/admin/services`),
  adminCreateServiceWorker: (payload) => req('POST', `/admin/services`, payload),
  adminUpdateServiceWorker: (ids, payload) => req('PUT', `/admin/services/${ids.serviceWorkerId}`, payload),
  adminDeleteServiceWorker: (ids) => req('DELETE', `/admin/services/${ids.serviceWorkerId}`, { delete: true }),
  adminUnDeleteServiceWorker: (ids) => req('DELETE', `/admin/services/${ids.serviceWorkerId}`, { delete: false }),
  adminTestServiceWorker: (ids) => req('GET', `/admin/services/${ids.serviceWorkerId}/test`),
  listActiveServiceWorkers: () => req('GET', `/services`),

  // TODO these should go into their respective authschemes:
  createRecoveryCode: ids => req('POST', '/auth/recovery/generate', ids),
  deleteExpiredRecoveryCodes: () => req('DELETE', '/auth/recovery/expired'),
  getRecoveryMetrics: () => req('GET', '/auth/recovery/metrics'),
  adminChangePassword: i => req('PUT', '/auth/local/admin/password', i),
  adminCreateLocalUser: i => req('POST', '/auth/local/admin/register', i),
  adminInviteUser: i => req('POST', '/auth/recovery/admin/register', i),
  getTotpForUser: ids => req('GET', '/auth/local/totp', ids),
  deleteTotpForUser: ids => req('DELETE', '/auth/local/totp', ids),

  listGlobalVars: () => req('GET', '/global-vars'),
  createGlobalVar: payload => req('POST', '/global-vars', payload),
  updateGlobalVar: (ids, payload) => req('PUT', `/global-vars/${ids.name}`, payload),
  deleteGlobalVar: (ids) => req('DELETE', `/global-vars/${ids.name}`),
}
