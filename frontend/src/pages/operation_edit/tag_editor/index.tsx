import { useCallback } from 'react'
import { useWiredData } from 'src/helpers'
import { getTags } from 'src/services'

import SettingsSection from 'src/components/settings_section'
import { OperationTagTable } from 'src/components/tag_editor'

export default function TagEditor(props: { operationSlug: string }) {
  const wiredTags = useWiredData(
    useCallback(() => getTags({ operationSlug: props.operationSlug }), [props.operationSlug]),
  )

  return (
    <SettingsSection title="Operation Tags">
      {wiredTags.render((tags) => (
        <OperationTagTable
          operationSlug={props.operationSlug}
          tags={tags}
          onUpdate={wiredTags.reload}
        />
      ))}
    </SettingsSection>
  )
}
