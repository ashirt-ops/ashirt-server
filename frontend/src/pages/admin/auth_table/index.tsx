// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import SettingsSection from 'src/components/settings_section'
import Table from 'src/components/table'
import { AuthSchemeDetails } from 'src/global_types'
import { DeleteGlobalAuthSchemeModal } from 'src/pages/admin_modals'
import { default as Button, ButtonGroup } from 'src/components/button'
import { formatDistanceToNow } from 'date-fns'
import { getSupportedAuthenticationDetails } from 'src/services'
import { useWiredData } from 'src/helpers'

type PurgableScheme = {
  schemeCode: string,
  uniqueUsers: number
}

export default (props: {}) => {
  const [purgeScheme, setPurgingScheme] = React.useState<PurgableScheme | null>(null)

  const wiredSchemes = useWiredData<Array<AuthSchemeDetails>>(getSupportedAuthenticationDetails)
  const columns = [
    'Scheme Name',
    { label: '# Users', title: "Number of users who can use this method" },
    { label: '# Unique Users', title: "Number of users who only use this method" },
    'Last Used',
    'Notes',
    'Actions',
  ]

  return (
    <SettingsSection title="Authentication Methods" width="wide">
      {wiredSchemes.render(supportedSchemes => {
        const nonEmptySchemes = supportedSchemes.filter( s => s.userCount > 0).length
        return (
          <Table columns={columns}>
            {supportedSchemes.map((item) => renderTableRow(item, nonEmptySchemes, setPurgingScheme))}
          </Table>
        )
      })}

      {purgeScheme && <DeleteGlobalAuthSchemeModal {...purgeScheme} onRequestClose={() => { setPurgingScheme(null); wiredSchemes.reload() }} />}
    </SettingsSection>
  )
}

const renderTableRow = (data: AuthSchemeDetails, nonEmptySchemeQuantity: number, purgeScheme: (i: PurgableScheme) => void) => {
  return (
    <tr key={data.schemeCode}>
      <td>{data.schemeName}</td>
      <td>{data.userCount}</td>
      <td>{data.uniqueUserCount}</td>
      <td>{data.lastUsed ? formatDistanceToNow(data.lastUsed, { addSuffix: true }) : "Never"}</td>
      <td>{data.labels.join(", ")}</td>
      <td>
        <ButtonGroup>
          <Button small danger
            disabled={nonEmptySchemeQuantity == 1 || data.userCount == 0}
            onClick={() => purgeScheme({ schemeCode: data.schemeCode, uniqueUsers: data.uniqueUserCount })}
          >Remove All Users</Button>
        </ButtonGroup>
      </td>
    </tr>
  )
}
