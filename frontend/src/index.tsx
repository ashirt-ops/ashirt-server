// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import { backendDataSource } from 'src/services/data_sources/backend'
import { DataSourceContext, getCurrentUser } from 'src/services'
import { render } from 'react-dom'
import { useWiredData } from 'src/helpers'

import Layout from 'src/components/layout'
import Routes from 'src/routes'
import { BrowserRouter } from 'react-router-dom'

require('./base_css')

const RootComponent = () => {
  const wiredUser = useWiredData(React.useCallback(() => getCurrentUser(backendDataSource), []))

  return wiredUser.render(user => (
    <AuthContext.Provider value={{user}}>
      <DataSourceContext.Provider value={backendDataSource}>
        <BrowserRouter>
          <Layout>
            <Routes />
          </Layout>
        </BrowserRouter>
      </DataSourceContext.Provider>
    </AuthContext.Provider>
  ))
}

const root = document.createElement('div')
document.body.appendChild(root)
root.style.height = '100%'
render(<RootComponent /> , root)
