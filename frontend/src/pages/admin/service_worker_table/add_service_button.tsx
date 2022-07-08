// Copyright 2022, Yahoo Inc
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { renderModals, useModal, } from 'src/helpers'

import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddEditServiceWorkerModal } from './modals'


export default (props: {
  requestReload?: () => void
}) => {
  const editModal = useModal<{}>(mProps => (
    <AddEditServiceWorkerModal {...mProps} />
  ), props.requestReload)

  return (
    <SettingsSection title="Add Service Worker">
      <Button primary onClick={() => editModal.show({})}>Create New Service Worker</Button>

      {renderModals(editModal)}
    </SettingsSection>
  )
}
