// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import classnames from 'classnames/bind'
import List from './list'
import { getOperations, createOperation, hasFlag, setFavorite } from 'src/services'
import { useForm, useFormField } from 'src/helpers/use_form'
import { useWiredData, useModal, renderModals } from 'src/helpers'
import { Operation } from 'src/global_types'
const cx = classnames.bind(require('./stylesheet'))

export default () => {
  const { user } = React.useContext(AuthContext) // user should never be null
  const wiredData = useWiredData(React.useCallback(() => Promise.all([
    getOperations(),
    hasFlag("welcome-message")
  ]), []))

  const newOperationModal = useModal<{}>(modalProps => (
    <NewOperationModal {...modalProps} onCreated={wiredData.reload} />
  ))
  const filterText = useFormField<string>('')

  const [ops, setOps] = React.useState<Operation[]>([])
  const [welcomeFlag, setWelcomeFlag] = React.useState<boolean>(false)

  React.useEffect(() => {
    wiredData.expose(data => {
      if (data) {
        const [ops, welcomeFlag] = data
        setOps(ops)
        setWelcomeFlag(welcomeFlag)
      }
    })
  }, [wiredData])

  return (
    <div className={cx('root')}>
      {wiredData.render(() => <>
        {welcomeFlag && (
          <h1 className={cx('welcomeMessage')}>
            Welcome Back, {user ? `${user.firstName} ${user.lastName}` : "Kotter"}!
          </h1>
        )}
        <Input
          placeholder="Filter Operations"
          className={cx('filterInput')}
          icon={require('./search.svg')}
          {...filterText}
        />
        <List
          ops={ops}
          newOperationModal={newOperationModal}
          filterText={filterText}
          onFavoriteToggled={async (slug, isFav) => {
            await setFavorite(slug, isFav)
            wiredData.reload()
          }}
        />
        {renderModals(newOperationModal)}
      </>)}
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
