import { ListObjectForAdminQuery, PaginationResult, UserGroupAdminView } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

// TODO TN do these naming conventions line up with other examples?
export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroupAdminView>> {
  return await ds.adminListUserGroups(i)
}

export async function adminCreateUserGroup(i: {
  name: string,
  userSlugs: string[],
}): Promise<void> {
  return await ds.adminCreateUserGroup(i)
}

export async function adminDeleteUserGroup(i : { userGroupSlug:string}): Promise<void> {
  return await ds.adminDeleteUserGroup(i)
}
