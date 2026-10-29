import { type ReactNode, useState } from 'react'
import { render, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { Link, MemoryRouter } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import AuthContext from './auth_context'
import ErrorBoundary from './components/error_boundary'
import Layout from './components/layout'
import AppRoutes from './routes'

vi.mock('src/components/error_display', () => ({
  default: ({ err, children }: { err: Error; children: ReactNode }) => (
    <div role="alert">
      {err.message}
      {children}
    </div>
  ),
}))
vi.mock('src/components/layout/nav_bar', () => ({
  default: () => (
    <nav>
      <Link to="/operations">Operations</Link>
    </nav>
  ),
}))
vi.mock('src/pages/admin', () => {
  throw new Error('Failed to download the admin chunk')
})
vi.mock('src/pages/operation_list', () => ({
  default: function OperationList() {
    const [draft, setDraft] = useState('')
    return (
      <>
        <h1>Operation list</h1>
        <input aria-label="Draft" value={draft} onChange={(e) => setDraft(e.target.value)} />
        <Link to="/operations?q=new">Change query</Link>
        <Link to="/admin/users">Admin</Link>
      </>
    )
  },
}))
vi.mock('src/pages/login', () => ({ default: () => <h1>Login</h1> }))

const user = {
  slug: 'reviewer',
  firstName: 'Test',
  lastName: 'User',
  email: 'test@example.com',
  admin: true,
  authSchemes: [],
  headless: false,
}

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
})
afterEach(() => vi.restoreAllMocks())

function renderRoutes(path: string, loggedIn = true) {
  return render(
    <ErrorBoundary>
      <AuthContext.Provider value={{ user: loggedIn ? user : null }}>
        <MemoryRouter initialEntries={[path]}>
          <Layout>
            <AppRoutes />
          </Layout>
        </MemoryRouter>
      </AuthContext.Provider>
    </ErrorBoundary>,
  )
}

describe('page error recovery', () => {
  it('keeps navigation after a rejected lazy import and recovers when leaving the failed route', async () => {
    const browser = userEvent.setup()
    renderRoutes('/admin/users')
    await screen.findByRole('alert')
    expect(screen.getByRole('navigation')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Reload page' })).toBeInTheDocument()
    await browser.click(screen.getByRole('link', { name: 'Operations' }))
    expect(await screen.findByRole('heading', { name: 'Operation list' })).toBeInTheDocument()
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    // React caches the rejected import. Revisiting it must still leave navigation usable.
    await browser.click(screen.getByRole('link', { name: 'Admin' }))
    await screen.findByRole('alert')
    expect(screen.getByRole('navigation')).toBeInTheDocument()
  })

  it('does not remount a healthy page when only the query changes', async () => {
    const browser = userEvent.setup()
    renderRoutes('/operations')
    const draft = await screen.findByRole('textbox', { name: 'Draft' })
    await browser.type(draft, 'Unsaved input')
    await browser.click(screen.getByRole('link', { name: 'Change query' }))
    expect(screen.getByRole('textbox', { name: 'Draft' })).toHaveValue('Unsaved input')
  })

  it('still redirects unauthenticated visitors to the lazy login page', async () => {
    renderRoutes('/operations', false)
    expect(await screen.findByRole('heading', { name: 'Login' })).toBeInTheDocument()
  })
})
