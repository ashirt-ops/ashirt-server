import { describe, it, expect, vi } from 'vitest'
import { render, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Modal from './index'

describe('Modal', () => {
  it('renders the title', () => {
    render(
      <Modal title="My Modal" onRequestClose={() => {}}>
        {'.'}
      </Modal>,
    )
    expect(screen.getByText('My Modal')).toBeInTheDocument()
  })

  it('renders children', () => {
    render(
      <Modal title="Test" onRequestClose={() => {}}>
        <p>Modal content here</p>
      </Modal>,
    )
    expect(screen.getByText('Modal content here')).toBeInTheDocument()
  })

  it('renders via a portal into document.body', () => {
    const { baseElement } = render(
      <Modal title="Portal Test" onRequestClose={() => {}}>
        <span>portal child</span>
      </Modal>,
    )
    expect(within(baseElement).getByText('Portal Test')).toBeInTheDocument()
  })

  it('calls onRequestClose when the backdrop is clicked', async () => {
    const user = userEvent.setup()
    const onRequestClose = vi.fn()
    render(
      <Modal title="Close test" onRequestClose={onRequestClose}>
        <button>inside</button>
      </Modal>,
    )
    const backdrop = document.body.querySelector('[aria-modal]')?.parentElement
    if (backdrop) {
      await user.pointer({ target: backdrop, keys: '[MouseLeft]' })
    }
    expect(onRequestClose).toHaveBeenCalled()
  })

  it('does not call onRequestClose when the inner modal is clicked', async () => {
    const user = userEvent.setup()
    const onRequestClose = vi.fn()
    render(
      <Modal title="No close" onRequestClose={onRequestClose}>
        <button>inside modal</button>
      </Modal>,
    )
    await user.click(screen.getByText('No close'))
    expect(onRequestClose).not.toHaveBeenCalled()
  })

  it('has role="dialog" and aria-modal on the dialog element', () => {
    render(
      <Modal title="ARIA test" onRequestClose={() => {}}>
        {'.'}
      </Modal>,
    )
    const dialog = screen.getByRole('dialog')
    expect(dialog).toHaveAttribute('aria-modal', 'true')
  })

  it('dialog is labelled by the title', () => {
    render(
      <Modal title="Labelled Modal" onRequestClose={() => {}}>
        {'.'}
      </Modal>,
    )
    const dialog = screen.getByRole('dialog', { name: 'Labelled Modal' })
    expect(dialog).toBeInTheDocument()
  })
})
