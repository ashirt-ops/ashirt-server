// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { UserOwnView, SupportedAuthenticationScheme, AuthenticationInfo } from 'src/global_types'
import { useWiredData } from 'src/helpers'
import { getSupportedAuthentications } from 'src/services'
import { format } from 'date-fns'

import { DeleteAuthModal } from '../modals'
import Modal from 'src/components/modal'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import { useAuthFrontendComponent } from 'src/authschemes'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  profile: UserOwnView,
  allowLinking: boolean,
  requestReload?: () => void,
}) => {
  const [removeAuth, setRemovingAuth] = React.useState<null | string>(null)
  const wiredSchemes = useWiredData<Array<SupportedAuthenticationScheme>>(getSupportedAuthentications)

  const columns = [
    "Method Name",
    "Status",
    "Last Login",
    "Actions",
  ]

  return (
    <SettingsSection title="Authentication Methods" width="wide">
      {wiredSchemes.render(supportedSchemes => (
        <Table className={cx('table')} columns={columns}>
          {supportedSchemes.map(scheme => (
            <TableRow
              allowLinking={props.allowLinking}
              key={scheme.schemeCode}
              supportedScheme={scheme}
              authInfo={props.profile.authSchemes}
              removeAuth={() => setRemovingAuth(scheme.schemeCode)}
              requestReload={props.requestReload}
            />
          ))}
        </Table>
      ))}

      {removeAuth && <DeleteAuthModal
        userSlug={props.profile.slug}
        schemeCode={removeAuth}
        onRequestClose={() => { setRemovingAuth(null); props.requestReload && props.requestReload() }}
      />}
    </SettingsSection>
  )
}

const TableRow = (props: {
  supportedScheme: SupportedAuthenticationScheme,
  authInfo: Array<AuthenticationInfo>,
  removeAuth: () => void,
  allowLinking: boolean,
  requestReload?: () => void,
}) => {
  const [linking, setLinking] = React.useState<boolean>(false)
  const Linker = useAuthFrontendComponent(props.supportedScheme.schemeCode, 'Linker')

  const userScheme = props.authInfo.find(x => x.schemeCode === props.supportedScheme.schemeCode)
  const canDeleteAuth = () => {
    switch (true) {
      case (!userScheme): return { disabled: true, title: "Auth scheme has not been linked" }
      default: return {}
    }
  }
  const canLink = () => {
    switch (true) {
      case (userScheme != undefined): return { disabled: true, title: "Auth scheme has already been linked" }
      default: return {}
    }
  }

  return (
    <tr>
      <td>{props.supportedScheme.schemeName}</td>
      <td>{userScheme ? "Linked" : "Not Linked"}</td>
      <td>{userScheme && userScheme.lastLogin
        ? format(userScheme.lastLogin, "MMMM do, yyyy")
        : "Never"} </td>
      <td>
        <ButtonGroup>
          {props.allowLinking && <Button small {...canLink()} onClick={() => setLinking(true)}>Link</Button>}
          <Button danger small {...canDeleteAuth()} onClick={() => props.removeAuth()}>Delete</Button>
        </ButtonGroup>
      </td>
      {linking && (
        <Modal onRequestClose={() => setLinking(false)} title={"Link Account"}>
          <Linker
            onSuccess={() => {
              setLinking(false)
              props.requestReload && props.requestReload()
            }}
            authFlags={props.supportedScheme.schemeFlags}
          />
        </Modal>
      )}
    </tr>
  )
}
