// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { DataSource } from './data_sources/data_source'
import { PaginationResult, User, UserAdminView, ListUsersForAdminQuery, UserFilter } from 'src/global_types'

export async function listUsers(ds: DataSource, i: {
  query: string,
  includeDeleted?: boolean,
}): Promise<Array<User>> {
  return await ds.listUsers(i.query, i.includeDeleted || false)
}

export async function listUsersAdminView(ds: DataSource, i: ListUsersForAdminQuery & UserFilter): Promise<PaginationResult<UserAdminView>> {
  return await ds.adminListUsers(i)
}
