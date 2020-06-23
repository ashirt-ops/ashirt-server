// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { DataSource } from '../data_source'
import { default as req, xhrText as reqText, reqMultipart } from './request_helper'

export const backendDataSource: DataSource = {
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
  readEvidenceContent: ids => reqText('GET', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}/media`),
  updateEvidence: (ids, formData) => reqMultipart('PUT', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, formData),
  deleteEvidence: (ids, payload) => req('DELETE', `/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, payload),
  getEvidenceMigrationDifference: (ids, fromOperationSlug) => req('GET', `/move/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, null, { sourceOperationSlug: fromOperationSlug }),
  moveEvidence: (ids, fromOperationSlug) => req('PUT', `/move/operations/${ids.operationSlug}/evidence/${ids.evidenceUuid}`, null, { sourceOperationSlug: fromOperationSlug }),

  listFindings: (ids, query) => req('GET', `/operations/${ids.operationSlug}/findings`, null, { query }),
  createFinding: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/findings`, payload),
  readFinding: ids => req('GET', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`),
  updateFinding: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`, payload),
  deleteFinding: ids => req('DELETE', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}`),
  readFindingEvidence: ids => req('GET', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}/evidence`),
  updateFindingEvidence: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/findings/${ids.findingUuid}/evidence`, payload),

  listOperations: () => req('GET', '/operations'),
  adminListOperations: () => req('GET', '/operations'),
  createOperation: payload => req('POST', '/operations', payload),
  readOperation: ids => req('GET', `/operations/${ids.operationSlug}`),
  updateOperation: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}`, payload),
  listUserPermissions: (ids, query) => req('GET', `/operations/${ids.operationSlug}/users`, null, query),
  updateUserPermissions: (ids, payload) => req('PATCH', `/operations/${ids.operationSlug}/users`, payload),

  listUsers: (query, includeDeleted) => req('GET', '/users', null, { query, includeDeleted }),
  readUser: ids => req('GET', `/user`, null, ids),
  updateUser: (ids, payload) => req('POST', `/user/profile/${ids.userSlug}`, payload),
  deleteUserAuthScheme: ids => req('DELETE', `/user/${ids.userSlug}/scheme/${ids.authSchemeName}`),
  adminListUsers: query => req('GET', '/admin/users', null, query),
  adminCreateHeadlessUser: payload => req('POST', "/admin/user/headless", payload),

  listQueries: ids => req('GET', `/operations/${ids.operationSlug}/queries`),
  createQuery: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/queries`, payload),
  updateQuery: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/queries/${ids.queryId}`, payload),
  deleteQuery: ids => req('DELETE', `/operations/${ids.operationSlug}/queries/${ids.queryId}`),

  listTags: ids => req('GET', `/operations/${ids.operationSlug}/tags`),
  createTag: (ids, payload) => req('POST', `/operations/${ids.operationSlug}/tags`, payload),
  updateTag: (ids, payload) => req('PUT', `/operations/${ids.operationSlug}/tags/${ids.tagId}`, payload),
  deleteTag: ids => req('DELETE', `/operations/${ids.operationSlug}/tags/${ids.tagId}`),

  // TODO these should go into their respective authschemes:
  createRecoveryCode: ids => req('POST', '/auth/recovery/generate', ids),
  deleteExpiredRecoveryCodes: () => req('DELETE', '/auth/recovery/expired'),
  getRecoveryMetrics: () => req('GET', '/auth/recovery/metrics'),
  adminChangePassword: i => req('PUT', '/auth/local/admin/password', i),
}
