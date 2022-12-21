// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import {getOperation, saveOperation} from 'src/services'
import {useForm, useFormField} from 'src/helpers/use_form'
import {useWiredData} from 'src/helpers'

const EditForm = (props: {
  name: string,
  onSave: (op: {name: string }) => Promise<void>,
}) => {
  const nameField = useFormField(props.name)
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
}) => {
  const wiredOperation = useWiredData(React.useCallback(() => getOperation(props.operationSlug), [props.operationSlug]))

  wiredOperation.expose(operation => props.setCanViewGroups(!!operation?.userCanViewGroups))
  return (
    <SettingsSection title="Operation Settings">
      {wiredOperation.render(operation => {
        props.setCanViewGroups(!!operation?.userCanViewGroups)
        return (
          <EditForm
            name={operation.name}
            onSave={({name}) => saveOperation(props.operationSlug, {name})}
          />
        )
      })}
    </SettingsSection>
  )
}
