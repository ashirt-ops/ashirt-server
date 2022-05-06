// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import * as dateFns from 'date-fns'

import { Evidence, ViewName } from "src/global_types"
import { useForm, useFormField, useModal, renderModals } from 'src/helpers'

import {
  EvidenceTypeChooser,
  ManagedCreatorChooser as CreatorChooser,
  TagPicker,
} from 'src/components/bullet_chooser'
import { ButtonGroup, default as Button } from 'src/components/button'
import { ComboBoxItem, default as ComboBox } from 'src/components/combobox'
import DateRangePicker from 'src/components/date_range_picker'
import { DateRange } from 'src/components/date_range_picker/range_picker_helpers'
import EvidenceChooser from 'src/components/evidence_chooser'
import ModalForm from 'src/components/modal_form'
import Input from 'src/components/input'
import { SearchOptions, stringifySearch } from "src/components/search_query_builder/helpers"
import WithLabel from 'src/components/with_label'

import Modal from 'src/components/modal'
import { upsertQuery } from 'src/services'

const cx = classnames.bind(require('./stylesheet'))

const renderDateRange = (d: DateRange) => `${toEnUSDate(d[0])} to ${toEnUSDate(d[1])}`

export const FilterFieldsGrid = (props: {
  operationSlug: string
  viewName: ViewName
  value: SearchOptions
  onChange: (result: SearchOptions) => void
  withButtonRow: boolean
  queryName?: string
  cancelText?: string
  onCanceled?: () => void
  submitText?: string
  onSubmit?: (data: SearchOptions) => void
  className?: string
  requestQueriesReload?: () => void
}) => {
  const [state, dispatch] = React.useReducer<SearchOptionsReducer>(searchOptionsReducer, props.value)
  const chooseEvidenceModal = useModal<void>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={
        state.withEvidenceUuid
          ? state.withEvidenceUuid.map(uuidToBasicEvidence)
          : []
      }
      operationSlug={props.operationSlug}
      onChanged={uuid => dispatch({
        type: "set-value",
        field: "withEvidenceUuid",
        newValue: uuid
      })}
      {...modalProps}
    />
  ))
  const saveQueryModal = useModal<void>(modalProps=> (
    <SaveQueryModal
      query={stringifySearch(state)}
      onSaved={() => {
        props.requestQueriesReload?.()
        props.onSubmit?.(state)
      }}
      operationSlug={props.operationSlug}
      view={props.viewName}
      name={props.queryName}
      {...modalProps}
    />
  ))

  React.useEffect(() => {
    dispatch({
      type: 'set-full-state',
      value: props.value
    })
  }, [props.value])

  const mkUpdateState = (field: keyof SearchOptions) =>
    (newValue: SearchOptions[typeof field]) =>
      dispatch({
        type: 'set-value',
        field,
        newValue
      })

  const dateRange = state.dateRange ?? null
  return (
    <div className={cx('grid-container', props.className)}>
      <Cell startOfRow span={2}>
        <Input label="Description" value={state.text} onChange={mkUpdateState('text')} />
      </Cell>

      <Cell startOfRow>
        <TagPicker
          label="Tags"
          operationSlug={props.operationSlug}
          value={state.tags ?? []}
          onChange={mkUpdateState('tags')}
          enableNot
        />
      </Cell>
      <Cell>
        <CreatorChooser
          label='Creators'
          operationSlug={props.operationSlug}
          value={state.operator ?? []}
          onChange={mkUpdateState('operator')}
          enableNot
        />
      </Cell>

      <Cell startOfRow>
        <EvidenceTypeChooser label="Evidence Type"
          value={state.type ?? []}
          onChange={mkUpdateState('type')}
          enableNot
        />
      </Cell>
      <Cell>
        <SplitInputRow label="Date Range" inputValue={dateRange ? renderDateRange(dateRange) : ''}>
          <DateRangePicker range={dateRange} onSelectRange={v => mkUpdateState('dateRange')(v ?? undefined)} />
        </SplitInputRow>
      </Cell>

      {props.viewName == 'findings' && (
        <Cell startOfRow span={2}>
          <SplitInputRow label="Includes Evidence (uuid)" className={'linked-evidence-input'}
            inputValue={state.withEvidenceUuid?.join(', ') ?? ''}>
            <Button doNotSubmit onClick={() => chooseEvidenceModal.show()}>Browse</Button>
          </SplitInputRow>
        </Cell>
      )}

      <Cell startOfRow>
        <ComboBox label="Sort Direction"
          className={cx('grid-col-1')}
          options={[
            { name: "Newest First (default)", value: false },
            { name: 'Oldest First', value: true }
          ]}
          value={state.sortAsc}
          onChange={mkUpdateState('sortAsc')}
        />
      </Cell>
      <Cell>
        {props.viewName == 'evidence' && (
          <ComboBox
            label="Exists in Finding"
            options={supportedLinking}
            value={state.hasLink}
            onChange={mkUpdateState('hasLink')}
          />
        )}
      </Cell>

      {/* Always the last row */}
      {props.withButtonRow && (
        <>
          <ButtonGroup className={cx('button-row')}>
            <Button onClick={() => props.onCanceled?.()}>{props.cancelText ?? "Close"}</Button>
            <Button primary onClick={() => props.onSubmit?.(state)}>{props.cancelText ?? "Apply"}</Button>
          </ButtonGroup>
          <Button className={cx('save-button')} primary onClick={() =>saveQueryModal.show()}>Save</Button>
        </>
      )}
      {renderModals(chooseEvidenceModal, saveQueryModal)}
    </div>
  )
}

