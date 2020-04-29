// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { DataSourceContext, DataSource } from './data_sources/data_source'

export * from './api_keys'
export * from './auth'
export * from './evidence'
export * from './findings'
export * from './operations'
export * from './queries'
export * from './tags'
export * from './users'
export * from './user'
export { DataSourceContext, DataSource }

export function useDataSource(): DataSource {
  return React.useContext(DataSourceContext)
}
