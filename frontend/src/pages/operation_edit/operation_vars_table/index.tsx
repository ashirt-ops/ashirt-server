import * as React from 'react'

import { OperationVar } from 'src/global_types'
import { getOperationVars } from 'src/services'

import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import SettingsSection from 'src/components/settings_section'
import { DeleteVarModal, ModifyVarModal, } from 'src/pages/admin_modals'
import { useWiredData } from 'src/helpers'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import CreateVarButton from "src/components/add_variable"

export default (props: {
  operationSlug: string,
  isAdmin: boolean,
}) => {
  const bus = BuildReloadBus()

  const [deletingOperationVar, setDeletingOperationVar] = React.useState<null | OperationVar>(null)
  const [modifyingOperationVar, setModifyingOperationVar] = React.useState<null | OperationVar>(null)

  const columns = Object.keys(rowBuilder(null, <span />))

  const wiredOperationVars = useWiredData<OperationVar[]>(
    React.useCallback(() => getOperationVars(props.operationSlug), []),
    (err: Error) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )

  React.useEffect(() => {
    bus.onReload(wiredOperationVars.reload)
    return () => { bus.offReload(wiredOperationVars.reload) }
  })

  return (
    <>
      <SettingsSection title="Operation Variables" width="wide">
        <Table columns={columns}>
          {wiredOperationVars.render(data => <>
            {data?.map((operationVar) => <TableRow key={operationVar.name} operationVar={operationVar} data={rowBuilder(operationVar, modifyActions(operationVar, setDeletingOperationVar, setModifyingOperationVar))} />
            )}
          </>)}
        </Table>
        {deletingOperationVar && <DeleteVarModal variableData={{variable: deletingOperationVar, operationSlug: props.operationSlug}} onRequestClose={() => { setDeletingOperationVar(null); wiredOperationVars.reload() }} />}
        {modifyingOperationVar && <ModifyVarModal variableData={{variable: modifyingOperationVar, operationSlug: props.operationSlug}} onRequestClose={() => { setModifyingOperationVar(null); wiredOperationVars.reload() }} />}
      </SettingsSection>
      <CreateVarButton {...bus} operationSlug={props.operationSlug} />
    </>
  )
}


const TableRow = (props: { 
  data: Rowdata, 
  operationVar: OperationVar, 
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

const rowBuilder = (u: OperationVar | null, actions: JSX.Element): Rowdata => ({
  "Name": u ? u.name : "",
  "Actions": actions,
})

const modifyActions = (
  u: OperationVar,
  onDeleteClick: (u: OperationVar) => void,
  onEditClick: (u: OperationVar) => void
) => {
  return (
    <ButtonGroup>
      <Button small onClick={() => onEditClick(u)}>Edit</Button>
      <Button small danger={true} onClick={() => onDeleteClick(u)}>Delete</Button>
    </ButtonGroup>
  )
}
