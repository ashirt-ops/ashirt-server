// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import { PaginationResult, User, UserAdminView, ListUsersForAdminQuery, UserFilter } from 'src/global_types'

export async function listUsers(i: {
  query: string,
  includeDeleted?: boolean,
}): Promise<Array<User>> {
  return await req('GET', '/users', null, i)
}

export async function listUsersAdminView(i: ListUsersForAdminQuery & UserFilter): Promise<PaginationResult<UserAdminView>> {
  return await req('GET', '/admin/users', null, i)
}
