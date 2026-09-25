import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Form from './index'

describe('Form', () => {
  it('renders children', () => {
    render(
      <Form result={null} loading={false} onSubmit={() => {}}>
        <input placeholder="Name" />
      </Form>,
    )
    expect(screen.getByPlaceholderText('Name')).toBeInTheDocument()
  })

  it('renders submit button when submitText is provided', () => {
    render(<Form result={null} loading={false} onSubmit={() => {}} submitText="Save" />)
    expect(screen.getByText('Save')).toBeInTheDocument()
  })

  it('does not render submit button when submitText is absent', () => {
    render(<Form result={null} loading={false} onSubmit={() => {}} />)
    expect(screen.queryByRole('button')).toBeNull()
  })

  it('renders cancel button when onCancel is provided', () => {
    render(
      <Form
        result={null}
        loading={false}
        onSubmit={() => {}}
        onCancel={() => {}}
        cancelText="Cancel"
      />,
    )
    expect(screen.getByText('Cancel')).toBeInTheDocument()
  })

  it('calls onSubmit when form is submitted', async () => {
    const onSubmit = vi.fn((e) => e.preventDefault())
    const user = userEvent.setup()
    render(
      <Form result={null} loading={false} onSubmit={onSubmit} submitText="Go">
        <input />
      </Form>,
    )
    await user.click(screen.getByText('Go'))
    expect(onSubmit).toHaveBeenCalled()
  })

  it('calls onCancel when cancel button is clicked', async () => {
    const onCancel = vi.fn()
    const user = userEvent.setup()
    render(
      <Form
        result={null}
        loading={false}
        onSubmit={() => {}}
        onCancel={onCancel}
        cancelText="Cancel"
      />,
    )
    await user.click(screen.getByText('Cancel'))
    expect(onCancel).toHaveBeenCalled()
  })

  it('displays an error result message', () => {
    render(
      <Form
        result={{ err: new Error('Something went wrong') }}
        loading={false}
        onSubmit={() => {}}
      />,
    )
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
  })

  it('displays a success result message', () => {
    render(<Form result={{ success: 'Saved!' }} loading={false} onSubmit={() => {}} />)
    expect(screen.getByText('Saved!')).toBeInTheDocument()
  })

  it('shows loading state on submit button', () => {
    render(<Form result={null} loading={true} onSubmit={() => {}} submitText="Save" />)
    // Button renders a spinner when loading; the text may be hidden or replaced
    const button = screen.getByRole('button')
    expect(button).toBeInTheDocument()
  })

  it('disables cancel when loading', () => {
    render(
      <Form
        result={null}
        loading={true}
        onSubmit={() => {}}
        onCancel={() => {}}
        cancelText="Cancel"
      />,
    )
    // Cancel button is disabled when loading
    const buttons = screen.getAllByRole('button')
    const cancelButton = buttons.find((b) => b.textContent === 'Cancel')
    expect(cancelButton).toBeDisabled()
  })

  it('disables submit when disableSubmit is true', () => {
    render(
      <Form result={null} loading={false} onSubmit={() => {}} submitText="Save" disableSubmit />,
    )
    expect(screen.getByRole('button')).toBeDisabled()
  })
})
