import * as React from 'react'
import classnames from 'classnames/bind'
import { format } from 'date-fns'

import Button from 'src/components/button'
import Form from 'src/components/form'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { ApiKey } from 'src/global_types'
import { UserWithAuth } from 'src/global_types'
import { getApiKeys, createApiKey } from 'src/services'
import { useWiredData, useForm } from 'src/helpers'

import { NewApiKeyModal, DeleteApiKeyModal } from './modals'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  profile: UserWithAuth
}) => {
  const [deleteKey, setDeleteKey] = React.useState<null | ApiKey>(null)
  const wiredApiKeys = useWiredData<Array<ApiKey>>(React.useCallback(() => getApiKeys({ userSlug: props.profile.slug }), [props.profile.slug]))

  return (
    <SettingsSection title="API Key Management" width="wide">
      {wiredApiKeys.render(apiKeys => (
        <Table columns={['Access Key', 'Secret Key', 'Last Used', 'Actions']} className={cx('table')}>
          {apiKeys.map(apiKey => (
            <tr key={apiKey.accessKey}>
              <td><span className={cx('monospace')}>{apiKey.accessKey}</span></td>
              <td><span className={cx('monospace')}>{apiKey.secretKey || '**************'}</span></td>
              <td>{apiKey.lastAuth ? format(apiKey.lastAuth, "MMMM do, yyyy 'at' HH:mm:ss") : 'Never'}</td>
              <td><Button small onClick={() => setDeleteKey(apiKey)}>Delete</Button></td>
            </tr>
          ))}
        </Table>
      ))}
      <br />
      <GenerateKeyButton userSlug={props.profile.slug} onKeyCreated={wiredApiKeys.reload} />
      {deleteKey && (
        <DeleteApiKeyModal
          userSlug={props.profile.slug}
          apiKey={deleteKey}
          onRequestClose={() => setDeleteKey(null)}
          onDeleted={wiredApiKeys.reload}
        />
      )}
    </SettingsSection>
  )
}

const GenerateKeyButton = (props: {
  userSlug: string,
  onKeyCreated: () => void,
}) => {
  const [apiKey, setApiKey] = React.useState<null | ApiKey>(null)
  const generateKeyForm = useForm({
    onSuccess: props.onKeyCreated,
    handleSubmit: () => createApiKey({ userSlug: props.userSlug }).then(setApiKey),
  })

  return <>
    <Form submitText="Generate New API Key" {...generateKeyForm} />
    {apiKey && (
      <NewApiKeyModal apiKey={apiKey} onRequestClose={ () => setApiKey(null)} />
    )}
  </>
}
