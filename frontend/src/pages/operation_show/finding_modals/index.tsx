// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import EvidenceChooser from 'src/components/evidence_chooser'
import ModalForm from 'src/components/modal_form'
import Select from 'src/components/select'
import { Evidence, Finding } from 'src/global_types'
import { createFinding, removeEvidenceFromFinding, updateFinding, deleteFinding, changeEvidenceOfFinding, getFindingCategories } from 'src/services'
import { default as Input, TextArea } from 'src/components/input'
import { useForm, useFormField, useWiredData } from 'src/helpers'
import Checkbox from 'src/components/checkbox'

const CategorySelect = (props: {
  disabled: boolean,
  onChange: (v: string) => void,
  value: string,
}) => {
  const wiredCategories = useWiredData(getFindingCategories)
  return wiredCategories.render(categories => (
    <Select label="Category" {...props}>
      <option value="">- Select a category -</option>
      {categories.map(category => (
        <option key={category}>{category}</option>
      ))}
    </Select>
  ))
}

export const CreateFindingModal = (props: {
  fromEvidence?: Evidence,
  onCreated: (f: Finding) => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const categoryField = useFormField<string>('')
  const titleField = useFormField<string>('')
  const descriptionField = useFormField<string>('')
  const formComponentProps = useForm({
    fields: [categoryField, titleField, descriptionField],
    onSuccess: () => props.onRequestClose(),
    handleSubmit: async () => {
      const finding = await createFinding({
        operationSlug: props.operationSlug,
        category: categoryField.value,
        title: titleField.value,
        description: descriptionField.value,
      })
      props.onCreated(finding)
    },
  })

  return (
    <ModalForm title="New Finding" submitText="Create Finding" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <Input label="Title" {...titleField} />
      <CategorySelect {...categoryField} />
      <TextArea label="Description" {...descriptionField} />
    </ModalForm>
  )
}

export const EditFindingModal = (props: {
  finding: Finding,
  onEdited: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const categoryField = useFormField<string>(props.finding.category)
  const titleField = useFormField<string>(props.finding.title)
  const ticketField = useFormField<string>(props.finding.ticketLink || "")
  const descriptionField = useFormField<string>(props.finding.description)
  const readyToReportField = useFormField(props.finding.readyToReport)
  const formComponentProps = useForm({
    fields: [categoryField, titleField, descriptionField],
    onSuccess: () => { props.onEdited(); props.onRequestClose() },
    handleSubmit: () => updateFinding({
      operationSlug: props.operationSlug,
      findingUuid: props.finding.uuid,
      category: categoryField.value,
      title: titleField.value,
      description: descriptionField.value,
      readyToReport: readyToReportField.value,
      ticketLink: ticketField.value === "" ? null : ticketField.value,
    }),
  })
  return (
    <ModalForm title="Edit Finding" submitText="Save" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <Input label="Title" {...titleField} />
      <CategorySelect {...categoryField} />
      <Checkbox label="Ready to Report" {...readyToReportField}/>
      <Input label="Ticket URL" {...ticketField} disabled={!readyToReportField.value}/>
      <TextArea label="Description" {...descriptionField} />
    </ModalForm>
  )
}

export const DeleteFindingModal = (props: {
  finding: Finding,
  onDeleted: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => deleteFinding({
      findingUuid: props.finding.uuid,
      operationSlug: props.operationSlug,
    }),
  })
  return (
    <ModalForm title="Delete Finding" submitText="Delete Finding" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <p>Are you sure you want to delete this finding?</p>
    </ModalForm>
  )
}

export const ChangeEvidenceOfFindingModal = (props: {
  finding: Finding,
  initialEvidence: Array<Evidence>,
  onChanged: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const evidenceField = useFormField<Array<Evidence>>(props.initialEvidence)
  const formComponentProps = useForm({
    fields: [evidenceField],
    onSuccess: () => { props.onChanged(); props.onRequestClose() },
    handleSubmit: () => changeEvidenceOfFinding({
      operationSlug: props.operationSlug,
      findingUuid: props.finding.uuid,
      oldEvidence: props.initialEvidence,
      newEvidence: evidenceField.value,
    }),
  })

  return (
    <ModalForm title="Select Evidence For Finding" submitText="Save" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <EvidenceChooser operationSlug={props.operationSlug} {...evidenceField} />
    </ModalForm>
  )
}

export const RemoveEvidenceFromFindingModal = (props: {
  finding: Finding,
  evidence: Evidence,
  onRemoved: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onRemoved(); props.onRequestClose() },
    handleSubmit: () => removeEvidenceFromFinding({
      evidenceUuid: props.evidence.uuid,
      findingUuid: props.finding.uuid,
      operationSlug: props.operationSlug,
    }),
  })
  return (
    <ModalForm title="Remove Evidence From Finding" submitText="Remove Evidence" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <p>Are you sure you want to remove the selected evidence from this finding?</p>
    </ModalForm>
  )
}
