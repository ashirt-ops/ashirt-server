// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'
import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import NewOperationButton from './new_operation_button'
import OperationCard from './operation_card'
import classnames from 'classnames/bind'
import { getOperations, createOperation, hasFlag } from 'src/services'
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

  type Header = "Other" | "Favorites" | null
  const showOperationList = (ops: Operation[], header: Header) => (
    <>
      {header && <h1 className={cx('opTitle')}>
          {header}
      </h1>}
      <div className={cx('operationList')}>
      {
        ops
          .filter(op => normalizedInclude(op.name, filterText.value))
          .map(op => {
            return (
            <OperationCard
              slug={op.slug}
              status={op.status}
              numUsers={op.numUsers}
              key={op.slug}
              name={op.name}
              favorite={op.favorite}
              className={cx('card')}
            />
          )})
      }
      {header !== "Other" && <NewOperationButton onClick={() => newOperationModal.show({})} />}
    </div>
    </>
  )

  const favoriteOps = ops?.filter(op => op.favorite)
  const otherOps = ops?.filter(op => !op.favorite)

  const favOpsExist = favoriteOps?.length > 0
  const bothOpsExist = favOpsExist && otherOps.length > 0

  const returnBothOpTypes = [showOperationList(favoriteOps, "Favorites"), showOperationList(otherOps, "Other")].map(oplist => oplist)
  const returnOneCateogry = favOpsExist ? showOperationList(favoriteOps, null) : showOperationList(otherOps, null)

  const renderBoth = bothOpsExist
    ? returnBothOpTypes
    : returnOneCateogry

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
        {renderBoth}
        {renderModals(newOperationModal)}
        </>)}
    </div>
  )
}

const normalizedInclude = (baseString: string, term: string) => {
  return baseString.toLowerCase().includes(term.toLowerCase())
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
