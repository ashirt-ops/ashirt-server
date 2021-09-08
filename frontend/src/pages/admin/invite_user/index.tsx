// Copyright 2021, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { InviteUserModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
}) => {
  const [newUser, setNewUser] = React.useState<boolean>(false)

  return (
    <SettingsSection title="Invite user">
      <em>
        Pre-provision accounts for expected users.
        This will create a new account, and place that account in "recovery" mode. Provide
        the generated URL to the user and instruct them to link an account once signing in.
      </em>
      <Button primary onClick={() => setNewUser(true)}>Invite a User</Button>
      {newUser && <InviteUserModal onRequestClose={() => {
        setNewUser(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
