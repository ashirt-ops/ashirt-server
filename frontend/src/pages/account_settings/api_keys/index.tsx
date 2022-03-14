// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { format } from 'date-fns'

import Button from 'src/components/button'
import Form from 'src/components/form'
import Modal from 'src/components/modal'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { ApiKey } from 'src/global_types'
import { InputWithCopyButton } from 'src/components/text_copiers'
import { UserWithAuth } from 'src/global_types'
import { getApiKeys, createApiKey, deleteApiKey, rotateApiKey } from 'src/services'
import { useWiredData, useForm } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  profile: UserWithAuth
}) => {
  const [deleteKey, setDeleteKey] = React.useState<null | ApiKey>(null)
  const [rotateKey, setRotateKey] = React.useState<null | ApiKey>(null)
  const wiredApiKeys = useWiredData<Array<ApiKey>>(React.useCallback(() => getApiKeys({ userSlug: props.profile.slug }), [props.profile.slug]))

  return (
    <SettingsSection title="API Key Management" width="wide">
      {wiredApiKeys.render(apiKeys => (
        <Table columns={['Access Key', 'Secret Key', 'Last Used', 'Actions']} className={cx('table')}>
          {apiKeys.map(apiKey => (
            <tr key={apiKey.accessKey}>
              <td><span className={cx('monospace')}>{apiKey.accessKey}</span></td>
              <td><span className={cx('monospace')}>{apiKey.secretKey || '**************'}</span></td>
              <td>{apiKey.lastAuth ? format(apiKey.lastAuth, "MMMM do, yyyy 'at' HH:mm:ss") : 'Never'}</td>
              <td>
                <Button small onClick={() => setRotateKey(apiKey)}>Rotate</Button>
                <Button small danger onClick={() => setDeleteKey(apiKey)}>Delete</Button>
              </td>
            </tr>
          ))}
        </Table>
      ))}
      <br />
      <GenerateKeyButton userSlug={props.profile.slug} onKeyCreated={wiredApiKeys.reload} />
      {deleteKey && (
        <DeleteApiKeyModal
          userSlug={props.profile.slug}
          apiKey={deleteKey}
          onRequestClose={() => setDeleteKey(null)}
          onDeleted={wiredApiKeys.reload}
        />
      )}
      {rotateKey && (
        <RotateApiKeyModal
          userSlug={props.profile.slug}
          apiKey={rotateKey}
          onRequestClose={() => setRotateKey(null)}
          onRotated={wiredApiKeys.reload}
        />
      )}
    </SettingsSection>
  )
}

const GenerateKeyButton = (props: {
  userSlug: string,
  onKeyCreated: () => void,
}) => {
  const [apiKey, setApiKey] = React.useState<null | ApiKey>(null)
  const generateKeyForm = useForm({
    onSuccess: props.onKeyCreated,
    handleSubmit: () => createApiKey({ userSlug: props.userSlug }).then(setApiKey),
  })

  return <>
    <Form submitText="Generate New API Key" {...generateKeyForm} />
    {apiKey && (
      <Modal title="New API Key" onRequestClose={() => setApiKey(null)}>
        <div className={cx('new-api-key-modal')}>
          <p>
            Below are your seceret and access keys.
            Once you close this modal, the seceret key will no longer be available.
          </p>
          <InputWithCopyButton label="Access Key" value={apiKey.accessKey} />
          <InputWithCopyButton label="Secret Key" value={apiKey.secretKey || ''} />
          <Button primary onClick={() => setApiKey(null)}>Close</Button>
        </div>
      </Modal>
    )}
  </>
}

const DeleteApiKeyModal = (props: {
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

const RotateApiKeyModal = (props: {
  apiKey: ApiKey,
  userSlug: string,
  onRequestClose: () => void,
  onRotated: () => void,
}) => {
  const [updatedApiKey, setUpdatedApiKey] = React.useState<null | ApiKey>(null)

  const closeModal = () => {
    setUpdatedApiKey(null)
    props.onRequestClose()
  }
  const formComponentProps = useForm({
    onSuccess: () => {
      props.onRotated()
    },
    handleSubmit: () => rotateApiKey({ userSlug: props.userSlug, accessKey: props.apiKey.accessKey }).then(setUpdatedApiKey),
  })

  return (
    <Modal title="Rotate API Key" onRequestClose={closeModal}>
      <Form
        submitText="Rotate API Key"
        cancelText="Close"
        disableSubmit={!!updatedApiKey}
        onCancel={closeModal}
        {...formComponentProps}>
        {updatedApiKey != null
          ? (
            <div className={cx('new-api-key-modal')}>
              <p>
                Below are your seceret and access keys.
                Once you close this modal, the seceret key will no longer be available.
              </p>
              <InputWithCopyButton label="Access Key" value={updatedApiKey.accessKey} />
              <InputWithCopyButton label="Secret Key" value={updatedApiKey.secretKey || ''} />
            </div>
          )
          : (
            <p>This will delete and re-create an API key. Are you sure you want to do this?</p>
          )
        }
      </Form>
    </Modal>
  )
}
