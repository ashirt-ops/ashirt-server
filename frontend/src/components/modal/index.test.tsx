import { describe, it, expect, vi } from 'vitest'
import { render, screen, within, fireEvent } from '@testing-library/react'
import { StrictMode } from 'react'
import { userEvent } from '@testing-library/user-event'
import Modal from './index'
import Popover from 'src/components/popover'

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
    await user.pointer({ target: screen.getByRole('dialog'), keys: '[MouseLeft]' })
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

  it('prevents the backdrop from taking focus after dismissal', () => {
    const onRequestClose = vi.fn()
    render(
      <Modal title="Backdrop" onRequestClose={onRequestClose}>
        <button>Inside</button>
      </Modal>,
    )
    expect(fireEvent.mouseDown(screen.getByRole('dialog'))).toBe(false)
    expect(onRequestClose).toHaveBeenCalledTimes(1)
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

  it('opens as a native modal and closes on unmount', () => {
    const show = vi.spyOn(HTMLDialogElement.prototype, 'showModal')
    const close = vi.spyOn(HTMLDialogElement.prototype, 'close')
    try {
      const view = render(
        <StrictMode>
          <Modal title="Native modal" onRequestClose={() => {}}>
            <button>Close</button>
          </Modal>
        </StrictMode>,
      )
      const dialog = screen.getByRole('dialog')
      expect(dialog.tagName).toBe('DIALOG')
      expect(dialog).toHaveAttribute('open')
      expect(show).toHaveBeenCalled()
      view.unmount()
      expect(close).toHaveBeenCalled()
      expect(dialog).not.toHaveAttribute('open')
    } finally {
      show.mockRestore()
      close.mockRestore()
    }
  })

  it('requests close on native cancellation without closing behind the parent state', () => {
    const onRequestClose = vi.fn()
    render(
      <Modal title="Cancel" onRequestClose={onRequestClose}>
        <button>Inside</button>
      </Modal>,
    )
    const dialog = screen.getByRole('dialog')
    const cancel = new Event('cancel', { cancelable: true })
    fireEvent(dialog, cancel)
    expect(onRequestClose).toHaveBeenCalledTimes(1)
    expect(cancel.defaultPrevented).toBe(true)
    expect(dialog).toHaveAttribute('open')
  })

  it('keeps popup controls inside the modal and does not dismiss it when selecting them', async () => {
    const user = userEvent.setup()
    const onRequestClose = vi.fn()
    const onSelect = vi.fn()
    render(
      <Modal title="Picker" onRequestClose={onRequestClose}>
        <Popover isOpen content={<button onClick={onSelect}>Popup option</button>}>
          <button>Open picker</button>
        </Popover>
      </Modal>,
    )
    const option = within(screen.getByRole('dialog')).getByRole('button', { name: 'Popup option' })
    await user.click(option)
    expect(onSelect).toHaveBeenCalledTimes(1)
    expect(onRequestClose).not.toHaveBeenCalled()
  })
})
