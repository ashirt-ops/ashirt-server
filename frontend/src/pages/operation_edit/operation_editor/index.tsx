// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import { saveOperation } from 'src/services'
import {useForm, useFormField} from 'src/helpers/use_form'

const EditForm = (props: {
  name: string,
  onSave: (op: {name: string }) => Promise<void>,
}) => {

  const nameField = useFormField(props.name)
  React.useEffect(() => {
    nameField.onChange(props.name)
  }, [props.name, nameField])

  const formComponentProps = useForm({
    fields: [nameField],
    handleSubmit: () => props.onSave({name: nameField.value }),
  })

  return (
    <Form submitText="Save Changes" {...formComponentProps}>
      <Input label="Name" {...nameField} />
    </Form>
  )
}

export default (props: {
  operationSlug: string,
  setCanViewGroups: (canViewGroups: boolean) => void,
  operationName: string,
}) => {
  return (
    <SettingsSection title="Operation Settings">
      <EditForm
        name={props.operationName}
        onSave={({name}) => saveOperation(props.operationSlug, {name})}
      />
    </SettingsSection>
  )
}
