// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import NewOperationButton from './new_operation_button'
import OperationCard from './operation_card'
import classnames from 'classnames/bind'
import {Operation} from 'src/global_types'
import {RouteComponentProps} from 'react-router-dom'
import {getOperations, createOperation} from 'src/services'
import {useForm, useFormField} from 'src/helpers/use_form'
import {useWiredData, useModal, renderModals} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{}>) => {
  const wiredOperations = useWiredData<Array<Operation>>(getOperations)

  const newOperationModal = useModal<void>(modalProps => (
    <NewOperationModal {...modalProps} onCreated={wiredOperations.reload} />
  ))

  return (
    <div className={cx('root')}>
      {wiredOperations.render(ops => <>
        {ops.map(op => (
          <OperationCard
            slug={op.slug}
            status={op.status}
            numUsers={op.numUsers}
            key={op.slug}
            name={op.name}
            className={cx('card')}
          />
        ))}
        <NewOperationButton onClick={() => newOperationModal.show()} />
      </>)}
      {renderModals(newOperationModal)}
    </div>
  )
}

const NewOperationModal = (props: {
  onRequestClose: () => void,
  onCreated: () => void,
}) => {
  const nameField = useFormField('')
  const formComponentProps = useForm({
    fields: [nameField],
    handleSubmit: () => createOperation(nameField.value),
    onSuccess: () => { props.onCreated(); props.onRequestClose() },
  })

  return (
    <Modal title="New Operation" onRequestClose={props.onRequestClose}>
      <Form submitText="Create Operation" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <Input label="Operation Name" {...nameField} />
      </Form>
    </Modal>
  )
}
