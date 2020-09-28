// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { totpIsEnabled, generateTotpSecret, setTotp, deleteTotp } from 'src/services'
import { useWiredData, useForm, useFormField } from 'src/helpers'

import Form from 'src/components/form'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import SettingsSection from 'src/components/settings_section'
import Button from 'src/components/button'
import { InputWithCopyButton } from 'src/components/text_copiers'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
}) => {
  const wiredTotpIsEnabled = useWiredData(totpIsEnabled)

  return (
    <SettingsSection title="Two-Factor Authentication" className={cx('root')}>
      {wiredTotpIsEnabled.render(isEnabled => (
        isEnabled
          ? <TotpEnabled onSuccess={wiredTotpIsEnabled.reload} />
          : <TotpDisabled onSuccess={wiredTotpIsEnabled.reload} />
      ))}
    </SettingsSection>
  )
}

// This component renders a button to disable totp from the user's account. It is used
// for users who currently have totp enabled
const TotpEnabled = (props: {
  onSuccess: () => void,
}) => {
  const removeTotpForm = useForm({
    handleSubmit: deleteTotp,
    onSuccess: props.onSuccess,
  })

  return (
    <Form submitDanger submitText="Remove Two-Factor Auth" {...removeTotpForm}>
      Your account currently has TOTP enabled
    </Form>
  )
}

// This component renders a button to open the totp setup modal. It is used for users
// who don't currently have totp enabled
const TotpDisabled = (props: {
  onSuccess: () => void,
}) => {
  const [showTotpSetupModal, setShowTotpSetupModal] = React.useState(false)

  return <>
    <p>Two-Factor auth is not currently setup on your account.</p>
    <p>
      Two-factor authentication adds an additional layer of security to
      your account by requiring more than just a password to log in.
    </p>
    <Button primary onClick={() => setShowTotpSetupModal(true)}>Setup Two-Factor Authentication</Button>

    {showTotpSetupModal && (
      <TotpSetupModal
        {...showTotpSetupModal}
        onSuccess={props.onSuccess}
        onRequestClose={() => setShowTotpSetupModal(false)}
      />
    )}
  </>
}

// This component renders a modal that guides a user to setting up totp on their account.
// On mount it makes a call to the backend to generate a new totp secret.
// It displays the qr code and form to finish totp setup.
const TotpSetupModal = (props: {
  onSuccess: () => void,
  onRequestClose: () => void,
}) => {
  const wiredGeneratedTotp = useWiredData(generateTotpSecret)

  return (
    <Modal title="Enable Two-Factor Authentication" onRequestClose={props.onRequestClose}>
      {wiredGeneratedTotp.render(generatedTotp => (
        <div className={cx('setup-modal')}>
          <h2>1. Set up your authenticator app</h2>
          <p>
            Set up two factor auth with your authenticator app of choice by scanning
            the QR code below.
          </p>
          <img className={cx('qr')} src={generatedTotp.qr} />

          <p>If you are unable to scan the QR code you can enter this URI instead:</p>
          <InputWithCopyButton label="OTP Auth URI" value={generatedTotp.url} />

          <br /><br />

          <h2>2. Enter the code from your authenticator app</h2>
          <p>
            To finalize setup, verify your app is generating correct two-factor codes
            by entering a generated code below:
          </p>
          <TotpSetupForm {...props} secret={generatedTotp.secret} onCancel={props.onRequestClose} />
        </div>
      ))}
    </Modal>
  )
}

// This component renders a form that actually sets up totp on a user's account with a given secret
// It requests the user enter a valid one time passcode and sends up the totp passcode with the passed
// secret for the backend to validate. If the passcode is a valid one for the secret for the current time
// the backend enables totp on the current user's account.
const TotpSetupForm = (props: {
  onSuccess: () => void,
  onCancel: () => void,
  secret: string,
}) => {
  const passcodeField = useFormField("")

  const setupTotpForm = useForm({
    fields: [passcodeField],
    handleSubmit: () => {
      return setTotp({
        secret: props.secret,
        passcode: passcodeField.value,
      })
    },
    onSuccess: props.onSuccess,
  })

  return (
    <Form submitText="Finish" cancelText="Cancel" onCancel={props.onCancel} {...setupTotpForm}>
      <Input label="One Time Passcode" placeholder="123456" {...passcodeField} />
    </Form>
  )
}
