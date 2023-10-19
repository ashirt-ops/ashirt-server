// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

export type SuccessfulResult<T> = { success: T }
export type ErrorResult = { err: Error }

export type Result<T> =
  | ErrorResult
  | SuccessfulResult<T>

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

export type TopContrib = {
  slug: string,
  count: number,
}

export type EvidenceCount = {
  imageCount: number,
  codeblockCount: number,
  recordingCount: number,
  eventCount: number,
  harCount: number,
}

export type Operation = {
  slug: string,
  name: string,
  numUsers: number,
  numEvidence: number,
  numTags: number,
  favorite: boolean,
  topContribs: Array<TopContrib>,
  evidenceCount: EvidenceCount,
  userCanViewGroups?: boolean,
  userCanExportData?: boolean,
}

export type Evidence = {
  uuid: string,
  description: string,
  operator: User,
  occurredAt: Date,
  tags: Array<Tag>,
  contentType: SupportedEvidenceType
  sendImageInfo: boolean,
}

export type ExportedEvidence = Omit<Evidence, 'tags' | 'uuid'> & {
  filename?: string,
  sourceFilename?: string,
  tags: Array<string | Tag>,
  uuid?: string,
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

export type DenormalizedTag = {
  name: string;
} & Partial<Tag>;

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
  username: string,
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

export type UserGroup = {
  slug: string,
  name: string,
  userSlugs: Array<string>,
}

export type UserGroupAdminView = UserGroup & IncludeDeleted 

export type UserOperationRole = {
  user: User,
  role: UserRole,
}

export type UserGroupOperationRole = {
  userGroup: UserGroupAdminView,
  role: UserRole,
}

export type PaginationQuery = {
  page: number,
  pageSize: number,
}

export type UserFilter = {
  name?: string
}

export type IncludeDeleted = {
  deleted: boolean,
}

export type ListUsersForAdminQuery = PaginationQuery & IncludeDeleted

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

export type FilterText = {
  value: string;
  onChange: React.Dispatch<React.SetStateAction<string>>;
  disabled: boolean;
  setDisabled: React.Dispatch<React.SetStateAction<boolean>>;
}

export type ContentType = "image" | "terminal-recording" | "http-request-cycle" | "event" | "none" | "codeblock"

type Languages = "" | "abap" | "actionscript" | "ada" | "c_cpp" | "csharp" | "cobol" | "d" | "dart" | "dockerfile" | "elixir" | "elm" | "erlang" | "fsharp" | "fortran" | "golang" | "groovy" | "haskell" | "java" | "javascript" | "julia" | "kotlin" | "lisb" | "lua" | "matlab" | "markdown" | "objectivec" | "pascal" | "php" | "perl" | "prolog" | "properties" | "python" | "r" | "ruby" | "rust" | "sass" | "scala" | "scheme" | "sh" | "sql" | "swift" | "tcl" | "terraform" | "toml" | "typescript" | "vbscript" | "xml"

export interface Media {
  filename: string,
  contentType: ContentType,
  contentSubtype?: Languages,
  sourceFilename?: string,
  blob: Blob
}

export interface Codeblock {
  contentType: string,
  contentSubtype: Languages,
  content: string 
  metadata: {
    source: string,
  }
}

export type GlobalVar = {
  name: string,
  value: string,
}

export type ImageInfo = {
  url: string
}
