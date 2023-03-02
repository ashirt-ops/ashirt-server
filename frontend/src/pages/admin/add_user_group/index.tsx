// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddUserGroupModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
}) => {
  const [newUserGroup, setNewUserGroup] = React.useState<boolean>(false)

  return (
    <SettingsSection title="New Group Creation">
      <em>
        Creates a new user group, which allows multiple users to be managed as a single entity.
      </em>
      <Button primary onClick={() => setNewUserGroup(true)}>Create New Group</Button>
      {newUserGroup && <AddUserGroupModal onRequestClose={() => {
        setNewUserGroup(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
