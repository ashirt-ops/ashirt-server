// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Tag, TagWithUsage } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

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
