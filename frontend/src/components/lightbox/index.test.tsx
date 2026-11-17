import { useState } from 'react'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import Modal from 'src/components/modal'
import Lightbox from './index'

describe('Lightbox', () => {
  it('keeps modal previews inside the dialog and closes only the preview on Escape', async () => {
    const closeModal = vi.fn()
    function Preview() {
      const [open, setOpen] = useState(true)
      return (
        <Modal title="Select evidence" onRequestClose={closeModal}>
          <Lightbox isOpen={open} onRequestClose={() => setOpen(false)}>
            <button>Preview content</button>
          </Lightbox>
        </Modal>
      )
    }
    render(<Preview />)
    const dialog = screen.getByRole('dialog')
    const preview = await within(dialog).findByRole('button', { name: 'Preview content' })
    expect(fireEvent.keyDown(preview, { key: 'Escape' })).toBe(false)
    await waitFor(() => expect(screen.queryByText('Preview content')).not.toBeInTheDocument())
    expect(dialog).toHaveAttribute('open')
    expect(closeModal).not.toHaveBeenCalled()
    expect(fireEvent.keyDown(dialog, { key: 'Escape' })).toBe(true)
  })

  it('portals standalone previews to the body', async () => {
    const { container } = render(
      <Lightbox isOpen onRequestClose={() => {}}>
        <button>Standalone preview</button>
      </Lightbox>,
    )
    const preview = await screen.findByRole('button', { name: 'Standalone preview' })
    expect(document.body).toContainElement(preview)
    expect(container).not.toContainElement(preview)
  })
})
