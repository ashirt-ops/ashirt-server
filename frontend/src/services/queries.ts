// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import {SavedQuery} from 'src/global_types'

export async function getSavedQueries(i: {
  operationSlug: string,
}): Promise<Array<SavedQuery>> {
  return await req('GET', `/operations/${i.operationSlug}/queries`)
}

export async function saveQuery(i: {
  operationSlug: string,
  name: string,
  query: string,
  type: 'evidence'|'findings',
}): Promise<void> {
  await req('POST', `/operations/${i.operationSlug}/queries`, {
    name: i.name.trim(),
    query: i.query.trim(),
    type: i.type,
  })
}

export async function updateSavedQuery(i: {
  operationSlug: string,
  queryId: number,
  name: string,
  query: string,
}): Promise<void> {
  await req('PUT', `/operations/${i.operationSlug}/queries/${i.queryId}`, {
    name: i.name.trim(),
    query: i.query.trim(),
  })
}

export async function deleteSavedQuery(i: {
  operationSlug: string,
  queryId: number,
}): Promise<void> {
  await req('DELETE', `/operations/${i.operationSlug}/queries/${i.queryId}`)
}
