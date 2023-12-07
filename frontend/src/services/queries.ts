import { SavedQuery } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'
import { queryFromDto } from './data_sources/converters'

export async function getSavedQueries(i: {
  operationSlug: string,
}): Promise<Array<SavedQuery>> {
  const queries = await ds.listQueries(i)
  return queries.map(queryFromDto)
}

export async function saveQuery(i: {
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

export async function upsertQuery(i: {
  operationSlug: string,
  name: string,
  query: string,
  type: 'evidence' | 'findings',
  replaceName?: boolean
}): Promise<void> {
  const {operationSlug, ...rest} = i
  await ds.upsertQuery({operationSlug}, rest)
}

export async function updateSavedQuery(i: {
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

export async function deleteSavedQuery(i: {
  operationSlug: string,
  queryId: number,
}): Promise<void> {
  await ds.deleteQuery(i)
}
