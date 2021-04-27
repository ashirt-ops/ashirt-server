// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Input from 'src/components/input'
import ModalForm from 'src/components/modal_form'
import { FindingCategory } from 'src/global_types'
import { createFindingCategory, deleteFindingCategory, updateFindingCategory } from 'src/services'
import { useForm, useFormField } from 'src/helpers'

export const DeleteFindingCategoryModal = (props: {
  onDeleted: () => void,
  onRequestClose: () => void,
  category: FindingCategory,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => deleteFindingCategory({
      id: props.category.id,
      delete: true
    }),
  })

  const affectedText = props.category.usageCount == 0 
    ? "" 
    : props.category.usageCount == 1
      ? "This will affect one finding."
      : `This will affect ${props.category.usageCount} findings.`

  return (
    <ModalForm title="Delete Finding Category"
      submitDanger submitText="Delete"
      cancelText="Close"
      onRequestClose={props.onRequestClose} {...formComponentProps}>
      <p>
        Are you sure you want to delete this finding category? {affectedText}
      </p>
    </ModalForm>
  )
}

export const EditFindingCategoryModal = (props: {
  onEdited: () => void,
  onRequestClose: () => void,
  category?: FindingCategory,
}) => {
  const nameField = useFormField<string>(props.category?.category || "")
  const isUpdate = props.category !== undefined

  const formComponentProps = useForm({
    fields: [nameField],
    onSuccess: () => { props.onEdited(); props.onRequestClose() },
    handleSubmit: () => {
      if (nameField.value.length == 0) {
        return Promise.reject(new Error("Please provide a name for the the category"))
      }

      if (isUpdate) {
        return updateFindingCategory({
          findingCategoryId: props.category!.id,
          category: nameField.value.trim(),
        })
      }
      else {
        return (async () => {await createFindingCategory(nameField.value)})()
      }
    },
  })

  const modalProps = isUpdate
    ? {
      title: "Edit Finding Category",
      submitText: "Save"
    }
    : {
      title: "Create Finding Category",
      submitText: "Create"
    }

  return (
    <ModalForm {...modalProps}
      cancelText="Close"
      onRequestClose={props.onRequestClose} {...formComponentProps}>
      <Input label="Name" {...nameField} />
    </ModalForm>
  )
}
