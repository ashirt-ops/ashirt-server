import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddUserModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
}) => {
  const [newUser, setNewUser] = React.useState<boolean>(false)

  return (
    <SettingsSection title="New User Creation">
      <em>
        Pre-provision accounts for expected users.
        This will create a local authentication account, and provide their first password.
        Users will be forced to reset their password on their next login.
      </em>
      <Button primary onClick={() => setNewUser(true)}>Create New User</Button>
      {newUser && <AddUserModal onRequestClose={() => {
        setNewUser(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
