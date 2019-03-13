// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { deleteUserAuthenticationScheme } from 'src/services'

import ChallengeModalForm from 'src/components/challenge_modal_form'


export const DeleteAuthModal = (props: {
  userSlug: string,
  schemeCode: string,
  onRequestClose: () => void,
}) => <ChallengeModalForm
    modalTitle="Delete Authentication"
    warningText="This will remove this authentication method, preventing login. This cannot be undone."
    submitText="Delete"
    challengeText={props.schemeCode}
    handleSubmit={() => deleteUserAuthenticationScheme({
      userSlug: props.userSlug,
      authSchemeName: props.schemeCode
    })}
    onRequestClose={props.onRequestClose}
  />
