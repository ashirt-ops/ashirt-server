import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Pager, { StandardPager } from './index'

describe('Pager', () => {
  it('renders children', () => {
    render(
      <Pager pageNumber={1} onPageUp={() => {}} onPageDown={() => {}}>
        <span>Page 1</span>
      </Pager>,
    )
    expect(screen.getByText('Page 1')).toBeInTheDocument()
  })

  it('renders previous and next buttons', () => {
    render(
      <Pager pageNumber={2} onPageUp={() => {}} onPageDown={() => {}}>
        <span>2</span>
      </Pager>,
    )
    expect(screen.getByRole('button', { name: 'previous' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'next' })).toBeInTheDocument()
  })

  it('disables previous button on page 1', () => {
    render(
      <Pager pageNumber={1} onPageUp={() => {}} onPageDown={() => {}}>
        <span>1</span>
      </Pager>,
    )
    expect(screen.getByRole('button', { name: 'previous' })).toBeDisabled()
  })

  it('enables previous button on page > 1', () => {
    render(
      <Pager pageNumber={2} onPageUp={() => {}} onPageDown={() => {}}>
        <span>2</span>
      </Pager>,
    )
    expect(screen.getByRole('button', { name: 'previous' })).not.toBeDisabled()
  })

  it('disables next button when at max page', () => {
    render(
      <Pager pageNumber={3} maxPageNumber={3} onPageUp={() => {}} onPageDown={() => {}}>
        <span>3</span>
      </Pager>,
    )
    expect(screen.getByRole('button', { name: 'next' })).toBeDisabled()
  })

  it('enables next button when below max page', () => {
    render(
      <Pager pageNumber={2} maxPageNumber={5} onPageUp={() => {}} onPageDown={() => {}}>
        <span>2</span>
      </Pager>,
    )
    expect(screen.getByRole('button', { name: 'next' })).not.toBeDisabled()
  })

  it('calls onPageUp when next is clicked', async () => {
    const onPageUp = vi.fn()
    const user = userEvent.setup()
    render(
      <Pager pageNumber={1} onPageUp={onPageUp} onPageDown={() => {}}>
        <span>1</span>
      </Pager>,
    )
    await user.click(screen.getByRole('button', { name: 'next' }))
    expect(onPageUp).toHaveBeenCalled()
  })

  it('calls onPageDown when previous is clicked', async () => {
    const onPageDown = vi.fn()
    const user = userEvent.setup()
    render(
      <Pager pageNumber={2} onPageUp={() => {}} onPageDown={onPageDown}>
        <span>2</span>
      </Pager>,
    )
    await user.click(screen.getByRole('button', { name: 'previous' }))
    expect(onPageDown).toHaveBeenCalled()
  })
})

describe('StandardPager', () => {
  it('renders the current page number', () => {
    render(<StandardPager page={3} onPageChange={() => {}} />)
    expect(screen.getByText('3')).toBeInTheDocument()
  })

  it('calls onPageChange with incremented page when next is clicked', async () => {
    const onPageChange = vi.fn()
    const user = userEvent.setup()
    render(<StandardPager page={2} maxPages={5} onPageChange={onPageChange} />)
    await user.click(screen.getByRole('button', { name: 'next' }))
    expect(onPageChange).toHaveBeenCalledWith(3)
  })

  it('calls onPageChange with decremented page when previous is clicked', async () => {
    const onPageChange = vi.fn()
    const user = userEvent.setup()
    render(<StandardPager page={3} onPageChange={onPageChange} />)
    await user.click(screen.getByRole('button', { name: 'previous' }))
    expect(onPageChange).toHaveBeenCalledWith(2)
  })

  it('disables previous button on page 1', () => {
    render(<StandardPager page={1} onPageChange={() => {}} />)
    expect(screen.getByRole('button', { name: 'previous' })).toBeDisabled()
  })

  it('disables next when at max pages', () => {
    render(<StandardPager page={5} maxPages={5} onPageChange={() => {}} />)
    expect(screen.getByRole('button', { name: 'next' })).toBeDisabled()
  })
})