const Cell = (props: {
  startOfRow?: true
  span?: 2
  children?: React.ReactNode
}) => {
  const span = props.span ?? 1
  return (
    <div className={cx(
      props.startOfRow ? 'grid-col-1' : null,
      span == 2 ? 'grid-span-2' : null
    )}>
      {props.children}
    </div>
  )
}

const supportedLinking: Array<ComboBoxItem<boolean | undefined>> = [
  { name: "Any", value: undefined },
  { name: "Only Included", value: true },
  { name: "Only Non-included", value: false },
]

type SearchOptionsReducer = (state: SearchOptions, action: SearchOptionAction) => SearchOptions

type SearchOptionAction =
  | SetSearchOption
  | SetNewState

type SetSearchOption = {
  type: 'set-value',
  field: keyof SearchOptions,
  newValue: SearchOptions[keyof SearchOptions]
}

type SetNewState = {
  type: 'set-full-state',
  value: SearchOptions
}

const searchOptionsReducer = (state: SearchOptions, action: SearchOptionAction) => {
  if (action.type == 'set-value') {
    return {
      ...state,
      [action.field]: action.newValue
    }
  }
  if (action.type == 'set-full-state') {
    return action.value
  }
  return state
}

const SplitInputRow = (props: {
  label: string,
  inputValue: string,
  className?: string,
  children: React.ReactNode,
}) => (
  <WithLabel label={props.label}>
    <div className={cx('multi-item-row')}>
      <Input readOnly className={cx('flex-input', props.className)} value={props.inputValue} />
      {props.children}
    </div>
  </WithLabel>
)

const toEnUSDate = (d: Date) => dateFns.format(d, "MMM dd, yyyy")

const uuidToBasicEvidence = (uuid: string): Evidence => ({
  uuid: uuid,
  description: "",
  operator: { slug: "", firstName: "", lastName: "", },
  occurredAt: new Date(),
  tags: [],
  contentType: 'none'
})

const ChooseEvidenceModal = (props: {
  initialEvidence: Array<Evidence>,
  onRequestClose: () => void,
  onChanged: (uuid: Array<string>) => void,
  operationSlug: string,
}) => {
  const evidenceField = useFormField<Array<Evidence>>(props.initialEvidence)

  return (
    <Modal title="Search for evidence" onRequestClose={props.onRequestClose}>
      <EvidenceChooser operationSlug={props.operationSlug} {...evidenceField} />
      <Button primary className={cx('submit-button')} onClick={() => {
        props.onChanged(evidenceField.value.map(evi => evi.uuid))
        props.onRequestClose()
      }}>Select</Button>
    </Modal>
  )
}

// TODO: we probably want to move this, but I'm not sure where we would put this yet.
export const SaveQueryModal = (props: {
  name?: string,
  query: string,
  operationSlug: string,
  view: ViewName,
  onRequestClose: () => void,
  onSaved: (name: string) => void,
}) => {
  const nameField = useFormField(props.name ?? '')
  const queryForm = useForm({
    fields: [nameField],
    onSuccess: () => {
      props.onRequestClose()
      props.onSaved(nameField.value)
    },
    handleSubmit: async () => {
      const replaceName = nameField.value !== props.name
      const { operationSlug, query, view } = props
      upsertQuery({
        operationSlug, query, replaceName,
        type: view,
        name: nameField.value,
      })
    }
  })

  return (
    <ModalForm title="Save Query" submitText="Save" onRequestClose={props.onRequestClose} {...queryForm}>
      <Input label="Name" {...nameField} />
    </ModalForm>
  )
}
