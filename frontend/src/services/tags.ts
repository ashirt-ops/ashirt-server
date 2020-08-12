// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Tag, TagWithUsage, TagByEvidenceDate } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'
import { tagEvidenceDateFromDto } from './data_sources/converters'

export async function createTag(i: {
  operationSlug: string,
  name: string,
  colorName: string,
}): Promise<Tag> {
  return await ds.createTag(
    { operationSlug: i.operationSlug },
    { name: i.name, colorName: i.colorName },
  )
}

export async function getTags(i: {
  operationSlug: string,
}): Promise<Array<TagWithUsage>> {
  return await ds.listTags(i)
}

export async function deleteTag(i: {
  id: number,
  operationSlug: string,
}): Promise<void> {
  await ds.deleteTag({ operationSlug: i.operationSlug, tagId: i.id })
}

export async function updateTag(i: {
  id: number,
  operationSlug: string,
  name: string,
  colorName: string,
}): Promise<void> {
  await ds.updateTag(
    { operationSlug: i.operationSlug, tagId: i.id },
    { name: i.name, colorName: i.colorName },
  )
}

export async function getTagsByEvidenceUsage(i: {
  operationSlug: string,
}): Promise<Array<TagByEvidenceDate>> {
  const data = await ds.listTagsByEvidenceDate(i)
  return data.map(tagEvidenceDateFromDto)
}
