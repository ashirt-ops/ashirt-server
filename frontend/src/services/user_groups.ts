import { ListObjectForAdminQuery, PaginationResult, UserGroup, UserGroupAdminView } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

// TODO TN rename these later?
export async function listUserGroups(i: {
  query: string,
  includeDeleted?: boolean,
}): Promise<Array<UserGroup>> {
  return await ds.listUserGroups(i.query, i.includeDeleted || false)
}

// TODO TN removing group from op doesn't work
// TODO TN add tests for newly added functions
// TODO TN editing group doesn't seem to work

export async function listUserGroupsAdminView(i: ListObjectForAdminQuery): Promise<PaginationResult<UserGroupAdminView>> {
  return await ds.adminListUserGroups(i)
}

export async function createUserGroup(i: {
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

export async function deleteUserGroup(i : { userGroupSlug:string}): Promise<void> {
  return await ds.adminDeleteUserGroup(i)
}

export async function modifyUserGroup(i : {
  slug: string,
  newName: string | null,
  userSlugsToAdd: string[],
  userSlugsToRemove: string[],
}): Promise<void> {
  return await ds.adminModifyUserGroup({ userGroupSlug: i.slug }, i)
}
