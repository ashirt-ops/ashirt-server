import * as React from 'react'
import classnames from 'classnames/bind'
import { useForm, useFormField } from 'src/helpers'
import { SearchOptions, SearchType } from './helpers'
import { Tag, Evidence } from 'src/global_types'
import { addOrUpdateDateRangeInQuery, getDateRangeFromQuery, useModal, renderModals } from 'src/helpers'
import { MaybeDateRange } from 'src/components/date_range_picker/range_picker_helpers'

import DateRangePicker from 'src/components/date_range_picker'
import Input from 'src/components/input'
import TagChooser from 'src/components/tag_chooser'
import { default as ComboBox, ComboBoxItem } from 'src/components/combobox'
import EvidenceChooser from 'src/components/evidence_chooser'
import ModalForm from 'src/components/modal_form'

import Button from '../button'

const cx = classnames.bind(require('./stylesheet'))


export default (props: {
  operationSlug: string
  searchOptions: SearchOptions
  searchType: SearchType
  onChanged: (result: SearchOptions) => void
}) => {
  const initialEvidence: [Evidence] | [] = props.searchOptions.withEvidenceUuid == null
    ? []
    : [uuidToBasicEvidence(props.searchOptions.withEvidenceUuid)]

  const descriptionField = useFormField<string>(props.searchOptions.text)
  const tagsField = useFormField<Array<Tag>>([])
  const [dateRangeStr, setDateRangeStr] = React.useState<string>(
    dateRangeToStr(toMaybeDateRange(props.searchOptions.dateRange))
  )
  const sortingField = useFormField<boolean>(props.searchOptions.sortAsc)
  const linkedField = useFormField<boolean | undefined>(props.searchOptions.hasLink)

  const chooseEvidenceModal = useModal<void>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={initialEvidence}
      operationSlug={props.operationSlug}
      onChanged={() => { }}
      {...modalProps}
    />
  ))

  const onFormSubmit = () => {
    const rtn: SearchOptions = {
      text: descriptionField.value,
      sortAsc: sortingField.value,
    }
    if (props.searchOptions.uuid) { // just forward it along for now
      rtn.uuid = props.searchOptions.uuid
    }
    if (tagsField.value.length > 0) {
      rtn.tags = tagsField.value.map(tag => tag.name)
    }
    if (dateRangeStr.length > 0) {
      //  rtn.dateRange = dateRangeStr
    }
    if (linkedField) {
      rtn.hasLink = linkedField.value
    }
    // operator ?: string, // TODO

    // withEvidenceUuid ?: string,
    props.onChanged(rtn)
  }

  return (
    <div className={cx('root')}>
      <Input label="Description" {...descriptionField} />
      <TagChooser label="Include Tags" operationSlug={props.operationSlug} {...tagsField} />
      <DateRangePicker
        range={getDateRangeFromQuery(dateRangeStr)}
        onSelectRange={r => setDateRangeStr(dateRangeToStr(r))}
      />
      <ComboBox
        label="Exists in Finding"
        options={supportedLinking}
        {...linkedField}
      />
      <ComboBox
        label="Sort Direction"
        options={supportedSortDirections}
        {...sortingField}
      />
      <Input label="Includes Evidence (uuid)" />
      <Button onClick={() => chooseEvidenceModal.show()}>Search for evidence</Button>
      <Button primary onClick={onFormSubmit}>Submit</Button>
      {renderModals(chooseEvidenceModal)}
    </div>
  )
}

const ChooseEvidenceModal = (props: {
  initialEvidence: [Evidence] | [],
  onRequestClose: () => void,
  onChanged: (uuid: string) => void,
  operationSlug: string,
}) => {
  const evidenceField = useFormField<Array<Evidence>>(props.initialEvidence)
  const formComponentProps = useForm({
    fields: [evidenceField],
    onSuccess: () => { props.onChanged(evidenceField.value.length > 0 ? evidenceField.value[0].uuid : ''); props.onRequestClose() },
    handleSubmit: () => {
      return Promise.resolve()
    },
  })

  return (
    <ModalForm title="Search for evidence" submitText="Select" onRequestClose={props.onRequestClose} {...formComponentProps}>
      <EvidenceChooser operationSlug={props.operationSlug} {...evidenceField} />
    </ModalForm>
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


const dateRangeToStr = (r: MaybeDateRange) => addOrUpdateDateRangeInQuery('', r)
const toMaybeDateRange = (range: [Date, Date] | undefined) => range ? range : null

const uuidToBasicEvidence = (uuid: string): Evidence => ({
  uuid: uuid,
  description: "",
  operator: { slug: "", firstName: "", lastName: "", },
  occurredAt: new Date(),
  tags: [],
  contentType: 'none'
})
