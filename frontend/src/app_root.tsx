// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

require('./base_css')
import * as React from 'react'
import AuthContext from 'src/auth_context'
import { DataSourceContext, DataSource } from 'src/services'
import { getCurrentUser } from 'src/services'
import { render } from 'react-dom'
import { useWiredData } from 'src/helpers'

import Layout from 'src/components/layout'
import Routes from 'src/routes'

export function bootApp(
  RouterProvider: React.ComponentType,
  ds: DataSource,
) {
  const AppRoot = () => {
    const wiredUser = useWiredData(React.useCallback(() => (
      getCurrentUser(ds)
    ), []))

    return wiredUser.render(user => (
      <AuthContext.Provider value={{user}}>
        <DataSourceContext.Provider value={ds}>
          <RouterProvider>
            <Layout>
              <Routes />
            </Layout>
          </RouterProvider>
        </DataSourceContext.Provider>
      </AuthContext.Provider>
    ))
  }

  const root = document.createElement('div')
  document.body.appendChild(root)
  root.style.height = '100%'
  render(<AppRoot />, root)
}
