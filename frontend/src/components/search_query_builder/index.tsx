import * as React from 'react'
import classnames from 'classnames/bind'
import { useForm, useFormField } from 'src/helpers'
import { SearchOptions, SearchType } from './helpers'
import { Tag, Evidence, User } from 'src/global_types'
import { addOrUpdateDateRangeInQuery, getDateRangeFromQuery, useModal, renderModals } from 'src/helpers'
import { MaybeDateRange } from 'src/components/date_range_picker/range_picker_helpers'
import { useWiredData } from 'src/helpers'
import { listEvidenceCreators } from 'src/services'

import * as dateFns from 'date-fns'

import DateRangePicker from 'src/components/date_range_picker'
import Input from 'src/components/input'
import TagChooser from 'src/components/tag_chooser'
import { default as ComboBox, ComboBoxItem } from 'src/components/combobox'
import EvidenceChooser from 'src/components/evidence_chooser'
import ModalForm from 'src/components/modal_form'
import WithLabel from 'src/components/with_label'

import Button from '../button'

const cx = classnames.bind(require('./stylesheet'))


export default (props: {
  operationSlug: string
  searchOptions: SearchOptions
  searchType: SearchType
  onChanged: (result: SearchOptions) => void
}) => {
  const descriptionField = useFormField<string>(props.searchOptions.text)
  const tagsField = useFormField<Array<Tag>>([])
  const [dateRange, setDateRange] = React.useState<MaybeDateRange>(
    props.searchOptions.dateRange || null
  )
  const sortingField = useFormField<boolean>(props.searchOptions.sortAsc)
  const linkedField = useFormField<boolean | undefined>(props.searchOptions.hasLink)
  const creatorField = useFormField<string | undefined>(props.searchOptions.operator)
  const [evidenceUuid, setEvidenceUuid] = React.useState<string | null>(props.searchOptions.withEvidenceUuid || null)

  const wiredCreators = useWiredData<Array<User>>(
    React.useCallback(() => listEvidenceCreators({ operationSlug: props.operationSlug }), [props.operationSlug])
  )

  const chooseEvidenceModal = useModal<void>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={evidenceUuid == null ? [] : [uuidToBasicEvidence(evidenceUuid)]}
      operationSlug={props.operationSlug}
      onChanged={setEvidenceUuid}
      {...modalProps}
    />
  ))

  const onFormSubmit = () => {
    props.onChanged({
      uuid: props.searchOptions.uuid, // forward along the value
      text: descriptionField.value,
      tags: tagsField.value.map(tag => tag.name),
      operator: creatorField.value,
      dateRange: dateRange || undefined,
      hasLink: linkedField.value,
      sortAsc: sortingField.value,
      withEvidenceUuid: evidenceUuid || undefined,
    })
  }

  return (
    <div className={cx('root')}>
      {wiredCreators.render(users => {
        const creators = users.map(user => ({ name: `${user.firstName} ${user.lastName}`, value: user.slug }))
        return (<>
          <Input label="Description" {...descriptionField} />
          <TagChooser label="Include Tags" operationSlug={props.operationSlug} {...tagsField} />

          <WithLabel label="Date Range">
            <div className={cx('multi-item-row')}>
              <Input className={cx('flex-input', 'date-range-input')} readOnly
                value={dateRange ? `${toEnUSDate(dateRange[0])} to ${toEnUSDate(dateRange[1])}` : ''} />
              <DateRangePicker
                range={dateRange}
                onSelectRange={r => setDateRange(r)}
              />
            </div>
          </WithLabel>

          <ComboBox label="Sort Direction" options={supportedSortDirections} {...sortingField} />
          <ComboBox label="Exists in Finding" options={supportedLinking} {...linkedField} />
          <ComboBox label="Creator" options={[{ name: 'Any', value: undefined }, ...creators]} {...creatorField} />

          <WithLabel label="Includes Evidence (uuid)">
            <div className={cx('multi-item-row')}>
              <Input className={cx('flex-input', 'linked-evidence-input')} readOnly
                value={evidenceUuid || ''} />
              <Button onClick={() => chooseEvidenceModal.show()}>Browse</Button>
            </div>
          </WithLabel>

          <Button primary className={cx('submit-button')} onClick={onFormSubmit}>Submit</Button>
        </>)
      })}
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

const toEnUSDate = (d: Date) => dateFns.format(d, "MMM dd, yyyy")

const uuidToBasicEvidence = (uuid: string): Evidence => ({
  uuid: uuid,
  description: "",
  operator: { slug: "", firstName: "", lastName: "", },
  occurredAt: new Date(),
  tags: [],
  contentType: 'none'
})
