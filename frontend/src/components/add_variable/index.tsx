// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddVarModal } from 'src/pages/admin_modals'

export default (props: {
  requestReload?: () => void
  operationSlug?: string
}) => {
  const [newVar, setNewVar] = React.useState<boolean>(false)

  return (
    <SettingsSection title="New Variable Creation">
      <em>
        Creates a global variable, which can be used by service workers to customize the behavior of the application.
      </em>
      <Button primary onClick={() => setNewVar(true)}>Create New Variable</Button>
      {newVar && <AddVarModal operationSlug={props.operationSlug} onRequestClose={() => {
        setNewVar(false)
        props.requestReload && props.requestReload()
      }} />}
    </SettingsSection>
  )
}
