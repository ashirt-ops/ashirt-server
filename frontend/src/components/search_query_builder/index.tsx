// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import * as dateFns from 'date-fns'
import { useFormField, useModal, renderModals, useWiredData } from 'src/helpers'
import { listEvidenceCreators } from 'src/services'
import { SearchOptions } from './helpers'
import { Evidence, Tag, User, ViewName } from 'src/global_types'
import { MaybeDateRange } from 'src/components/date_range_picker/range_picker_helpers'

import Button from 'src/components/button'
import { default as ComboBox, ComboBoxItem } from 'src/components/combobox'
import DateRangePicker from 'src/components/date_range_picker'
import EvidenceChooser from 'src/components/evidence_chooser'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import TagChooser from 'src/components/tag_chooser'
import WithLabel from 'src/components/with_label'

const cx = classnames.bind(require('./stylesheet'))


export default (props: {
  operationSlug: string
  searchOptions: SearchOptions
  viewName: ViewName
  onChanged: (result: SearchOptions) => void,
  buttonName?: string
}) => {
  const [searchOptions, setSearchOptions] = React.useState<SearchOptions>(props.searchOptions)

  const wiredCreators = useWiredData<Array<User>>(
    React.useCallback(() => listEvidenceCreators({ operationSlug: props.operationSlug }), [props.operationSlug])
  )

  const onFormSubmit = () => {
    props.onChanged(searchOptions)
  }

  return (
    <div className={cx('root')}>
      {wiredCreators.render(users => (<>
        <FilterFields
          operationSlug={props.operationSlug}
          viewName={props.viewName}
          searchOptions={searchOptions}
          allCreators={users}
          onChanged={options => setSearchOptions({ ...options })}
        />
        <Button primary className={cx('submit-button')} onClick={onFormSubmit}>{props.buttonName || "Submit"}</Button>
      </>)
      )}
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

  const dateRange = props.searchOptions.dateRange

  const allCreators = props.allCreators.map(user => ({ name: `${user.firstName} ${user.lastName}`, value: user.slug }))

  const descriptionProps = {
    value: props.searchOptions.text,
    onChange: (text: string) => props.onChanged({ ...props.searchOptions, text })
  }
  const tagProps = {
    operationSlug: props.operationSlug,
    value: props.searchOptions.tags || [],
    onChange: (tags: Array<Tag>) => props.onChanged({ ...props.searchOptions, tags }),
  }
  const dateProps = {
    value: dateRange ? `${toEnUSDate(dateRange[0])} to ${toEnUSDate(dateRange[1])}` : '',
    range: dateRange || null,
    onSelectRange: (r: MaybeDateRange) => props.onChanged({ ...props.searchOptions, dateRange: r || undefined })
  }
  const sortProps = {
    options: supportedSortDirections,
    value: props.searchOptions.sortAsc,
    onChange: (sortAsc: boolean) => props.onChanged({ ...props.searchOptions, sortAsc }),
  }
  const creatorProps = {
    options: [{ name: 'Any', value: undefined }, ...allCreators],
    value: props.searchOptions.operator,
    onChange: (creator?: string) => props.onChanged({ ...props.searchOptions, operator: creator }),
  }
  const linkedProps = {
    options: supportedLinking,
    value: props.searchOptions.hasLink,
    onChange: (hasLink?: boolean) => props.onChanged({ ...props.searchOptions, hasLink }),
  }

  const chooseEvidenceModal = useModal<void>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={props.searchOptions.withEvidenceUuid == null
        ? []
        : [uuidToBasicEvidence(props.searchOptions.withEvidenceUuid)]
      }
      operationSlug={props.operationSlug}
      onChanged={uuid => props.onChanged({ ...props.searchOptions, withEvidenceUuid: uuid })}
      {...modalProps}
    />
  ))

  return (
    <div className={cx('root')}>
      <Input label="Description" {...descriptionProps} />
      <TagChooser label="Include Tags" {...tagProps} />

      <SplitInputRow label="Date Range" inputValue={dateProps.value} className={'date-range-input'}>
        <DateRangePicker {...dateProps} />
      </SplitInputRow>

      <ComboBox label="Sort Direction" {...sortProps} />
      <ComboBox label="Creator" {...creatorProps} />

      { props.viewName == 'evidence'
        ? <ComboBox label="Exists in Finding" {...linkedProps} />
        : (
          <SplitInputRow label="Includes Evidence (uuid)" className={'linked-evidence-input'}
            inputValue={props.searchOptions.withEvidenceUuid || ''}>
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
  initialEvidence: [Evidence] | [],
  onRequestClose: () => void,
  onChanged: (uuid: string) => void,
  operationSlug: string,
}) => {
  const evidenceField = useFormField<Array<Evidence>>(props.initialEvidence)

  return (
    <Modal title="Search for evidence" onRequestClose={props.onRequestClose}>
      <EvidenceChooser operationSlug={props.operationSlug} {...evidenceField} />
      <Button primary className={cx('submit-button')} onClick={() => {
        props.onChanged(evidenceField.value.length > 0 ? evidenceField.value[0].uuid : '')
        props.onRequestClose()
      }}>Select</Button>
    </Modal>
  )
}

const supportedSortDirections: Array<ComboBoxItem<boolean>> = [
  { name: "Descending (default)", value: false },
  { name: "Ascending", value: true },
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
