// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { OperationStatus, operationStatusToLabel } from 'src/global_types'
import { useDataSource, getOperation, saveOperation } from 'src/services'
import { useWiredData, useForm, useFormField } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import RadioGroup from 'src/components/radio_group'
import SettingsSection from 'src/components/settings_section'

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
  const ds = useDataSource()
  const wiredOperation = useWiredData(React.useCallback(() => (
    getOperation(ds, props.operationSlug)
  ), [ds, props.operationSlug]))

  return (
    <SettingsSection title="Operation Settings">
      {wiredOperation.render(operation => (
        <EditForm
          name={operation.name}
          status={operation.status}
          onSave={({name, status}) => saveOperation(ds, props.operationSlug, {name, status})}
        />
      ))}
    </SettingsSection>
  )
}
