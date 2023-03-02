// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Tag, TagWithUsage, DefaultTag } from 'src/global_types'
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

export async function getDefaultTags(): Promise<Array<DefaultTag>> {
  return await ds.listDefaultTags()
}

export async function createDefaultTag(i: {
  name: string,
  colorName: string,
}): Promise<DefaultTag> {
  return await ds.createDefaultTag(
    { name: i.name, colorName: i.colorName },
  )
}

export async function deleteDefaultTag(i: {
  id: number,
}): Promise<void> {
  await ds.deleteDefaultTag({ tagId: i.id })
}

export async function updateDefaultTag(i: {
  id: number,
  name: string,
  colorName: string,
}): Promise<void> {
  await ds.updateDefaultTag(
    { tagId: i.id },
    { name: i.name, colorName: i.colorName },
  )
}

export async function mergeDefaultTags(i: Array<{
  name: string,
  colorName: string,
}>) {
  await ds.mergeDefaultTags(i)
}
