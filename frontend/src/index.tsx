// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

require('./base_css')
import * as React from 'react'
import AuthContext from 'src/auth_context'
import Layout from 'src/components/layout'
import Routes from 'src/routes'
import {BrowserRouter} from 'react-router-dom'
import {getCurrentUser} from 'src/services'
import {render} from 'react-dom'
import {useWiredData} from 'src/helpers'

const RootComponent = () => {
  const wiredUser = useWiredData(getCurrentUser)

  return wiredUser.render(user => (
    <AuthContext.Provider value={{user}}>
      <BrowserRouter>
        <Layout>
          <Routes />
        </Layout>
      </BrowserRouter>
    </AuthContext.Provider>
  ))
}

const root = document.createElement('div')
document.body.appendChild(root)
root.style.height = '100%'
render(<RootComponent /> , root)
