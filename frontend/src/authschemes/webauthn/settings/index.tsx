// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import * as dateFns from 'date-fns'

import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import classnames from 'classnames/bind'
import { useForm, useFormField } from 'src/helpers/use_form'
import { renderModals, useModal, useWiredData } from 'src/helpers'
import { beginAddKey, deleteWebauthnKey, finishAddKey, listWebauthnKeys } from '../services'
import Table from 'src/components/table'
import Button from 'src/components/button'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import ModalForm from 'src/components/modal_form'
import { convertToCredentialCreationOptions, encodeAsB64 } from '../helpers'
import ChallengeModalForm from 'src/components/challenge_modal_form'
const cx = classnames.bind(require('./stylesheet'))

const toEnUSDate = (d: Date) => dateFns.format(d, "MMM dd, yyyy")

export default (props: {
  username: string,
  authFlags?: Array<string>
}) => {
  const bus = BuildReloadBus()
  return (
    <SettingsSection className={cx('security-keys-section')} title="WebAuthn Security Keys" width="narrow">
      <KeyList {...bus} />
      <AddKeyButton {...bus} />
    </SettingsSection>
  )
}

const KeyList = (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const wiredKeys = useWiredData(listWebauthnKeys)

  React.useEffect(() => {
    props.onReload(wiredKeys.reload)
    return () => { props.offReload(wiredKeys.reload) }
  })

  const deleteModal = useModal<{ keyName: string }>(mProps => <DeleteKeyModal {...mProps} />, wiredKeys.reload)

  return (<>
    {wiredKeys.render(data => {
      return (
        <div>
          <Table columns={['Key Name', 'Date Created', 'Actions']}>
            {data.keys.map(keyEntry => {
              const { keyName, dateCreated } = keyEntry
              return (
                <tr key={keyName}>
                  <td>{keyName}</td>
                  <td>{toEnUSDate(dateCreated)}</td>
                  <td>
                    <Button small danger onClick={() => {
                      deleteModal.show({ keyName })
                    }}>
                      Delete
                    </Button>
                  </td>
                </tr>
              )
            })}
          </Table>
          {renderModals(deleteModal)}
        </div>
      )
    })}
  </>)
}

const AddKeyButton = (props: {
  requestReload: () => void
}) => {
  const createModal = useModal(mProps => (
    <AddKeyModal {...mProps} />
  ), props.requestReload)

  return (
    <div>
      <Button primary onClick={createModal.show}>Register new security key</Button>
      {renderModals(createModal)}
    </div>
  )
}

const AddKeyModal = (props: {
  onRequestClose: () => void,
}) => {
  const keyName = useFormField("")

  const formComponentProps = useForm({
    fields: [keyName],
    handleSubmit: async () => {
      if (keyName.value === '') {
        return Promise.reject(new Error("Key name must be populated"))
      }
      const reg = await beginAddKey({
        keyName: keyName.value,
      })
      const credOptions = convertToCredentialCreationOptions(reg)

      const signed = await navigator.credentials.create(credOptions)

      if (signed == null || signed.type != 'public-key') {
        throw new Error("WebAuthn is not supported")
      }
      const pubKeyCred = signed as PublicKeyCredential
      const pubKeyResponse = pubKeyCred.response as AuthenticatorAttestationResponse

      await finishAddKey({
        type: 'public-key',
        id: pubKeyCred.id,
        rawId: encodeAsB64(pubKeyCred.rawId),
        response: {
          attestationObject: encodeAsB64(pubKeyResponse.attestationObject),
          clientDataJSON: encodeAsB64(pubKeyResponse.clientDataJSON),
        },
      })
    },
    onSuccess: props.onRequestClose
  })

  return (
    <ModalForm
      title={"Add Security Key"}
      submitText={"Create"}
      cancelText="Cancel"
      onRequestClose={props.onRequestClose}
      {...formComponentProps}
    >
      <Input label="Key name" {...keyName} />
    </ModalForm>
  )
}

const DeleteKeyModal = (props: {
  keyName: string,
  onRequestClose: () => void,
}) => (
  <ChallengeModalForm
    modalTitle="Delete Key"
    warningText="Are you sure you want to delete this security key?"
    submitText="Delete"
    challengeText={props.keyName}
    handleSubmit={() => deleteWebauthnKey({ keyName: props.keyName })}
    onRequestClose={props.onRequestClose}
  />
)
