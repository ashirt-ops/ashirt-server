// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { SavedQuery } from 'src/global_types'
import { DataSource } from './data_sources/data_source'
import { queryFromDto } from './data_sources/converters'

export async function getSavedQueries(ds: DataSource, i: {
  operationSlug: string,
}): Promise<Array<SavedQuery>> {
  const queries = await ds.listQueries(i)
  return queries.map(queryFromDto)
}

export async function saveQuery(ds: DataSource, i: {
  operationSlug: string,
  name: string,
  query: string,
  type: 'evidence'|'findings',
}): Promise<void> {
  await ds.createQuery({ operationSlug: i.operationSlug }, {
    name: i.name.trim(),
    query: i.query.trim(),
    type: i.type,
  })
}

export async function updateSavedQuery(ds: DataSource, i: {
  operationSlug: string,
  queryId: number,
  name: string,
  query: string,
}): Promise<void> {
  await ds.updateQuery({ operationSlug: i.operationSlug, queryId: i.queryId }, {
    name: i.name.trim(),
    query: i.query.trim(),
  })
}

export async function deleteSavedQuery(ds: DataSource, i: {
  operationSlug: string,
  queryId: number,
}): Promise<void> {
  await ds.deleteQuery(i)
}
