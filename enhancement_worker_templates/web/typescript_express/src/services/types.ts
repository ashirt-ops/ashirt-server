import { SupportedContentTypes } from 'src/helpers/request_validation'

export type CheckConnectionOutput = {
  ok: boolean
}

export type ReadEvidenceOutput = {
  uuid: string
  description: string
  contentType: string
  occurredAt: Date
}

export type EvidenceOutput = {
  uuid: string
  description: string
  occurredAt: Date
  operator: UserOutput
  tags: Array<TagOutputItem>
  contentType: string
}

export type ResponseWrapper<T> = {
  responseCode: number,
  contentType: string
  data: T
}

export type CreateOperationInput = {
  slug: string
  name: string
}

export type OperationOutputItem = {
  slug: string
  name: string
  numUsers: number // int
  status: 1 | 2 | 3
}

export type UserOutput = {
  firstName: string
  lastName: string
  slug: string
}

export type TagOutputItem = {
  id: number // int
  colorName: string
  name: string
}

export type TagWithUsageOutputItem = TagOutputItem & {
  evidenceCount: number // int
}
export type ListOperationTagsOutput = Array<TagWithUsageOutputItem>
export type ListOperationsOutput = Array<OperationOutputItem>

export type UpsertMetadataInput = {
  source: string
  body: string
  status: string
  message?: string
  canProcess?: boolean
}

export type CreateEvidenceInput = {
  notes: string
  file: FileData
  contentType?: typeof SupportedContentTypes[number]
  occurred_at?: string
  tagIds: Array<number> // int array
}

export type UpdateEvidenceInput = {
  notes: string
  file?: FileData
  contentType?: typeof SupportedContentTypes[number]
  occurred_at?: string
  tagsToAdd?: Array<number> // int array
  tagsToRemove?: Array<number> // int array
}

export type CreateTagInput = {
  name: string
  colorName: string
}

export type FileData = {
  filename: string
  mimetype: string
  content: Buffer
}
