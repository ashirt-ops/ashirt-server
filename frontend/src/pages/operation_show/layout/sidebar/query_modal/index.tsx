// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { SavedQuery } from 'src/global_types'
import { useDataSource, saveQuery, updateSavedQuery, deleteSavedQuery } from 'src/services'
import { useForm, useFormField } from 'src/helpers'

import Input from  'src/components/input'
import ModalForm from 'src/components/modal_form'

export const NewQueryModal = (props: {
  onRequestClose: () => void,
  onCreated: () => void,
  operationSlug: string,
  query: string,
  type: 'evidence'|'findings',
}) => {
  const ds = useDataSource()
  const nameField = useFormField('')
  const queryField = useFormField(props.query)
  const newQueryForm = useForm({
    fields: [nameField, queryField],
    onSuccess: () => {props.onRequestClose(); props.onCreated()},
    handleSubmit: () => saveQuery(ds, {
      operationSlug: props.operationSlug,
      name: nameField.value,
      query: queryField.value,
      type: props.type,
    }),
  })

  return (
    <ModalForm title="Save New Query" submitText="Save Query" onRequestClose={props.onRequestClose} {...newQueryForm}>
      <Input label="Name" {...nameField} />
      <Input label="Query" {...queryField} />
    </ModalForm>
  )
}

export const EditQueryModal = (props: {
  onEdited: (before: SavedQuery, after: SavedQuery) => void,
  onRequestClose: () => void,
  operationSlug: string,
  savedQuery: SavedQuery,
}) => {
  const ds = useDataSource()
  const nameField = useFormField(props.savedQuery.name)
  const queryField = useFormField(props.savedQuery.query)
  const editQueryForm = useForm({
    fields: [nameField, queryField],
    onSuccess: () => {
      props.onRequestClose()
      props.onEdited(props.savedQuery, {...props.savedQuery, name: nameField.value, query: queryField.value})
    },
    handleSubmit: () => updateSavedQuery(ds, {
      operationSlug: props.operationSlug,
      queryId: props.savedQuery.id,
      name: nameField.value,
      query: queryField.value,
    }),
  })

  return (
    <ModalForm title="Update Saved Query" submitText="Save" onRequestClose={props.onRequestClose} {...editQueryForm}>
      <Input label="Name" {...nameField} />
      <Input label="Query" {...queryField} />
    </ModalForm>
  )
}

export const DeleteQueryModal = (props: {
  onDeleted: (before: SavedQuery) => void,
  onRequestClose: () => void,
  operationSlug: string,
  savedQuery: SavedQuery,
}) => {
  const ds = useDataSource()
  const deleteQueryForm = useForm({
    onSuccess: () => {props.onRequestClose(); props.onDeleted(props.savedQuery)},
    handleSubmit: () => deleteSavedQuery(ds, {
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
