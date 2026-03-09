import { type ReactNode, type FormEvent } from 'react'
import Form from 'src/components/form'
import Modal from 'src/components/modal'
import { type Result } from 'src/global_types'

// A convenience wrapper around <Modal><Form>{...}</Form></Modal> since it's such a common UI piece
const ModalForm = (props: {
  children: ReactNode
  result: Result<string> | null
  loading: boolean
  onRequestClose: () => void
  onSubmit: (e: FormEvent) => void
  submitText?: string
  title: string
  submitDanger?: boolean
  cancelText?: string
  disableSubmit?: boolean
  disableCancel?: boolean
}) => (
  <Modal title={props.title} onRequestClose={props.onRequestClose}>
    <Form
      cancelText={props.cancelText || 'Cancel'}
      children={props.children}
      result={props.result}
      loading={props.loading}
      onCancel={props.onRequestClose}
      onSubmit={props.onSubmit}
      submitText={props.submitText}
      submitDanger={props.submitDanger}
      disableSubmit={props.disableSubmit}
      disableCancel={props.disableCancel}
    />
  </Modal>
)
export default ModalForm
