import { backendDataSource as ds } from './data_sources/backend'
import { PaginationResult, User, UserAdminView, ListUsersForAdminQuery, UserFilter } from 'src/global_types'

export async function listUsers(i: {
  query: string,
  includeDeleted?: boolean,
}): Promise<Array<User>> {
  return await ds.listUsers(i.query, i.includeDeleted || false)
}

export async function listUsersAdminView(i: ListUsersForAdminQuery & UserFilter): Promise<PaginationResult<UserAdminView>> {
  return await ds.adminListUsers(i)
}

export async function listEvidenceCreators(i: {
  operationSlug: string
}): Promise<Array<User>> {
  return await ds.listEvidenceCreators({ operationSlug: i.operationSlug })
}
