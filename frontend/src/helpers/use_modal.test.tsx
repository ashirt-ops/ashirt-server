import { describe, it, expect, vi } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { renderHook } from '@testing-library/react'
import { useModal, renderModals } from './use_modal'

describe('useModal', () => {
  it('starts with node as null', () => {
    const { result } = renderHook(() => useModal(() => <div>modal</div>))
    expect(result.current.node).toBeNull()
  })

  it('renders the modal node after show()', () => {
    const { result } = renderHook(() => useModal(() => <div>modal content</div>))
    act(() => result.current.show({}))
    const { getByText } = render(<>{result.current.node}</>)
    expect(getByText('modal content')).toBeInTheDocument()
  })

  it('passes props to the modal renderer', () => {
    const { result } = renderHook(() =>
      useModal<{ name: string }>((p) => <div>Hello {p.name}</div>),
    )
    act(() => result.current.show({ name: 'Alice' }))
    const { getByText } = render(<>{result.current.node}</>)
    expect(getByText('Hello Alice')).toBeInTheDocument()
  })

  it('hides the modal after onRequestClose is called', () => {
    let capturedClose: (() => void) | null = null
    const { result } = renderHook(() =>
      useModal<object>((p) => {
        capturedClose = p.onRequestClose
        return <div>modal</div>
      }),
    )
    act(() => result.current.show({}))
    render(<>{result.current.node}</>)
    act(() => capturedClose?.())
    expect(result.current.node).toBeNull()
  })

  it('calls the onClose callback when modal closes', () => {
    const onClose = vi.fn()
    let capturedClose: (() => void) | null = null
    const { result } = renderHook(() =>
      useModal<object>((p) => {
        capturedClose = p.onRequestClose
        return <div>modal</div>
      }, onClose),
    )
    act(() => result.current.show({}))
    render(<>{result.current.node}</>)
    act(() => capturedClose?.())
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('updates when show() is called again with new props', () => {
    const { result } = renderHook(() =>
      useModal<{ name: string }>((p) => <div>Hello {p.name}</div>),
    )
    act(() => result.current.show({ name: 'Alice' }))
    act(() => result.current.show({ name: 'Bob' }))
    const { getByText } = render(<>{result.current.node}</>)
    expect(getByText('Hello Bob')).toBeInTheDocument()
  })
})

describe('renderModals', () => {
  it('returns null when all modals have null nodes', () => {
    const a = { node: null }
    const b = { node: null }
    expect(renderModals(a, b)).toBeNull()
  })

  it('returns the first non-null modal node', () => {
    const a = { node: null }
    const b = { node: <div>first visible</div> }
    const c = { node: <div>second visible</div> }
    const result = renderModals(a, b, c)
    const { getByText } = render(<>{result}</>)
    expect(getByText('first visible')).toBeInTheDocument()
    expect(screen.queryByText('second visible')).toBeNull()
  })
})
