import { useState } from 'react'
import Button from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { AddVarModal } from 'src/pages/admin_modals'
export default function AddVariable(props: { requestReload?: () => void; operationSlug?: string }) {
  const [newVar, setNewVar] = useState<boolean>(false)

  return (
    <SettingsSection title="New Variable Creation">
      <em>
        Creates a variable, which can be used by service workers to customize the behavior of the
        application.
      </em>
      {newVar && (
        <AddVarModal
          operationSlug={props.operationSlug}
          onRequestClose={() => {
            setNewVar(false)
            props.requestReload && props.requestReload()
          }}
        />
      )}
      <Button primary onClick={() => setNewVar(true)}>
        Create New Variable
      </Button>
    </SettingsSection>
  )
}
