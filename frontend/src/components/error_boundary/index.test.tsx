import { type ReactNode } from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ErrorBoundary from './index'

vi.mock('src/components/error_display', () => ({
  default: ({ err, children }: { err: Error; children: ReactNode }) => (
    <div role="alert">
      {err.message}
      {children}
    </div>
  ),
}))

function Broken({ error }: { error: unknown }): never {
  throw error
}

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

describe('ErrorBoundary', () => {
  it.each([new Error('Render failed'), 'Render failed', null])(
    'provides full-page recovery even for a non-Error throw: %s',
    (error) => {
      render(
        <ErrorBoundary>
          <Broken error={error} />
        </ErrorBoundary>,
      )
      expect(screen.getByRole('alert')).toHaveTextContent(
        error instanceof Error ? error.message : String(error),
      )
      const reload = vi.fn()
      vi.stubGlobal('window', { location: { reload } })
      fireEvent.click(screen.getByRole('button', { name: 'Reload page' }))
      expect(reload).toHaveBeenCalledTimes(1)
    },
  )

  it('recovers on navigation without retrying during unrelated rerenders', () => {
    const view = render(
      <ErrorBoundary resetKey="first">
        <Broken error={new Error('Broken')} />
      </ErrorBoundary>,
    )
    view.rerender(
      <ErrorBoundary resetKey="first">
        <div>Healthy</div>
      </ErrorBoundary>,
    )
    expect(screen.queryByText('Healthy')).not.toBeInTheDocument()
    view.rerender(
      <ErrorBoundary resetKey="second">
        <div>Healthy</div>
      </ErrorBoundary>,
    )
    expect(screen.getByText('Healthy')).toBeInTheDocument()
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })
})
