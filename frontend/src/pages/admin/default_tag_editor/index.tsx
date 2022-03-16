// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import SettingsSection from 'src/components/settings_section'
import { getDefaultTags } from 'src/services'
import { useWiredData } from 'src/helpers'
import { DefaultTagTable } from 'src/components/tag_editor'

export const DefaultTagEditor = (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const wiredTags = useWiredData(getDefaultTags)

  React.useEffect(() => {
    props.onReload(wiredTags.reload)
    return () => { props.offReload(wiredTags.reload) }
  })

  return (
    <SettingsSection title="Initial Operation Tags">
      {wiredTags.render(tags => (
        <DefaultTagTable
          tags={tags}
          onUpdate={wiredTags.reload}
        />
      ))}
    </SettingsSection>
  )
}
