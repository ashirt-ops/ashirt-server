// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddGlobalVarModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
}) => {
  const [newGlobalVar, setNewGlobalVar] = React.useState<boolean>(false)

  return (
    <SettingsSection title="New Global Variable Creation">
      <em>
        Creates a global variable, which can be used by service workers to customize the behavior of the application.
      </em>
      <Button primary onClick={() => setNewGlobalVar(true)}>Create New Global Variable</Button>
      {newGlobalVar && <AddGlobalVarModal onRequestClose={() => {
        setNewGlobalVar(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
