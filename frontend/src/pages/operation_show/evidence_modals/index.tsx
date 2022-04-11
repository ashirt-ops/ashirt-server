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
import BinaryUpload from 'src/components/binary_upload'
import ComboBox from 'src/components/combobox'
import TagList from 'src/components/tag_list'
import { CodeBlockEditor } from 'src/components/code_block'
import { Evidence, Finding, Tag, CodeBlock, SubmittableEvidence, Operation, TagDifference, SupportedEvidenceType } from 'src/global_types'
import { default as Input, TextArea } from 'src/components/input'
import { useForm, useFormField } from 'src/helpers/use_form'
import { codeblockToBlob } from 'src/helpers/codeblock_to_blob'
import { useWiredData } from 'src/helpers'
import {
  createEvidence, updateEvidence, deleteEvidence, changeFindingsOfEvidence,
  getFindingsOfEvidence, getEvidenceAsCodeblock, getOperations, getEvidenceMigrationDifference,
  moveEvidence
} from 'src/services'
import classnames from 'classnames/bind'
import { ExpandableSection } from 'src/components/expandable_area'
import { escapeRegExp } from 'lodash'
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

/**
 * highlightSubstring breaks a given string into words that match the given regex, joined with
 * the rest of the string. This should preserve case.
 * 
 * @example 
 * const result = highlightSubstring("The quick brown fox jumps over the lazy dog.", /the/gi, "highlight")
 * assert( result, [
 *   <span className="highlight">The</span>,
 *   <span> quick brown fox jumps over </span>,
 *   <span className="highlight">the</span>,
 *   <span> lazy dog.</span>,
 * ])
 * 
 * @param s The string with a substring to highlight
 * @param regex What part of the string to match. Must be a global match (/.../g)
 * @param className What class name to apply to the highlighted word
 * @returns An array of spans. Spans will either be plain, or with the given classname.
 */
const highlightSubstring = (s: string, regexAsStr: string, className: any, options?: { regexFlags: string }): Array<React.ReactNode> => {
  const rtn: Array<React.ReactNode> = []
  const matches = [...s.matchAll(new RegExp(escapeRegExp(regexAsStr), "g" + (options?.regexFlags ?? "") ))]

  const endOfWord = (match: RegExpMatchArray) => (match.index ?? 0) + match[0].length
  const highlight = (v: string) => <span className={className}>{v}</span>

  if (matches.length) {
    if ((matches[0].index ?? 0) > 0) {
      rtn.push(<span>{s.substring(0, matches[0].index)}</span>)
    }

    for (let i = 0; i < matches.length; i++) {
      const item = matches[i]
      const next = matches[i + 1]
      const [value] = item
      rtn.push(highlight(value))
      if (next) {
        const end = endOfWord(item)
        const startOfNextWord = next.index ?? end
        if (end != startOfNextWord) {
          rtn.push(<span>{s.substring(end, startOfNextWord)}</span>)
        }
      }
    }
    const lastEntry = (matches[matches.length - 1])
    rtn.push(<span>{s.substring(endOfWord(lastEntry))}</span>)
  }
  else {
    rtn.push(<span>{s}</span>)
  }

  return rtn
}

export const ViewEvidenceMetadataModal = (props: {
  evidence: Evidence,
  onRequestClose: () => void,
}) => {
  const filterField = useFormField<string>("")
  const initiallyExpanded = props.evidence.metadata.length == 1

  const formComponentProps = useForm({
    fields: [filterField],
    onSuccess: () => { props.onRequestClose() },
    handleSubmit: async () => { },
  })

  return (
    <ModalForm title="Evidence Metadata" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <div className={cx('view-metadata-root')}>
        <Input label="filter" {...filterField} />
        {props.evidence.metadata
          .map((meta) => {
            const content = highlightSubstring(meta.body, filterField.value, cx("content-important"), {regexFlags: "i"})

            return (
              <ExpandableSection
                key={meta.source}
                label={meta.source}
                initiallyExpanded={initiallyExpanded}
                labelClassName={cx(
                  (content.length == 1 && filterField.value.length > 0)
                    ? 'label-not-important'
                    : ''
                )}
              >
                <span className={cx('metadata-content')}>{...content}</span>

              </ExpandableSection>
            )
          }
          )}
      </div>
    </ModalForm>
  )
}
