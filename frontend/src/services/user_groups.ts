import { ListObjectForAdminQuery, PaginationResult, UserGroup, UserGroupAdminView } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

// TODO TN do these naming conventions line up with other examples?
export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroupAdminView>> {
  return await ds.adminListUserGroups(i)
}

export async function adminCreateUserGroup(i: {
  name: string,
  userSlugs: string[],
}): Promise<void> {
  let slug = i.name.toLowerCase().replace(/[^A-Za-z0-9]+/g, '-').replace(/^-|-$/g, '')
  if (slug === "") {
    return (i.name === ""
      ? Promise.reject(Error("User group name must not be empty"))
      : Promise.reject(Error("User group name must include letters or numbers"))
    )
  }
  try {
    return await ds.adminCreateUserGroup({...i, slug})
  } catch (err) {
    if (err.message.match(/slug already exists/g)) {
      slug += '-' + Date.now()
      return await ds.adminCreateUserGroup({...i, slug})
    }
    throw err
  }
}

export async function adminDeleteUserGroup(i : { userGroupSlug:string}): Promise<void> {
  return await ds.adminDeleteUserGroup(i)
}

export async function adminModifyUserGroup(i : {
  slug: string,
  newName: string | null,
  userSlugsToAdd: string[],
  userSlugsToRemove: string[],
}): Promise<void> {
  return await ds.adminModifyUserGroup({ userGroupSlug: i.slug }, i)
}
