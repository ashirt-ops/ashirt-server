import { ListObjectForAdminQuery, PaginationResult, UserGroupAdminView } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroupAdminView>> {
  return await ds.adminListUserGroups(i)
}

export async function adminCreateUserGroup(i: {
  name: string,
  userSlugs: string[],
}) {
  return await ds.adminCreateUserGroup(i)
}
