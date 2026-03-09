import { useState, type FormEvent } from 'react'
import { type Result } from 'src/global_types'

// Use form react hook helper to make form loading/error handling easier
// Similar to useWiredData except for submitting data
//
// handleSubmit: (required)
//   An async handler function to actually send the data to the API
//   Rejected promises are displayed to the user
//
// fields:
//   An array of useFormField return values to tie into the form.
//   This is optional and only used to disable the fields on submit
//
// onSuccess:
//   Optional calback to call after a successful handleSubmit
//
// Example:
//
// const nameField = useFormField("Alice")
// const emailField = useFormField("alice@example.com")
// const formProps = useForm(
//   [nameField, emailField],
//   () => postToApi({name: nameField.value, email: emailField.value}),
// )
//
// return (
//   <form {...formProps}>
//     {formProps.result && <Error>{formProps.result}</Error>}
//     <input {...nameField} />
//     <input {...emailField} />
//   </form>
// )

export function useForm<T>(i: {
  handleSubmit: () => Promise<T>
  fields?: Array<{ setDisabled: (v: boolean) => void }>
  onSuccessText?: string
  onSuccess?: () => void
}) {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<Result<string> | null>(null)

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (i.fields) i.fields.forEach((f) => f.setDisabled(true))
    setLoading(true)
    let submitSuccessful
    try {
      await i.handleSubmit()
      setResult(null)
      submitSuccessful = true
    } catch (err) {
      setResult({ err: err instanceof Error ? err : new Error(String(err)) })
      submitSuccessful = false
    }
    if (i.fields) i.fields.forEach((f) => f.setDisabled(false))
    setLoading(false)
    if (submitSuccessful) {
      if (i.onSuccessText) setResult({ success: i.onSuccessText })
      if (i.onSuccess) i.onSuccess()
    }
  }

  return { onSubmit, loading, result }
}

export function useFormField<T>(initialValue: T) {
  const [value, onChange] = useState<T>(initialValue)
  const [disabled, setDisabled] = useState(false)

  return { value, onChange, disabled, setDisabled }
}
