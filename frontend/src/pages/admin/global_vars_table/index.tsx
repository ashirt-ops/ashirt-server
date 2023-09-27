// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { WiredData} from 'src/helpers'

import { GlobalVar } from 'src/global_types'
import { getGlobalVars } from 'src/services'

import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import { StandardPager } from 'src/components/paging'
import SettingsSection from 'src/components/settings_section'
import { default as Menu } from 'src/components/menu'
import { ClickPopover } from 'src/components/popover'
import Input from 'src/components/input'
import { DeleteGlobalVarModal, ModifyGlobalVarModal } from 'src/pages/admin_modals'
import { useWiredData } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const [deletingGlobalVar, setDeletingGlobalVar] = React.useState<null | GlobalVar>(null)
  const [modifyingGlobalVar, setModifyingGlobalVar] = React.useState<null | GlobalVar>(null)
  const itemsPerPage = 10
  const [page, setPage] = React.useState(1)
  const [pageLength, setPageLength] = React.useState(0)

  const [usernameFilterValue, setUsernameFilterValue] = React.useState('')

  const columns = Object.keys(rowBuilder(null, <span />, <span />))

  const wiredGlobalVars = useWiredData<GlobalVar[]>(
    React.useCallback(() => getGlobalVars(), [usernameFilterValue]),
    (err: Error) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )

  React.useEffect(() => {
    props.onReload(wiredGlobalVars.reload)
    return () => { props.offReload(wiredGlobalVars.reload) }
  })
  React.useEffect(() => {
    wiredGlobalVars.expose(data => setPageLength(Math.ceil(data.length / itemsPerPage)))
  }, [wiredGlobalVars])

  return (
    <SettingsSection title="Global Variables List" width="wide">
      <Table className={cx('table')} columns={columns}>
        {wiredGlobalVars.render(data => <>
          {data?.map((globalVar, i) => <TableRow key={globalVar.name} data={rowBuilder(globalVar, usersInGroup(wiredGlobalVars, globalVar), modifyActions(globalVar, setDeletingGlobalVar, setModifyingGlobalVar))} />
          )}
        </>)}
      </Table>

      {deletingGlobalVar && <DeleteGlobalVarModal globalVar={deletingGlobalVar} onRequestClose={() => { setDeletingGlobalVar(null); wiredGlobalVars.reload() }} />}
      {modifyingGlobalVar && <ModifyGlobalVarModal globalVar={modifyingGlobalVar} onRequestClose={() => { setModifyingGlobalVar(null); wiredGlobalVars.reload() }} />}
    </SettingsSection>
  )
}

const TableRow = (props: { data: Rowdata }) => (
  <tr>
    <td>{props.data["Name"]}</td>
    <td>{props.data["Value"]}</td>
    <td>{props.data["Actions"]}</td>
  </tr>
)

type Rowdata = {
  "Name": string,
  "Value": string,
  "Actions": JSX.Element,
}

const rowBuilder = (u: GlobalVar | null, users: JSX.Element, actions: JSX.Element): Rowdata => ({
  "Name": u ? u.name : "",
  "Value": u ? u.value : "",
  "Actions": actions,
})

const usersInGroup = (
  wiredGlobalVars: WiredData<GlobalVar[]>,
  u: GlobalVar
) => {
  const count = wiredGlobalVars.render(data => <span>{data.length}</span>)
  return (
    <ButtonGroup>
      <ClickPopover className={cx('popover')} closeOnContentClick content={
        <Menu>
          {wiredGlobalVars.render(data => {
            const varList = data.map(globalVar => <p key={globalVar.name} className={cx('user')}>{globalVar.name}</p>)
            return <>{varList}</>
      })}
        </Menu>
      }>
        <Button small className={cx('arrow')}><p className={cx('button-text')}>{count} Global Variables</p></Button>
      </ClickPopover>
    </ButtonGroup>
  )
}

const modifyActions = (
  u: GlobalVar,
  onDeleteClick: (u: GlobalVar) => void,
  onEditClick: (u: GlobalVar) => void
) => {
  return (
    <ButtonGroup className={cx('row-buttons')}>
      <Button small onClick={() => onEditClick(u)}>Edit</Button>
      <Button small danger={true} onClick={() => onDeleteClick(u)}>Delete</Button>
    </ButtonGroup>
  )
}
