import * as React from 'react'
import { Link } from 'react-router'
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
  const descriptionField = useFormField<string | undefined>(props.tag?.description)

  const formComponentProps = useForm({
    fields: [nameField, colorField, descriptionField],
    onSuccess: () => { props.onEdited(); props.onRequestClose() },
    handleSubmit: () => {
      return (
        props.tag === undefined
          ? props.createFn({
            name: nameField.value.trim(),
            colorName: colorField.value,
            description: descriptionField.value,
          })
          : props.updateFn({
            id: props.tag.id,
            name: nameField.value.trim(),
            colorName: colorField.value,
            description: descriptionField.value,
          })
      )
    }
  })

  return (
    <Modal title={props.tag ? "Edit Tag" : "Create Tag"} onRequestClose={props.onRequestClose}>
      <Form submitText="Save" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <Input label="Name" {...nameField} />
        <Input label="Description" {...descriptionField} />
        <TagColorPicker label="Color" {...colorField} />
      </Form>
    </Modal>
  )
}

export const UpsertOperationTagModal = (props: {
  onEdited: () => void,
  onRequestClose: () => void,
  operationSlug: string,
  tag?: TagType,
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

const DeleteTagModal = (props: {
  tag: TagType,
  onDeleted: () => void,
  onRequestClose: () => void,
  deleteFn: (id: number) => Promise<void>
  children?: React.ReactNode
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onDeleted(); props.onRequestClose() },
    handleSubmit: () => props.deleteFn(props.tag.id)
  })

  return (
    <Modal title="Delete Tag" onRequestClose={props.onRequestClose}>
      <Form submitText="Delete Tag" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>
          Are you sure you want to delete <Tag name={props.tag.name} color={props.tag.colorName} />?
        </p>
        {props.children}
      </Form>
    </Modal>
  )
}

export const DeleteOperationTagModal = (props: {
  tag: TagType,
  onDeleted: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList({
    operationSlug: props.operationSlug,
    query: `tag:${JSON.stringify(props.tag.name)}`,
  }), [props.operationSlug, props.tag.name]))

  return (
    <DeleteTagModal
      {...props}
      deleteFn={(id) => deleteTag({ id, operationSlug: props.operationSlug })}
    >
      {wiredEvidence.render(evidence => (
        <p>
          {evidence.length > 0 && 'This tag belongs to the following evidence and will be removed from them on deletion:'}
          {evidence.map(evi => (
            <Link
              to={`/operations/${props.operationSlug}/evidence/${evi.uuid}`}
              key={evi.uuid}
              children={evi.description.substring(0, 50)}
              style={evidenceLinkStyle}
            />
          ))}
        </p>
      ))}
    </DeleteTagModal>
  )
}

export const DeleteDefaultTagModal = (props: {
  tag: TagType,
  onDeleted: () => void,
  onRequestClose: () => void,
}) => {
  return (
    <DeleteTagModal {...props}
      deleteFn={(id) => deleteDefaultTag({ id })}
    />
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
