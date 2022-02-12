// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { Link } from 'react-router-dom'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import Tag from 'src/components/tag'
import TagColorPicker from 'src/components/tag_color_picker'
import { Tag as TagType } from 'src/global_types'
import { deleteTag, updateTag, getEvidenceList, deleteDefaultTag, updateDefaultTag, createDefaultTag, createTag } from 'src/services'
import { randomTagColorName, useForm, useFormField, useWiredData } from 'src/helpers'

function UpsertTagModal(props: {
  onEdited: () => void,
  onRequestClose: () => void,
  createFn: (tag: Omit<TagType, "id">) => Promise<void>
  updateFn: (tag: TagType) => Promise<void>
  tag?: TagType,
}) {
  const nameField = useFormField<string>(props.tag?.name ?? "")
  const colorField = useFormField<string>(props.tag?.colorName ?? randomTagColorName())

  const formComponentProps = useForm({
    fields: [nameField, colorField],
    onSuccess: () => { props.onEdited(); props.onRequestClose() },
    handleSubmit: () => {
      return (
        props.tag === undefined
          ? props.createFn({
            name: nameField.value.trim(),
            colorName: colorField.value,
          })
          : props.updateFn({
            id: props.tag.id,
            name: nameField.value.trim(),
            colorName: colorField.value,

          })
      )
    }
  })

  return (
    <Modal title={props.tag ? "Edit Tag" : "Create Tag"} onRequestClose={props.onRequestClose}>
      <Form submitText="Save" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <Input label="Name" {...nameField} />
        <TagColorPicker label="Color" {...colorField} />
      </Form>
    </Modal>
  )
}

export const EditTagModal = (props: {
  onEdited: () => void,
  onRequestClose: () => void,
  operationSlug: string,
  tag: TagType,
}) => {
  return (
    <UpsertTagModal
      {...props}
      createFn={async (t) => { createTag({ ...t, operationSlug: props.operationSlug }) }}
      updateFn={async (t) => updateTag({ ...t, operationSlug: props.operationSlug })}
    />
  )
}

export const UpsertDefaultTagModal = (props: {
  onEdited: () => void,
  onRequestClose: () => void,
  tag?: TagType,
}) => {
  return (
    <UpsertTagModal
      {...props}
      createFn={async (t) => { await createDefaultTag(t) }}
      updateFn={updateDefaultTag}
    />
  )
}

export const DeleteTagModal = (props: {
  tag: TagType,
  onDeleted: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => deleteTag({
      id: props.tag.id,
      operationSlug: props.operationSlug,
    }),
  })

  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList({
    operationSlug: props.operationSlug,
    query: `tag:${JSON.stringify(props.tag.name)}`,
  }), [props.operationSlug, props.tag.name]))

  return (
    <Modal title="Delete Tag" onRequestClose={props.onRequestClose}>
      <Form submitText="Delete Tag" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>
          Are you sure you want to delete <Tag name={props.tag.name} color={props.tag.colorName} />?
        </p>
        {wiredEvidence.render(evidence => (
          <p>
            {evidence.length > 0 && 'This tag belongs to the following evidence and will be removed from them on deletion:'}
            {evidence.map(evi => (
              <Link
                to={`/operations/${props.operationSlug}/evidence/${evi.uuid}`}
                key={evi.uuid}
                children={evi.description.substr(0, 50)}
                style={evidenceLinkStyle}
              />
            ))}
          </p>
        ))}
      </Form>
    </Modal>
  )
}

export const DeleteDefaultTagModal = (props: {
  tag: TagType,
  onDeleted: () => void,
  onRequestClose: () => void,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => deleteDefaultTag({
      id: props.tag.id,
    }),
  })

  return (
    <Modal title="Delete Default Tag" onRequestClose={props.onRequestClose}>
      <Form submitText="Delete Tag" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>
          Are you sure you want to delete <Tag name={props.tag.name} color={props.tag.colorName} />?
        </p>
      </Form>
    </Modal>
  )
}

const evidenceLinkStyle: React.CSSProperties = {
  display: 'block',
  fontWeight: 800,
  textDecoration: 'underline',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
}
