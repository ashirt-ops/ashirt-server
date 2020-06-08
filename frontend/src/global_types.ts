// Copyright 2020, Verizon Media
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

export enum ExportStatus {
  PENDING = 0,
  IN_PROGRESS = 1,
  COMPLETE = 2,
  ERROR = 3,
  CANCELLED = 4,
}

export const exportStatusToLabel = {
  [ExportStatus.PENDING]: "Pending",
  [ExportStatus.IN_PROGRESS]: "In Progress",
  [ExportStatus.COMPLETE]: "Complete",
  [ExportStatus.ERROR]: "Error",
  [ExportStatus.CANCELLED]: "Cancelled",
}


export type ApiKey = {
  accessKey: string,
  secretKey: string | null,
  lastAuth: Date | null,
}

export type SupportedEvidenceType = 'codeblock' | 'terminal-recording' | 'image' | 'none'

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
  | ContentFreeEvidence

export type Operation = {
  slug: string,
  name: string,
  status: OperationStatus,
  numUsers: number,
}

export type OperationWithExportData = Operation & {
  lastCompletedExport: Date | null,
  exportStatus: ExportStatus | null,
}

export type Evidence = {
  uuid: string,
  description: string,
  operator: User,
  occurredAt: Date,
  tags: Array<Tag>,
  contentType: SupportedEvidenceType
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

export type SavedQueryType = 'evidence' | 'findings'

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
  schemeCode: string,
  lastLogin: Date | null
}

export type UserOwnView = UserWithAuth & {
  authSchemes: Array<AuthenticationInfo>,
  headless: boolean,
}

export type UserAdminView = UserWithAuth & {
  disabled: boolean,
  headless: boolean,
  deleted: boolean,
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
