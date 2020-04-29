// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Tag } from 'src/global_types'
import { DataSource } from './data_sources/data_source'

export async function createTag(ds: DataSource, i: {
  operationSlug: string,
  name: string,
  colorName: string,
}): Promise<Tag> {
  return await ds.createTag(
    { operationSlug: i.operationSlug },
    { name: i.name, colorName: i.colorName },
  )
}

export async function getTags(ds: DataSource, i: {
  operationSlug: string,
}): Promise<Array<Tag>> {
  return await ds.listTags(i)
}

export async function deleteTag(ds: DataSource, i: {
  id: number,
  operationSlug: string,
}): Promise<void> {
  await ds.deleteTag({ operationSlug: i.operationSlug, tagId: i.id })
}

export async function updateTag(ds: DataSource, i: {
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
