// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { ApiKey } from 'src/global_types'
import Button from 'src/components/button'
import Form from 'src/components/form'
import Modal from "src/components/modal"
import { InputWithCopyButton } from 'src/components/text_copiers'
import { useForm } from 'src/helpers'
import { deleteApiKey } from 'src/services'

const cx = classnames.bind(require('./stylesheet'))

export const NewApiKeyModal = (props: {
  apiKey: ApiKey,
  onRequestClose: () => void
}) => {
  const { apiKey, onRequestClose } = props
  return (
    <Modal title="New API Key" onRequestClose={onRequestClose}>
      <NewApiKeyModalContents apiKey={apiKey}>
        <Button primary onClick={() => onRequestClose()}>Close</Button>
      </NewApiKeyModalContents>
    </Modal>
  )
}

export const NewApiKeyModalContents = (props:{
  apiKey: ApiKey,
  children?: React.ReactNode
}) => (
  <div className={cx('new-api-key-modal')}>
    <p>
      Below are your seceret and access keys.
      Once you close this modal, the seceret key will no longer be available.
    </p>
    <InputWithCopyButton label="Access Key" value={props.apiKey.accessKey} />
    <InputWithCopyButton label="Secret Key" value={props.apiKey.secretKey || ''} />
    {props.children}
  </div>
)

export const DeleteApiKeyModal = (props: {
  apiKey: ApiKey,
  userSlug: string,
  onRequestClose: () => void,
  onDeleted: () => void,
}) => {
  const formComponentProps = useForm({
    onSuccess: () => { props.onRequestClose(); props.onDeleted() },
    handleSubmit: () => deleteApiKey({ userSlug: props.userSlug, accessKey: props.apiKey.accessKey }),
  })

  return (
    <Modal title="Delete API Key" onRequestClose={props.onRequestClose}>
      <Form submitText="Delete API Key" cancelText="Close" onCancel={props.onRequestClose} {...formComponentProps}>
        <p>Are you sure you want to delete this API key?</p>
      </Form>
    </Modal>
  )
}
