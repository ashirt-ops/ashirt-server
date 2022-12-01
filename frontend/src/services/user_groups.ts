import { ListObjectForAdminQuery, PaginationResult } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'
import { UserGroup } from './data_sources/dtos/dtos'

export async function createUserGroup(name: string) {
  return await ds.createUserGroup(name)
}

export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroup>> {
  return await ds.adminListUserGroups(i)
}
