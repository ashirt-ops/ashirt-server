import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Input from './index'

describe('Input', () => {
  it('renders with a label', () => {
    render(<Input label="Email" value="" />)
    expect(screen.getByText('Email')).toBeInTheDocument()
  })

  it('renders the current value', () => {
    render(<Input value="hello@example.com" />)
    expect(screen.getByRole('textbox')).toHaveValue('hello@example.com')
  })

  it('calls onChange with the new value when typed', async () => {
    const user = userEvent.setup()
    const handleChange = vi.fn()
    render(<Input value="" onChange={handleChange} />)
    await user.type(screen.getByRole('textbox'), 'abc')
    expect(handleChange).toHaveBeenCalledWith('a')
    expect(handleChange).toHaveBeenCalledWith('b')
    expect(handleChange).toHaveBeenCalledWith('c')
  })

  it('is disabled when disabled prop is set', () => {
    render(<Input value="" disabled />)
    expect(screen.getByRole('textbox')).toBeDisabled()
  })

  it('renders placeholder text', () => {
    render(<Input value="" placeholder="Search..." />)
    expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument()
  })

  it('is readonly when readOnly prop is set', () => {
    render(<Input value="fixed" readOnly />)
    expect(screen.getByRole('textbox')).toHaveAttribute('readonly')
  })
})
