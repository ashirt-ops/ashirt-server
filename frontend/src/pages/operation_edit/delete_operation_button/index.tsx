// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import ChallengeModalForm from 'src/components/challenge_modal_form'
import { deleteOperation } from 'src/services'

export default (props: {
  operationSlug: string
}) => {
  const [showDeleteOpModal, setShowDeleteOpModal] = React.useState<boolean>(false)
  return (
    <SettingsSection title="Delete Operation">
      <Button primary danger onClick={() => setShowDeleteOpModal(true)}>Delete Operation</Button>
      {showDeleteOpModal && <DeleteOperationModal
        operationSlug={props.operationSlug}
        onRequestClose={(success) => {
          setShowDeleteOpModal(false)
          if (success) {
            window.location.href = "/operations"
          }
        }} />
      }
    </SettingsSection>
  )
}

export const DeleteOperationModal = (props: {
  operationSlug: string,
  onRequestClose: (success: boolean) => void,
}) => <ChallengeModalForm
    modalTitle="Delete Operation"
    warningText="This will permanently remove the operation from the system, along with associated evidence and findings."
    submitText="Delete"
    challengeText={props.operationSlug}
    handleSubmit={() => deleteOperation(props.operationSlug)}
    onRequestClose={props.onRequestClose}
  />
