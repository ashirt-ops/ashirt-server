// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import ModalForm from 'src/components/modal_form'
import { FilterFields } from 'src/components/search_query_builder'
import { SavedQuery, Tag, User, ViewName } from 'src/global_types'
import { saveQuery, updateSavedQuery, deleteSavedQuery, listEvidenceCreators, getTags } from 'src/services'
import { useForm, useFormField, useWiredData } from 'src/helpers'
import { stringToSearch, SearchOptions, stringifySearch } from 'src/components/search_query_builder/helpers'

const SaveOrEditModal = (props: {
  modalTitle: string
  submitButtonName: string
  name?: string,
  query: string,
  operationSlug: string,
  view: ViewName,
  onRequestClose: () => void,
  onComplete: (name: string, query: string) => void,
  saveFn: (name: string, query: string) => Promise<void>,
}) => {
  const nameField = useFormField(props.name || '')
  const [searchOptions, setSearchOptions] = React.useState<SearchOptions | null>(null)
  const editedQueryOrOriginal = () => searchOptions ? stringifySearch(searchOptions) : props.query

  const queryForm = useForm({
    fields: [nameField],
    onSuccess: () => {
      props.onRequestClose()
      props.onComplete(nameField.value, editedQueryOrOriginal())
    },
    handleSubmit: () =>  props.saveFn(nameField.value, editedQueryOrOriginal())
  })

  const wiredData = useWiredData<[Array<Tag>, Array<User>]>(
    React.useCallback(() =>
      Promise.all([
        getTags({ operationSlug: props.operationSlug }),
        listEvidenceCreators({ operationSlug: props.operationSlug }),
      ]), [props.operationSlug]
    ))
  return (
    <ModalForm title={props.modalTitle} submitText={props.submitButtonName} onRequestClose={props.onRequestClose} {...queryForm}>
      <Input label="Name" {...nameField} />
      {wiredData.render(([tags, users]) => (
        <FilterFields
          operationSlug={props.operationSlug}
          viewName={props.view}
          allCreators={users}
          searchOptions={searchOptions || stringToSearch(props.query, tags)}
          onChanged={setSearchOptions}
        />
      ))}
    </ModalForm>
  )
}

export const NewQueryModal = (props: {
  onRequestClose: () => void,
  onCreated: () => void,
  operationSlug: string,
  query: string,
  type: 'evidence' | 'findings',
}) => (
  <SaveOrEditModal
    modalTitle="Save New Query"
    submitButtonName="Save Query"
    query={props.query}
    view={props.type}
    onRequestClose={props.onRequestClose}
    operationSlug={props.operationSlug}
    onComplete={() => props.onCreated()}
    saveFn={(name, query) =>
      saveQuery({
        operationSlug: props.operationSlug,
        name,
        query,
        type: props.type,
      })}
  />
)

export const EditQueryModal = (props: {
  onEdited: (before: SavedQuery, after: SavedQuery) => void,
  onRequestClose: () => void,
  operationSlug: string,
  savedQuery: SavedQuery,
  view: ViewName,
}) => (
  <SaveOrEditModal
    modalTitle="Update Saved Query"
    submitButtonName="Save"
    name={props.savedQuery.name}
    query={props.savedQuery.query}
    view={props.view}
    onRequestClose={props.onRequestClose}
    operationSlug={props.operationSlug}
    onComplete={(name, query) => props.onEdited(props.savedQuery, { ...props.savedQuery, name, query })}
    saveFn={(name, query) =>
      updateSavedQuery({
        operationSlug: props.operationSlug,
        queryId: props.savedQuery.id,
        name,
        query,
      })}
  />
)

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
