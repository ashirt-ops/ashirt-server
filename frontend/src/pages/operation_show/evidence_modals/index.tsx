// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import {
  Evidence,
  EvidenceMetadata,
  Finding,
  Tag,
  CodeBlock,
  SubmittableEvidence,
  Operation,
  TagDifference,
  SupportedEvidenceType,
} from 'src/global_types'
import {
  codeblockToBlob,
  highlightSubstring,
  useForm,
  useFormField,
  useWiredData,
} from 'src/helpers'
import {
  createEvidence,
  createEvidenceMetadata,
  updateEvidence,
  deleteEvidence,
  changeFindingsOfEvidence,
  getFindingsOfEvidence,
  getEvidenceAsCodeblock,
  getOperations,
  getEvidenceMigrationDifference,
  moveEvidence,
  updateEvidenceMetadata,
} from 'src/services'

import BinaryUpload from 'src/components/binary_upload'
import { default as Button, ButtonGroup } from 'src/components/button'
import Checkbox from 'src/components/checkbox'
import { CodeBlockEditor } from 'src/components/code_block'
import ComboBox from 'src/components/combobox'
import { ExpandableSection } from 'src/components/expandable_area'
import FindingChooser from 'src/components/finding_chooser'
import Form from 'src/components/form'
import ImageUpload from 'src/components/image_upload'
import { default as Input, TextArea } from 'src/components/input'
import ModalForm from 'src/components/modal_form'
import Modal from 'src/components/modal'
import TagChooser from 'src/components/bullet_chooser/tag_chooser'
import TabMenu from 'src/components/tabs'
import TagList from 'src/components/tag_list'


const cx = classnames.bind(require('./stylesheet'))

