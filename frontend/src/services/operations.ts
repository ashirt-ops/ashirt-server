// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import req from './request_helper'
import { Operation, OperationStatus, UserRole, UserOperationRole, UserFilter } from 'src/global_types'

export async function createOperation(name: string): Promise<Operation> {
  let slug = name.toLowerCase().replace(/[^A-Za-z0-9]+/g, '-').replace(/^-|-$/g, '')
  if (slug === "") {
    return (name === ""
      ? Promise.reject(Error("Operation Name must not be empty"))
      : Promise.reject(Error("Operation Name must include letters or numbers"))
    )
  }
  try {
    return await req('POST', '/operations', {slug, name})
  } catch (err) {
    if (err.message.match(/slug already exists/g)) {
      slug += '-' + Date.now()
      return await req('POST', '/operations', {slug, name})
    }
    throw err
  }
}

export async function getOperations(): Promise<Array<Operation>> {
  return await req('GET', '/operations')
}

export async function getOperationsForAdmin(): Promise<Array<Operation>> {
  return await req('GET', '/admin/operations')
}

export async function getOperation(slug: string): Promise<Operation> {
  return await req('GET', `/operations/${slug}`)
}

export async function saveOperation(slug: string, i: {name: string, status: OperationStatus}) {
  return await req('PUT', `/operations/${slug}`, i)
}

export async function getUserPermissions(i: UserFilter & {
  slug: string,
}): Promise<Array<UserOperationRole>> {
  return await req('GET', `/operations/${i.slug}/users`, null, i)
}

export async function setUserPermission(i: {operationSlug: string, userSlug: string, role: UserRole}) {
  await req('PATCH', `/operations/${i.operationSlug}/users`, {userSlug: i.userSlug, role: i.role})
}
