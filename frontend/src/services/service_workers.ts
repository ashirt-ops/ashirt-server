// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { ServiceWorker } from 'src/global_types'
import { backendDataSource as ds } from './data_sources/backend'

export async function createServiceWorker(i: {
  name: string,
  config: string,
}): Promise<void> {
  return ds.adminCreateServiceWorker(i)
}

export async function listServiceWorkers(): Promise<Array<ServiceWorker>> {
  return ds.adminListServiceWorkers()
}

export async function updateServiceWorker(i: {
  id: number,
  name: string,
  config: string,
}): Promise<void> {
  const { id, ...payload } = i
  return ds.adminUpdateServiceWorker({ serviceWorkerId: id }, payload)
}

export async function deleteServiceWorkers(i: {
  id: number
}): Promise<void> {
  return ds.adminDeleteServiceWorker({ serviceWorkerId: i.id })
}
