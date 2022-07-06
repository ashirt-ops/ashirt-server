// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

export type Result<T> =
  | { err: Error }
  | { success: T }

export enum OperationStatus {
  // Values here are the backend representations of OperationStatus defined in backend/models
  PLANNING = 0,
  ACTIVE = 1,
  COMPLETE = 2,
}
export const operationStatusToLabel = {
  [OperationStatus.PLANNING]: "Planning",
  [OperationStatus.ACTIVE]: "Active",
  [OperationStatus.COMPLETE]: "Complete",
}

export type EvidenceViewHint = 'small' | 'medium' | 'large'
export type InteractionHint = 'active' | 'inactive'

export enum UserRole {
  // Values here are the backend representations of OperationRole defined in backend/policy
  NO_ACCESS = "",
  READ = "read",
  WRITE = "write",
  ADMIN = "admin",
}
export const userRoleToLabel = {
  [UserRole.NO_ACCESS]: "No Access",
  [UserRole.READ]: "Read",
  [UserRole.WRITE]: "Write",
  [UserRole.ADMIN]: "Admin",
}

export type ApiKey = {
  accessKey: string,
  secretKey: string | null,
  lastAuth: Date | null,
}

export type SupportedEvidenceType =
  | 'image'
  | 'codeblock'
  | 'terminal-recording'
  | 'http-request-cycle'
  | 'event'
  | 'none'

export type Event = {
  type: 'event'
}

export type CodeBlock = {
  type: 'codeblock',
  language: string,
  code: string,
  source: string | null,
}

export type SubmittableCodeblock = {
  type: 'codeblock',
  file: Blob
}

export type TerminalRecording = {
  type: 'terminal-recording',
  file: File
}

export type HttpRequestCycle = {
  type: 'http-request-cycle',
  file: File
}

export type ImageEvidence = {
  type: 'image'
  file: File
}

export type ContentFreeEvidence = {
  type: 'none'
}

export type SubmittableEvidence =
  | SubmittableCodeblock
  | ImageEvidence
  | TerminalRecording
  | HttpRequestCycle
  | ContentFreeEvidence
  | Event

export type Operation = {
  slug: string,
  name: string,
  status: OperationStatus,
  numUsers: number,
}

export type Evidence = {
  uuid: string,
  description: string,
  operator: User,
  occurredAt: Date,
  tags: Array<Tag>,
  contentType: SupportedEvidenceType
}

export type EvidenceMetadata = {
  source: string,
  body: string,
  canProcess?: boolean,
  status?: string, // "Error" | "Queued" | "Completed"
}

export type Finding = {
  uuid: string,
  title: string,
  description: string,
  tags: Array<Tag>,
  numEvidence: number,
  category: string,
  occurredFrom?: Date,
  occurredTo?: Date,
  readyToReport: boolean,
  ticketLink?: string,
}

export type ViewName = 'evidence' | 'findings'
export type SavedQueryType = ViewName

export type SavedQuery = {
  id: number,
  name: string,
  query: string,
  type: SavedQueryType,
}

export type Tag = {
  id: number,
  name: string,
  colorName: string,
}

export type DefaultTag = Tag

export type TagWithUsage = {
  id: number,
  name: string,
  colorName: string,
  evidenceCount: number,
}

export type NewUser = {
  slug: string,
}

export type User = {
  slug: string,
  firstName: string,
  lastName: string,
}

export type UserWithAuth = User & {
  email: string,
  admin: boolean,
}

export type AuthenticationInfo = {
  userKey: string,
  schemeName: string | undefined,
  schemeCode: string,
  schemeType: string,
  lastLogin: Date | null
  authDetails: SupportedAuthenticationScheme | undefined,
}

export type UserOwnView = UserWithAuth & {
  authSchemes: Array<AuthenticationInfo>,
  headless: boolean,
}

export type UserAdminView = UserWithAuth & {
  disabled: boolean,
  headless: boolean,
  deleted: boolean,
  hasLocalTotp: boolean,
  authSchemes: Array<string>,
}

export type UserOperationRole = {
  user: User,
  role: UserRole,
}

export type PaginationQuery = {
  page: number,
  pageSize: number,
}

export type UserFilter = {
  name?: string
}

export type ListUsersForAdminQuery = PaginationQuery & {
  deleted: boolean,
}

export type PaginationResult<T> = PaginationQuery & {
  totalCount: number,
  totalPages: number,
  content: Array<T>,
}

export type SupportedAuthenticationScheme = {
  schemeName: string,
  schemeCode: string,
  schemeType: string,
  schemeFlags: Array<string>
}

export type AuthSchemeDetails = SupportedAuthenticationScheme & {
  userCount: number,
  uniqueUserCount: number,
  labels: Array<string>,
  lastUsed: Date | null,
}

export type RecoveryMetrics = {
  activeCount: number,
  expiredCount: number,
}

export type TagPair = {
  sourceTag: Tag,
  destinationTag: Tag,
}

export type TagDifference = {
  included: Array<TagPair>,
  excluded: Array<Tag>
}

export type TagByEvidenceDate = Tag & {
  usages: Array<Date>,
}

export type FindingCategory = {
  id: number,
  category: string,
  deleted: boolean,
  usageCount: number,
}

export type ServiceWorker = {
  id: number,
  name: string,
  config: string,
  deleted: boolean,
}

export type ActiveServiceWorker = {
  name: string,
}

export type ServiceWorkerTestOutput = {
  id: number,
  name: string,
  live: boolean,
  message: string,
}
