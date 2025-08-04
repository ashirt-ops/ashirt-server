import * as dtos from './dtos/dtos'
import * as types from 'src/global_types'

export function apiKeyFromDto(apiKey: dtos.APIKey): types.ApiKey {
  return {
    ...apiKey,
    lastAuth: apiKey.lastAuth ? new Date(apiKey.lastAuth) : null,
  }
}

export function authenticationInfoFromDto(authInfo: dtos.AuthenticationInfo): types.AuthenticationInfo {
  return {
    ...authInfo,
    lastLogin: authInfo.lastLogin ? new Date(authInfo.lastLogin) : null
  }
}

export function evidenceFromDto(evidence: dtos.Evidence): types.Evidence {
  if (!isValidSupportedEvidenceType(evidence.contentType)) throw Error(`Unknown content type ${evidence.contentType}`)
  return {
    ...evidence,
    adjustedAt: evidence.adjustedAt ? new Date(evidence.adjustedAt) : null,
    occurredAt: new Date(evidence.occurredAt),
    contentType: evidence.contentType,
  }
}

export function findingFromDto(finding: dtos.Finding): types.Finding {
  return {
    ...finding,
    occurredFrom: finding.occurredFrom ? new Date(finding.occurredFrom) : undefined,
    occurredTo: finding.occurredTo ? new Date(finding.occurredTo) : undefined,
  }
}

export function userOperationRoleFromDto({ user, role }: dtos.UserOperationRole): types.UserOperationRole {
  if (!isValidUserRole(role)) throw Error(`Unknown userrole ${role}`)
  return { user, role }
}

export function userGroupOperationRoleFromDto({ userGroup, role }: dtos.UserGroupOperationRole): types.UserGroupOperationRole {
  if (!isValidUserRole(role)) throw Error(`Unknown userrole ${role}`)
  return { userGroup, role }
}

export function userOwnViewFromDto(user: dtos.UserOwnView): types.UserOwnView {
  return { ...user, authSchemes: user.authSchemes.map(authenticationInfoFromDto) }
}

export function queryFromDto(query: dtos.Query): types.SavedQuery {
  if (!isValidQueryType(query.type)) throw Error(`Unknown query type ${query.type}`)
  return { ...query, type: query.type }
}

export function tagEvidenceDateFromDto(tag: dtos.TagByEvidenceDate): types.TagByEvidenceDate {
  return {
    ...tag,
    usages: tag.usages.map((strDate: string) => new Date(strDate))
  }
}

function isValidUserRole(maybeRole: string): maybeRole is types.UserRole {
  // @ts-ignore
  return Object.values(types.UserRole).indexOf(maybeRole) > -1
}

function isValidSupportedEvidenceType(maybeSupportedEvidence: string): maybeSupportedEvidence is types.SupportedEvidenceType {
  return ['codeblock', 'image', 'terminal-recording',
    'http-request-cycle', 'event',
    'none'].indexOf(maybeSupportedEvidence) > -1
}

function isValidQueryType(maybeQueryType: string): maybeQueryType is types.SavedQueryType {
  return ['evidence', 'findings'].indexOf(maybeQueryType) > -1
}
