// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import * as dateFns from 'date-fns'
import { useFormField, useModal, renderModals } from 'src/helpers'
import { SearchOptions } from './helpers'
import { Evidence, User, ViewName } from 'src/global_types'
import { MaybeDateRange } from 'src/components/date_range_picker/range_picker_helpers'

import Button from 'src/components/button'
import { default as ComboBox, ComboBoxItem } from 'src/components/combobox'
import DateRangePicker from 'src/components/date_range_picker'
import EvidenceChooser from 'src/components/evidence_chooser'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import WithLabel from 'src/components/with_label'
import {
  ManagedCreatorChooser as CreatorChooser,
  EvidenceTypeChooser,
  TagPicker,
} from 'src/components/bullet_chooser'

const cx = classnames.bind(require('./stylesheet'))


export default (props: {
  operationSlug: string
  searchOptions: SearchOptions
  viewName: ViewName
  onChanged: (result: SearchOptions) => void,
  buttonName?: string
}) => {
  const [searchOptions, setSearchOptions] = React.useState<SearchOptions>(props.searchOptions)

  const onFormSubmit = () => {
    props.onChanged(searchOptions)
  }

  return (
    <div className={cx('root')}>
      <FilterFields
        operationSlug={props.operationSlug}
        viewName={props.viewName}
        searchOptions={searchOptions}
        allCreators={[]}
        onChanged={options => setSearchOptions({ ...options })}
      />
      <Button primary className={cx('submit-button')} onClick={onFormSubmit}>{props.buttonName || "Submit"}</Button>
    </div>
  )
}

export const FilterFields = (props: {
  operationSlug: string
  viewName: ViewName
  searchOptions: SearchOptions
  onChanged: (result: SearchOptions) => void
  allCreators: Array<User>
}) => {
  const onChange = (part: Partial<SearchOptions>) => props.onChanged({
    ...props.searchOptions, ...part
  })

  const dateRange = props.searchOptions.dateRange
  const dateProps = {
    value: dateRange ? `${toEnUSDate(dateRange[0])} to ${toEnUSDate(dateRange[1])}` : '',
    range: dateRange || null,
    onSelectRange: (r: MaybeDateRange) => props.onChanged({ ...props.searchOptions, dateRange: r || undefined })
  }
  const linkedProps = {
    options: supportedLinking,
    value: props.searchOptions.hasLink,
    onChange: (hasLink?: boolean) => onChange({ hasLink }),
  }

  const chooseEvidenceModal = useModal<void>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={
        props.searchOptions.withEvidenceUuid
          ? props.searchOptions.withEvidenceUuid.map(uuidToBasicEvidence)
        : []
      }
      operationSlug={props.operationSlug}
      onChanged={uuid => props.onChanged({ ...props.searchOptions, withEvidenceUuid: uuid })}
      {...modalProps}
    />
  ))

  return (
    <div className={cx('root')}>
      <Input label="Description" value={props.searchOptions.text} onChange={text => onChange({ text })} />
      <TagPicker label="Include Tags"
        operationSlug={props.operationSlug}
        value={props.searchOptions.tags ?? []}
        onChange={tags => onChange({ tags })}
        enableNot
      />

      <SplitInputRow label="Date Range" inputValue={dateProps.value} className={'date-range-input'}>
        <DateRangePicker {...dateProps} />
      </SplitInputRow>

      <ComboBox label="Sort Direction"
        options={supportedSortDirections}
        value={props.searchOptions.sortAsc}
        onChange={(sortAsc) => onChange({ sortAsc })}
      />
      <CreatorChooser label='Creators'
        operationSlug={props.operationSlug}
        value={props.searchOptions.operator ?? []}
        onChange={operator => onChange({ operator })}
        enableNot
      />

      {props.viewName == 'evidence'
        ? (
          <>
            <EvidenceTypeChooser label="Evidence Type"
              value={props.searchOptions.type ?? []}
              onChange={type => onChange({ type })}
              enableNot
            />
            <ComboBox label="Exists in Finding" {...linkedProps} />
          </>
        )
        : (
          <SplitInputRow label="Includes Evidence (uuid)" className={'linked-evidence-input'}
            inputValue={props.searchOptions.withEvidenceUuid?.join(', ') ?? ''}>
            <Button doNotSubmit onClick={() => chooseEvidenceModal.show()}>Browse</Button>
          </SplitInputRow>
        )
      }
      {renderModals(chooseEvidenceModal)}
    </div>
  )
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

const supportedSortDirections: Array<ComboBoxItem<boolean>> = [
  { name: "Newest First (default)", value: false },
  { name: "Oldest First", value: true },
]

const supportedLinking: Array<ComboBoxItem<boolean | undefined>> = [
  { name: "Any", value: undefined },
  { name: "Only Included", value: true },
  { name: "Only Non-included", value: false },
]

const toEnUSDate = (d: Date) => dateFns.format(d, "MMM dd, yyyy")

const uuidToBasicEvidence = (uuid: string): Evidence => ({
  uuid: uuid,
  description: "",
  operator: { slug: "", firstName: "", lastName: "", },
  occurredAt: new Date(),
  tags: [],
  contentType: 'none'
})
