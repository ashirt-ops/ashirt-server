// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Checkbox from 'src/components/checkbox'
import FindingChooser from 'src/components/finding_chooser'
import Form from 'src/components/form'
import ImageUpload from 'src/components/image_upload'
import ModalForm from 'src/components/modal_form'
import Modal from 'src/components/modal'
import TagChooser from 'src/components/tag_chooser'
import TerminalRecordingUpload from 'src/components/termrec_upload'
import { CodeBlockEditor } from 'src/components/code_block'
import { Evidence, Finding, Tag, CodeBlock, SubmittableEvidence } from 'src/global_types'
import { TextArea } from 'src/components/input'
import { default as TabMenu, Tab } from 'src/components/tabs'
import { useForm, useFormField } from 'src/helpers/use_form'
import { codeblockToBlob } from 'src/helpers/codeblock_to_blob'
import { useWiredData } from 'src/helpers'

import { createEvidence, updateEvidence, deleteEvidence, changeFindingsOfEvidence, getFindingsOfEvidence, getEvidenceAsCodeblock } from 'src/services'

export const CreateEvidenceModal = (props: {
  onCreated: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const descriptionField = useFormField<string>("")
  const tagsField = useFormField<Array<Tag>>([])
  const binaryBlobField = useFormField<File | null>(null)
  const codeblockField = useFormField<CodeBlock>({ type: 'codeblock', language: '', code: '', source: null })

  const evidenceTypes: Array<Tab> = [
    { id: 'screenshot', label: 'Screenshot', content: <ImageUpload label='Screenshot' {...binaryBlobField} /> },
    { id: 'codeblock', label: 'Code Block', content: <CodeBlockEditor {...codeblockField} /> },
    { id: 'terminal-recording', label: 'Terminal Recording', content: <TerminalRecordingUpload label='Terminal Recording' {...binaryBlobField} /> },
  ]

  const [selectedTab, setSelectedTab] = React.useState<Tab>(evidenceTypes[0])

  const formComponentProps = useForm({
    fields: [descriptionField, binaryBlobField],
    onSuccess: () => { props.onCreated(); props.onRequestClose() },
    handleSubmit: () => {
      let data: SubmittableEvidence = { type: "none" }

      if (selectedTab.id === 'screenshot' && binaryBlobField.value != null) {
        data = { type: 'image', file: binaryBlobField.value }
      } else if (selectedTab.id === 'codeblock' && codeblockField.value !== null) {
        data = { type: 'codeblock', file: codeblockToBlob(codeblockField.value) }
      } else if (selectedTab.id === 'terminal-recording' && binaryBlobField.value !== null) {
        data = { type: 'terminal-recording', file: binaryBlobField.value }
      }

      return createEvidence({
        operationSlug: props.operationSlug,
        description: descriptionField.value,
        evidence: data,
        tagIds: tagsField.value.map(t => t.id),
      })
    },
  })
  return (
    <ModalForm title="New Evidence" submitText="Create Evidence" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <TextArea label="Description" {...descriptionField} />
      <TabMenu
        onTabChanged={(tab, tabIndex) => { setSelectedTab(tab) }}
        tabs={evidenceTypes}
      />
      <TagChooser operationSlug={props.operationSlug} label="Tags" {...tagsField} />
    </ModalForm>
  )
}
/*
const wiredCodeblock = useWiredData<CodeBlock|null>(React.useCallback(async () => {
    if (props.evidence.contentType !== 'codeblock') return null
    const jsonEvidence = await getJSONEvidence({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidence.uuid,
      previewOrMedia: 'media',
    })
    return {
      type: 'codeblock',
      language: jsonEvidence.contentSubtype,
      code: jsonEvidence.content,
      source: jsonEvidence.metadata ? jsonEvidence.metadata.source : null,
    }
  }, [props.operationSlug, props.evidence.uuid, props.evidence.contentType]))
  return (
    <Modal title="Edit Evidence" onRequestClose={props.onRequestClose}>
      {wiredCodeblock.render(codeBlock => (
        <InternalEditEvidenceModal {...props} codeBlock={codeBlock} />
      ))}
    </Modal>
  )
}

const InternalEditEvidenceModal = (props: {
  evidence: Evidence,
  onEdited: () => void,
  onRequestClose: () => void,
  operationSlug: string,
  codeBlock: CodeBlock|null,
}) => {
*/

export const EditEvidenceModal = (props: {
  evidence: Evidence,
  onEdited: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const descriptionField = useFormField<string>(props.evidence.description)
  const tagsField = useFormField<Array<Tag>>(props.evidence.tags)
  const codeblockField = useFormField<CodeBlock>({ type: 'codeblock', language: '', code: '', source: null })
  React.useEffect(() => {
    if (props.evidence.contentType !== 'codeblock') {
      return
    }
    getEvidenceAsCodeblock({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidence.uuid,
    }).then(codeblockField.onChange)
  }, [props.evidence.contentType, codeblockField.onChange, props.operationSlug, props.evidence.uuid])

  const formComponentProps = useForm({
    fields: [descriptionField, tagsField, codeblockField],
    onSuccess: () => { props.onEdited(); props.onRequestClose() },
    handleSubmit: () => updateEvidence({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidence.uuid,
      description: descriptionField.value,
      oldTags: props.evidence.tags,
      newTags: tagsField.value,
      updatedContent: props.evidence.contentType === 'codeblock' ? codeblockToBlob(codeblockField.value) : null,
    }),
  })
  return (
    <ModalForm title="Edit Evidence" submitText="Save" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <TextArea label="Description" {...descriptionField} />
      {props.evidence.contentType === 'codeblock' && (
        <CodeBlockEditor {...codeblockField} />
      )}
      <TagChooser operationSlug={props.operationSlug} label="Tags" {...tagsField} />
    </ModalForm>
  )
}

export const ChangeFindingsOfEvidenceModal = (props: {
  evidence: Evidence,
  onChanged: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const wiredFindings = useWiredData<Array<Finding>>(React.useCallback(() => getFindingsOfEvidence({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidence.uuid,
  }), [props.operationSlug, props.evidence.uuid]))

  return (
    <Modal title="Select Findings For Evidence" onRequestClose={props.onRequestClose}>
      {wiredFindings.render(initialFindings => (
        <InternalChangeFindingsOfEvidenceModal {...props} initialFindings={initialFindings} />
      ))}
    </Modal>
  )
}

const InternalChangeFindingsOfEvidenceModal = (props: {
  evidence: Evidence,
  onChanged: () => void,
  onRequestClose: () => void,
  operationSlug: string,
  initialFindings: Array<Finding>,
}) => {
  const oldFindingsField = useFormField<Array<Finding>>(props.initialFindings)
  const newFindingsField = useFormField<Array<Finding>>(props.initialFindings)
  const formComponentProps = useForm({
    fields: [newFindingsField],
    onSuccess: () => { props.onChanged(); props.onRequestClose() },
    handleSubmit: () => changeFindingsOfEvidence({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidence.uuid,
      oldFindings: oldFindingsField.value,
      newFindings: newFindingsField.value,
    }),
  })

  return (
    <Form submitText="Update Evidence" cancelText="Cancel" onCancel={props.onRequestClose} {...formComponentProps}>
      <FindingChooser operationSlug={props.operationSlug} {...newFindingsField} />
    </Form>
  )
}

export const DeleteEvidenceModal = (props: {
  evidence: Evidence,
  onDeleted: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const deleteAssociatedFindingsField = useFormField(false)
  const formComponentProps = useForm({
    fields: [deleteAssociatedFindingsField],
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => deleteEvidence({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidence.uuid,
      deleteAssociatedFindings: deleteAssociatedFindingsField.value,
    }),
  })

  return (
    <ModalForm title="Delete Evidence" submitText="Delete Evidence" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <p>Are you sure you want to delete this evidence?</p>
      <Checkbox label="Also delete any findings associated with this evidence" {...deleteAssociatedFindingsField} />
    </ModalForm>
  )
}
