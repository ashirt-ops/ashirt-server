import { ListObjectForAdminQuery, PaginationResult, UserGroupAdminView } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createUserGroup(name: string) {
  return await ds.createUserGroup({name: name.trim()})
}

export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroupAdminView>> {
  return await ds.adminListUserGroups(i)
}
