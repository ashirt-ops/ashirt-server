// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

require('./base_css')
import * as React from 'react'
import AuthContext from 'src/auth_context'
import Layout from 'src/components/layout'
import Routes from 'src/routes'
import { BrowserRouter } from 'react-router-dom'
import { getCurrentUser } from 'src/services'
import { createRoot } from 'react-dom/client'
import { useWiredData } from 'src/helpers'

const RootComponent = () => {
  const wiredUser = useWiredData(getCurrentUser)

  return wiredUser.render(user => (
    <AuthContext.Provider value={{ user }}>
      <BrowserRouter>
        <Layout>
          <Routes />
        </Layout>
      </BrowserRouter>
    </AuthContext.Provider>
  ))
}

const container = document.createElement('div')
document.body.appendChild(container)
container.style.height = '100%'
const root = createRoot(container)
root.render(<RootComponent />)
