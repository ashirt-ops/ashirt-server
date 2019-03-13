// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddHeadlessUserModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
}) => {
  const [newHeadlessUser, setNewHeadlessUser] = React.useState<boolean>(false)

  return (
    <SettingsSection title="Headless User Creation">
      <em>Headless users cannot login as regular users, but can access services via API keys.</em>
      <Button primary onClick={() => setNewHeadlessUser(true)}>Create New Headless User</Button>
      {newHeadlessUser && <AddHeadlessUserModal onRequestClose={() => {
        setNewHeadlessUser(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
