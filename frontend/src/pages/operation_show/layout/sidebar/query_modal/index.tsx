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

export const NewQueryModal = (props: {
  onRequestClose: () => void,
  onCreated: () => void,
  operationSlug: string,
  query: string,
  type: 'evidence'|'findings',
}) => {
  const nameField = useFormField('')
  const queryField = useFormField(props.query)
  const newQueryForm = useForm({
    fields: [nameField, queryField],
    onSuccess: () => {props.onRequestClose(); props.onCreated()},
    handleSubmit: () => saveQuery({
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
  view: ViewName,
}) => {
  const nameField = useFormField(props.savedQuery.name)
  const [searchOptions, setSearchOptions] = React.useState<SearchOptions | null>()

  const editedQueryOrOriginal = () => searchOptions ? stringifySearch(searchOptions) : props.savedQuery.query

  const editQueryForm = useForm({
    fields: [nameField],
    onSuccess: () => {
      props.onRequestClose()
      props.onEdited(props.savedQuery,
        { ...props.savedQuery, name: nameField.value, query: editedQueryOrOriginal() })
    },
    handleSubmit: () => updateSavedQuery({
      operationSlug: props.operationSlug,
      queryId: props.savedQuery.id,
      name: nameField.value,
      query: editedQueryOrOriginal()
    }),
  })

  const wiredData = useWiredData<[Array<Tag>, Array<User>]>(
    React.useCallback(() =>
      Promise.all([
        getTags({ operationSlug: props.operationSlug }),
        listEvidenceCreators({ operationSlug: props.operationSlug }),
      ]), [props.operationSlug]
    ))

  return (
    <ModalForm title="Update Saved Query" submitText="Save" onRequestClose={props.onRequestClose} {...editQueryForm}>
      <Input label="Name" {...nameField} />
      {wiredData.render(([tags, users]) => (
        <FilterFields
          operationSlug={props.operationSlug}
          viewName={props.view}
          allCreators={users}
          searchOptions={stringToSearch(props.savedQuery.query, tags)}
          onChanged={setSearchOptions}
        />
      ))}
    </ModalForm>
  )
}

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
