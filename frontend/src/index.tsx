import { Suspense } from 'react'
import AuthContext from 'src/auth_context'
import ErrorBoundary from 'src/components/error_boundary'
import Layout from 'src/components/layout'
import LoadingSpinner from 'src/components/loading_spinner'
import Routes from 'src/routes'
import { BrowserRouter } from 'react-router'
import { getCurrentUser } from 'src/services'
import { createRoot } from 'react-dom/client'
import { useWiredData } from 'src/helpers'
require('./base_css')

const RootComponent = () => {
  const wiredUser = useWiredData(getCurrentUser)

  return wiredUser.render((user) => (
    <AuthContext.Provider value={{ user }}>
      <BrowserRouter>
        <Layout>
          <Suspense fallback={<LoadingSpinner />}>
            <Routes />
          </Suspense>
        </Layout>
      </BrowserRouter>
    </AuthContext.Provider>
  ))
}

const container = document.createElement('div')
document.body.appendChild(container)
container.style.height = '100%'
const root = createRoot(container)
root.render(
  <ErrorBoundary>
    <RootComponent />
  </ErrorBoundary>,
)