export const CreateEvidenceModal = (props: {
  onCreated: () => void,
  onRequestClose: () => void,
  operationSlug: string,
}) => {
  const descriptionField = useFormField<string>("")
  const tagsField = useFormField<Array<Tag>>([])
  const binaryBlobField = useFormField<File | null>(null)
  const codeblockField = useFormField<CodeBlock>({ type: 'codeblock', language: '', code: '', source: null })

  const isATerminalRecording = (file: File) => file.type == ''
  const isAnHttpRequestCycle = (file: File) => file.name.endsWith("har")

  const evidenceTypeOptions: Array<{ name: string, value: SupportedEvidenceType, content?: React.ReactNode }> = [
    { name: 'Screenshot', value: 'image', content: <ImageUpload label='Screenshot' {...binaryBlobField} /> },
    { name: 'Code Block', value: 'codeblock', content: <CodeBlockEditor {...codeblockField} /> },
    { name: 'Event', value: 'event', content: <div /> },
    {
      name: 'Terminal Recording', value: 'terminal-recording',
      content: <BinaryUpload label='Terminal Recording' isSupportedFile={isATerminalRecording} {...binaryBlobField} />
    },
    {
      name: 'HTTP Request/Response', value: 'http-request-cycle',
      content: <BinaryUpload label='HAR File' isSupportedFile={isAnHttpRequestCycle} {...binaryBlobField} />
    },
  ]

  const [selectedCBValue, setSelectedCBValue] = React.useState<string>(evidenceTypeOptions[0].value)
  const getSelectedOption = () => evidenceTypeOptions.filter(opt => opt.value === selectedCBValue)[0]

  const formComponentProps = useForm({
    fields: [descriptionField, binaryBlobField],
    onSuccess: () => { props.onCreated(); props.onRequestClose() },
    handleSubmit: () => {
      let data: SubmittableEvidence = { type: "none" }
      const selectedOption = getSelectedOption()
      const fileBasedKeys = ['image', 'terminal-recording', 'http-request-cycle']

      if (selectedOption.value === 'codeblock' && codeblockField.value !== null) {
        data = { type: 'codeblock', file: codeblockToBlob(codeblockField.value) }
      } else if (fileBasedKeys.includes(selectedOption.value) && binaryBlobField.value != null) {
        data = { type: selectedOption.value, file: binaryBlobField.value }
      } else if (selectedOption.value === 'event') {
        data = { type: 'event' }
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
      <ComboBox
        label="Evidence Type"
        className={cx('dropdown')}
        options={evidenceTypeOptions}
        value={selectedCBValue}
        onChange={setSelectedCBValue}
      />
      {getSelectedOption().content}
      <TagChooser operationSlug={props.operationSlug} label="Tags" {...tagsField} />
    </ModalForm>
  )
}

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

export const MoveEvidenceModal = (props: {
  evidence: Evidence,
  operationSlug: string,
  onRequestClose: () => void,
  onEvidenceMoved: () => void,
}) => {

  const [selectedOperationSlug, setSelectedOperation] = React.useState(props.operationSlug)

  const wiredOps = useWiredData<Array<Operation>>(React.useCallback(getOperations, [props.operationSlug, props.evidence.uuid]))
  const wiredDiff = useWiredData<TagDifference>(React.useCallback(() =>
    getEvidenceMigrationDifference({
      fromOperationSlug: props.operationSlug,
      toOperationSlug: selectedOperationSlug,
      evidenceUuid: props.evidence.uuid,
    }), [selectedOperationSlug, props.evidence.uuid, props.operationSlug]))

  const formComponentProps = useForm({
    fields: [],
    onSuccess: () => { props.onEvidenceMoved(); props.onRequestClose() },
    handleSubmit: () => {
      if (selectedOperationSlug == props.operationSlug) {
        return Promise.resolve() // no need to do anything if the to and from destinations are the same
      }
      return moveEvidence({
        fromOperationSlug: props.operationSlug,
        toOperationSlug: selectedOperationSlug,
        evidenceUuid: props.evidence.uuid
      }).then(() => { window.location.href = `/operations/${props.operationSlug}/evidence` })
    },
  })

  return (
    <ModalForm title="Move Evidence To Another Operation" submitText="Move" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <div>
        Moving evidence will disconnect this evidence from any findings and some tags may be
        lost in the transition.
      </div>
      {wiredOps.render(operations => {
        operations.sort((a, b) => a.name.localeCompare(b.name))

        const mappedOperations = operations.map(op => ({ name: op.name, value: op }))
        return (
          <ComboBox
            label="Select a destination operation"
            options={mappedOperations}
            value={operations.filter(op => op.slug === selectedOperationSlug)[0]}
            onChange={op => setSelectedOperation(op.slug)} />
        )
      })}
      {wiredDiff.render(data => (
        <TagListRenderer sourceSlug={props.operationSlug} destSlug={selectedOperationSlug} tags={data.excluded} />
      ))}
    </ModalForm>
  )
}

const TagListRenderer = (props: {
  sourceSlug: string,
  destSlug: string
  tags: Array<Tag> | null
}) => {
  if (props.sourceSlug == props.destSlug) {
    return <div>This is the current operation, and so no changes will be made</div>
  }
  else if (props.tags == null || props.tags.length == 0) {
    return <div>All tags will carry over</div>
  }

  return (<>
    <div>The following tags will be removed:</div>
    <TagList tags={props.tags} />
  </>)
}

export const EvidenceMetadataModal = (props: {
  operationSlug: string,
  evidence: Evidence,
  onRequestClose: () => void,
  onUpdated: () => void,
}) => {

  // this switch controls whether the user is only able to view metadata, or if they can manually
  // create and edit existing metadata. We might get more precise control later.
  const viewOnly = false

  return (
    <Modal title='Evidence Metadata' onRequestClose={props.onRequestClose}>
      {viewOnly
        ? (
          <ViewEditEvidenceMetadataContainer
            evidence={props.evidence}
            operationSlug={props.operationSlug}
          />
        )
        : (
          <TabMenu className={cx('tab-menu')}
            tabs={[
              {
                id: 'view', label: 'View',
                content: (
                  <ViewEditEvidenceMetadataContainer
                    evidence={props.evidence}
                    operationSlug={props.operationSlug}
                    onEdited={() => { props.onUpdated(); props.onRequestClose() }}
                  />
                )
              },
              {
                id: 'create', label: 'Create',
                content: (
                  <AddEvidenceMetadataForm
                    evidence={props.evidence}
                    onCreated={() => { props.onUpdated(); props.onRequestClose() }}
                    operationSlug={props.operationSlug}
                  />
                )
              },
            ]}
          />
        )
      }
    </Modal>
  )
}

const ViewEditEvidenceMetadataContainer = (props: {
  evidence: Evidence,
  operationSlug: string,
  onEdited?: () => void
  onCancel?: () => void,
}) => {
  const [editedMetadata, setEditedMetadata] = React.useState<null | EvidenceMetadata>(null)
  const [filterText, setFilterText] = React.useState<string>("")

  return (
    editedMetadata
      ? (
        <EditEvidenceMetadataForm
          evidence={props.evidence}
          metadata={editedMetadata}
          onCancel={() => setEditedMetadata(null)}
          onEdited={() => {
            props.onEdited?.()
            setEditedMetadata(null)
          }}
          operationSlug={props.operationSlug}
        />
      )
      : (
        <ViewEvidenceMetadataForm
          evidence={props.evidence}
          onMetadataEdited={props.onEdited ? setEditedMetadata : undefined}
          filterText={filterText}
          onFilterUpdated={setFilterText}
        />
      )
  )
}

const EditEvidenceMetadataForm = (props: {
  metadata: EvidenceMetadata
  operationSlug: string
  evidence: Evidence
  onEdited: () => void
  onCancel: () => void
}) => (
  <EvidenceMetadataEditorForm
    metadata={props.metadata}
    readonlySource
    submitText="Save"
    onSubmit={(metadata: EvidenceMetadata) => {
      return updateEvidenceMetadata({
        operationSlug: props.operationSlug,
        evidenceUuid: props.evidence.uuid,
        body: metadata.body,
        source: metadata.source,
      })
    }}
    onEdited={props.onEdited}
    onCancel={props.onCancel}
  />
)

const AddEvidenceMetadataForm = (props: {
  operationSlug: string,
  evidence: Evidence,
  onCreated: () => void,
  onCancel?: () => void,
}) => (
  <EvidenceMetadataEditorForm
    metadata={{ body: "", source: "" }}
    submitText="Create"
    onSubmit={(metadata: EvidenceMetadata) => {
      return createEvidenceMetadata({
        operationSlug: props.operationSlug,
        evidenceUuid: props.evidence.uuid,
        body: metadata.body,
        source: metadata.source,
      })
    }}
    onEdited={props.onCreated}
  />
)

const ViewEvidenceMetadataForm = (props: {
  evidence: Evidence,
  onMetadataEdited?: (metadata: EvidenceMetadata) => void
  onCancel?: () => void,
  filterText: string,
  onFilterUpdated: (val: string) => void
}) => {
  return (
    <div className={cx('view-metadata-root')}>
        {props.evidence.metadata.length == 0
          ? <em>No metadata exists for this evidence</em>
          : (<>
            <Input label="Filter Metadata" value={props.filterText} onChange={props.onFilterUpdated} />

            {props.evidence.metadata
              .map((meta) => {
                return (
                  <EvidenceMetadataItem
                    key={meta.source}
                    meta={meta}
                    filterText={props.filterText}
                    onMetadataEdited={props.onMetadataEdited}
                    expanded={props.evidence.metadata.length == 1}
                  />
                )
              }
              )}
          </>)
        }
    </div>
  )
}

const EvidenceMetadataEditorForm = (props: {
  metadata: EvidenceMetadata
  onSubmit: (metadata: EvidenceMetadata) => Promise<void>
  onEdited: () => void
  onCancel?: () => void
  submitText: string
  readonlySource?: boolean
}) => {
  const sourceField = useFormField<string>(props.metadata.source)
  const contentField = useFormField<string>(props.metadata.body)

  const formComponentProps = useForm({
    fields: [sourceField, contentField],
    onSuccess: () => props.onEdited(),
    handleSubmit: () => {
      if (sourceField.value.trim() == "") {
        throw new Error("Must specify a source")
      }
      return props.onSubmit({
        source: sourceField.value,
        body: contentField.value,
      })
    },
  })

  return (
    <Form {...formComponentProps}
      submitText={props.submitText}
      onCancel={props.onCancel}
      cancelText="Cancel"
    >
      <Input label={'Source' + (props.readonlySource ? ' (locked)' : '')} readOnly={props.readonlySource} {...sourceField} />
      <TextArea label="Content" {...contentField} />
    </Form>
  )
}

const EvidenceMetadataItem = (props: {
  meta: EvidenceMetadata
  filterText: string
  onMetadataEdited?: (metadata: EvidenceMetadata) => void
  expanded?: boolean
}) => {
  const { onMetadataEdited, meta, filterText, expanded } = props
  const minLength = 3
  const content = highlightSubstring(
    meta.body, filterText, cx("content-important"), { regexFlags: "i", minLength }
  )

  const actions: Array<ExpandableSectionLabelActionItem> = onMetadataEdited
    ? [{
      label: 'Edit',
      action: (e) => {
        e.stopPropagation()
        onMetadataEdited(meta)
      }
    }]
    : []

  return (
    <ExpandableSection
      label={(
        <ExpandableSectionLabel label={meta.source} actions={actions} />
      )}
      initiallyExpanded={expanded}
      labelClassName={cx(
        (content.length == 1 && filterText.length >= minLength)
          ? 'label-not-important'
          : ''
      )}
    >
      <span className={cx('metadata-content')}>{...content}</span>
    </ExpandableSection>
  )
}

type ExpandableSectionLabelActionItem = { label: string, action: (e: React.MouseEvent<Element, MouseEvent>) => void }

const ExpandableSectionLabel = (props: {
  label: string
  actions: Array<ExpandableSectionLabelActionItem>
}) => {
  return (
    <div className={cx('expandable-section-label-container')}>
      <span className={cx('expandable-section-label')}>{props.label}</span>
      {props.actions.length > 0 && (
        <ButtonGroup className={cx('expandable-section-button-group')}>
          {props.actions.map(act => (
            <Button small key={act.label} onClick={act.action}>{act.label}</Button>
          ))}
        </ButtonGroup>
      )}
    </div>
  )
}
