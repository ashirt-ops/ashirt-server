// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { GlobalVar } from 'src/global_types'
import { getGlobalVars } from 'src/services'

import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { DeleteGlobalVarModal, ModifyGlobalVarModal, } from 'src/pages/admin_modals'
import { useWiredData } from 'src/helpers'

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const [deletingGlobalVar, setDeletingGlobalVar] = React.useState<null | GlobalVar>(null)
  const [modifyingGlobalVar, setModifyingGlobalVar] = React.useState<null | GlobalVar>(null)

  const columns = Object.keys(rowBuilder(null, <span />))

  const wiredGlobalVars = useWiredData<GlobalVar[]>(
    React.useCallback(() => getGlobalVars(), []),
    (err: Error) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )

  React.useEffect(() => {
    props.onReload(wiredGlobalVars.reload)
    return () => { props.offReload(wiredGlobalVars.reload) }
  })

  return (
    <SettingsSection title="Global Variables List" width="wide">
      <Table columns={columns}>
        {wiredGlobalVars.render(data => <>
          {data?.map((globalVar) => <TableRow key={globalVar.name} globalVar={globalVar} data={rowBuilder(globalVar, modifyActions(globalVar, setDeletingGlobalVar, setModifyingGlobalVar))} />
          )}
        </>)}
      </Table>

      {deletingGlobalVar && <DeleteGlobalVarModal globalVar={deletingGlobalVar} onRequestClose={() => { setDeletingGlobalVar(null); wiredGlobalVars.reload() }} />}
      {modifyingGlobalVar && <ModifyGlobalVarModal globalVar={modifyingGlobalVar} onRequestClose={() => { setModifyingGlobalVar(null); wiredGlobalVars.reload() }} />}
    </SettingsSection>
  )
}


const TableRow = (props: { 
  data: Rowdata, 
  globalVar: GlobalVar, 
}) => (
    <tr>
      <td>{props.data["Name"]}</td>
      <td>{props.data["Actions"]}</td>
    </tr>
  )

type Rowdata = {
  "Name": string,
  "Actions": JSX.Element,
}

const rowBuilder = (u: GlobalVar | null, actions: JSX.Element): Rowdata => ({
  "Name": u ? u.name : "",
  "Actions": actions,
})

const modifyActions = (
  u: GlobalVar,
  onDeleteClick: (u: GlobalVar) => void,
  onEditClick: (u: GlobalVar) => void
) => {
  return (
    <ButtonGroup>
      <Button small onClick={() => onEditClick(u)}>Edit</Button>
      <Button small danger={true} onClick={() => onDeleteClick(u)}>Delete</Button>
    </ButtonGroup>
  )
}
