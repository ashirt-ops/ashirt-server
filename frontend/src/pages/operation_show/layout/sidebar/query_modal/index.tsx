import * as React from 'react'
import ModalForm from 'src/components/modal_form'
import { SavedQuery } from 'src/global_types'
import { deleteSavedQuery } from 'src/services'
import { useForm } from 'src/helpers'

export const DeleteQueryModal = (props: {
  onDeleted: (before: SavedQuery) => void,
  onRequestClose: () => void,
  operationSlug: string,
  savedQuery: SavedQuery,
}) => {
  const deleteQueryForm = useForm({
    onSuccess: () => {props.onRequestClose(); props.onDeleted(props.savedQuery)},
    handleSubmit: () => deleteSavedQuery({
      operationSlug: props.operationSlug,
      queryId: props.savedQuery.id,
    }),
  })

  return (
    <ModalForm title="Delete Saved Query" submitText="Delete Query" onRequestClose={props.onRequestClose} {...deleteQueryForm}>
      <p>Are you sure you want to delete the saved query "{props.savedQuery.name}"?</p>
    </ModalForm>
  )
}
