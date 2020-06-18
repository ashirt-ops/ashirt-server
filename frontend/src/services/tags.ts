// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import {Tag} from 'src/global_types'

export async function createTag(i: {
  operationSlug: string,
  name: string,
  colorName: string,
}) {
  return await req('POST', `/operations/${i.operationSlug}/tags`, i)
}

export async function getTags(i: {
  operationSlug: string,
}): Promise<Array<Tag>> {
  return await req('GET', `/operations/${i.operationSlug}/tags`)
}

export async function deleteTag(i: {
  id: number,
  operationSlug: string,
}): Promise<void> {
  await req('DELETE', `/operations/${i.operationSlug}/tags/${i.id}`)
}

export async function updateTag(i: {
  id: number,
  operationSlug: string,
  name: string,
  colorName: string,
}): Promise<void> {
  await req('PUT', `/operations/${i.operationSlug}/tags/${i.id}`, {
    name: i.name,
    colorName: i.colorName,
  })
}
