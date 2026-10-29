import { describe, it, expect, vi } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { type FormEvent } from 'react'
import { useForm, useFormField } from './use_form'

const mockEvent = { preventDefault: vi.fn() } as unknown as FormEvent

describe('useFormField', () => {
  it('initializes with the given value', () => {
    const { result } = renderHook(() => useFormField('hello'))
    expect(result.current.value).toBe('hello')
  })

  it('updates value via onChange', () => {
    const { result } = renderHook(() => useFormField('hello'))
    act(() => result.current.onChange('world'))
    expect(result.current.value).toBe('world')
  })

  it('starts enabled', () => {
    const { result } = renderHook(() => useFormField(''))
    expect(result.current.disabled).toBe(false)
  })

  it('can be disabled via setDisabled', () => {
    const { result } = renderHook(() => useFormField(''))
    act(() => result.current.setDisabled(true))
    expect(result.current.disabled).toBe(true)
  })
})

describe('useForm', () => {
  it('starts not loading', () => {
    const { result } = renderHook(() => useForm({ handleSubmit: () => Promise.resolve() }))
    expect(result.current.loading).toBe(false)
  })

  it('starts with null result', () => {
    const { result } = renderHook(() => useForm({ handleSubmit: () => Promise.resolve() }))
    expect(result.current.result).toBeNull()
  })

  it('calls handleSubmit on submit', async () => {
    const handleSubmit = vi.fn(() => Promise.resolve())
    const { result } = renderHook(() => useForm({ handleSubmit }))
    await act(() => result.current.onSubmit(mockEvent))
    expect(handleSubmit).toHaveBeenCalledTimes(1)
  })

  it('sets result.err when handleSubmit rejects with an Error', async () => {
    const handleSubmit = vi.fn(() => Promise.reject(new Error('fail')))
    const { result } = renderHook(() => useForm({ handleSubmit }))
    await act(() => result.current.onSubmit(mockEvent))
    expect(result.current.result).toMatchObject({ err: expect.any(Error) })
  })

  it('sets result.err for non-Error throws', async () => {
    const handleSubmit = vi.fn(() => Promise.reject('string error'))
    const { result } = renderHook(() => useForm({ handleSubmit }))
    await act(() => result.current.onSubmit(mockEvent))
    expect(result.current.result).toMatchObject({ err: expect.any(Error) })
  })

  it('sets result.success after success when onSuccessText is provided', async () => {
    const { result } = renderHook(() =>
      useForm({ handleSubmit: () => Promise.resolve(), onSuccessText: 'Saved!' }),
    )
    await act(() => result.current.onSubmit(mockEvent))
    expect(result.current.result).toMatchObject({ success: 'Saved!' })
  })

  it('calls onSuccess callback after successful submit', async () => {
    const onSuccess = vi.fn()
    const { result } = renderHook(() =>
      useForm({ handleSubmit: () => Promise.resolve(), onSuccess }),
    )
    await act(() => result.current.onSubmit(mockEvent))
    expect(onSuccess).toHaveBeenCalledTimes(1)
  })

  it('disables and re-enables fields during submit', async () => {
    let resolveSubmit!: () => void
    const handleSubmit = () =>
      new Promise<void>((r) => {
        resolveSubmit = r
      })
    const field = { setDisabled: vi.fn() }
    const { result } = renderHook(() => useForm({ handleSubmit, fields: [field] }))

    const submitPromise = act(() => result.current.onSubmit(mockEvent))
    expect(field.setDisabled).toHaveBeenCalledWith(true)

    act(() => resolveSubmit())
    await submitPromise
    expect(field.setDisabled).toHaveBeenCalledWith(false)
  })

  it('is not loading before submit and is loading after submit resolves', async () => {
    const { result } = renderHook(() => useForm({ handleSubmit: () => Promise.resolve() }))
    expect(result.current.loading).toBe(false)
    await act(() => result.current.onSubmit(mockEvent))
    expect(result.current.loading).toBe(false)
  })
})
