// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddVarModal } from 'src/pages/admin_modals'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  requestReload?: () => void
  operationSlug?: string
}) => {
  const [newVar, setNewVar] = React.useState<boolean>(false)

  return (
    <SettingsSection title="New Variable Creation">
      <em>
        Creates a variable, which can be used by service workers to customize the behavior of the application.
      </em>
      {newVar && <AddVarModal operationSlug={props.operationSlug} onRequestClose={() => {
        setNewVar(false)
        props.requestReload && props.requestReload()
      }} />}
      <Button primary onClick={() => setNewVar(true)}>Create New Variable</Button>
    </SettingsSection>
  )
}
