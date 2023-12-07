import * as React from 'react'
import * as dateFns from 'date-fns'

import Input from 'src/components/input'
import SettingsSection from 'src/components/settings_section'
import classnames from 'classnames/bind'
import { useForm, useFormField } from 'src/helpers/use_form'
import { renderModals, useModal, useWiredData } from 'src/helpers'
import { beginAddCredential, deleteWebauthnCredential, finishAddCredential, listWebauthnCredentials, modifyCredentialName } from '../services'
import Table from 'src/components/table'
import Button, { ButtonGroup } from 'src/components/button'
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
    <SettingsSection className={cx('security-credentials-section')} title="WebAuthn Security Credentials" width="narrow">
      <CredentialList {...bus} />
      <AddCredentialButton {...bus} />
    </SettingsSection>
  )
}

const CredentialList = (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const wiredCredentials = useWiredData(listWebauthnCredentials)

  React.useEffect(() => {
    props.onReload(wiredCredentials.reload)
    return () => { props.offReload(wiredCredentials.reload) }
  })

  const deleteModal = useModal<{ credentialId: string, credentialName: string }>(mProps => <DeleteCredentialModal {...mProps} />, wiredCredentials.reload)
  const modifyModal = useModal<{ credentialName: string }>(mProps => <EditCredentialModal {...mProps} />, wiredCredentials.reload)

  return (<>
    {wiredCredentials.render(data => {
      return (
        <div>
          <Table columns={['Credential Name', 'Date Created', 'Actions']}>
            {data.credentials.map(credentialEntry => {
              const { credentialName, dateCreated, credentialId } = credentialEntry
              return (
                <tr key={credentialName}>
                  <td>{credentialName}</td>
                  <td>{toEnUSDate(dateCreated)}</td>
                  <td className={cx('button-cell')}>
                    <ButtonGroup className={cx('row-buttons')}>
                      <Button small onClick={() => {
                        modifyModal.show({ credentialName })
                      }}>Edit</Button>
                      <Button danger small onClick={() => {
                        deleteModal.show({ credentialId, credentialName })
                      }}>Delete</Button>
                    </ButtonGroup>
                  </td>
                </tr>
              )
            })}
          </Table>
          {renderModals(deleteModal, modifyModal)}
        </div>
      )
    })}
  </>)
}

const AddCredentialButton = (props: {
  requestReload: () => void
}) => {
  const createModal = useModal(mProps => (
    <AddCredentialModal {...mProps} />
  ), props.requestReload)

  return (
    <div>
      <Button primary onClick={createModal.show}>Register new security credential</Button>
      {renderModals(createModal)}
    </div>
  )
}

const AddCredentialModal = (props: {
  onRequestClose: () => void,
}) => {
  const credentialName = useFormField("")

  const formComponentProps = useForm({
    fields: [credentialName],
    handleSubmit: async () => {
      if (credentialName.value === '') {
        return Promise.reject(new Error("Credential name must be populated"))
      }
      const reg = await beginAddCredential({
        credentialName: credentialName.value,
      })
      const credOptions = convertToCredentialCreationOptions(reg)

      const signed = await navigator.credentials.create(credOptions)

      if (signed == null || signed.type != 'public-key') {
        throw new Error("WebAuthn is not supported")
      }
      const pubCredential = signed as PublicKeyCredential
      const pubCredentialResponse = pubCredential.response as AuthenticatorAttestationResponse

      await finishAddCredential({
        type: 'public-key',
        id: pubCredential.id,
        rawId: encodeAsB64(pubCredential.rawId),
        response: {
          attestationObject: encodeAsB64(pubCredentialResponse.attestationObject),
          clientDataJSON: encodeAsB64(pubCredentialResponse.clientDataJSON),
        },
      })
    },
    onSuccess: props.onRequestClose
  })

  return (
    <ModalForm
      title={"Add Security Credential"}
      submitText={"Create"}
      cancelText="Cancel"
      onRequestClose={props.onRequestClose}
      {...formComponentProps}
    >
      <Input label="Credential name" {...credentialName} />
    </ModalForm>
  )
}

const DeleteCredentialModal = (props: {
  credentialName: string,
  credentialId: string,
  onRequestClose: () => void,
}) => (
  <ChallengeModalForm
    modalTitle="Delete Credential"
    warningText="Are you sure you want to delete this security credential?"
    submitText="Delete"
    challengeText={props.credentialName}
    handleSubmit={() => deleteWebauthnCredential({ credentialId: props.credentialId })}
    onRequestClose={props.onRequestClose}
  />
)

const EditCredentialModal = (props: {
  credentialName: string,
  onRequestClose: () => void,
}) => {
  const credentialName = useFormField("")

  const formComponentProps = useForm({
    fields: [credentialName],
    handleSubmit: async () => {
      if (credentialName.value === '') {
        return Promise.reject(new Error("Credential name must be populated"))
      }
      await modifyCredentialName({
        newCredentialName: credentialName.value,
        credentialName: props.credentialName,
      })
    },
    onSuccess: props.onRequestClose
  })

  return (
    <ModalForm
      title={"Edit Credential Name"}
      submitText={"Edit"}
      cancelText="Cancel"
      onRequestClose={props.onRequestClose}
      {...formComponentProps}
    >
      <Input label="New Credential name" {...credentialName} />
    </ModalForm>
  )
}
