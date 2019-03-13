// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import RadioGroup from 'src/components/radio_group'
import SettingsSection from 'src/components/settings_section'
import {OperationStatus, operationStatusToLabel} from 'src/global_types'
import {getOperation, saveOperation} from 'src/services'
import {useForm, useFormField} from 'src/helpers/use_form'
import {useWiredData} from 'src/helpers'

const EditForm = (props: {
  name: string,
  status: OperationStatus,
  onSave: (op: {name: string, status: OperationStatus}) => Promise<void>,
}) => {
  const nameField = useFormField(props.name)
  const statusField = useFormField(props.status)
  const formComponentProps = useForm({
    fields: [nameField, statusField],
    handleSubmit: () => props.onSave({name: nameField.value, status: statusField.value}),
  })

  return (
    <Form submitText="Save Changes" {...formComponentProps}>
      <Input label="Name" {...nameField} />
      <RadioGroup
        groupLabel="Status"
        getLabel={(s: OperationStatus) => operationStatusToLabel[s]}
        options={[OperationStatus.PLANNING, OperationStatus.ACTIVE, OperationStatus.COMPLETE]}
        {...statusField}
      />
    </Form>
  )
}

export default (props: {
  operationSlug: string,
}) => {
  const wiredOperation = useWiredData(React.useCallback(() => getOperation(props.operationSlug), [props.operationSlug]))

  return (
    <SettingsSection title="Operation Settings">
      {wiredOperation.render(operation => (
        <EditForm
          name={operation.name}
          status={operation.status}
          onSave={({name, status}) => saveOperation(props.operationSlug, {name, status})}
        />
      ))}
    </SettingsSection>
  )
}
