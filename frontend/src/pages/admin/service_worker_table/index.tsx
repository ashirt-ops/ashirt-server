// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import { ServiceWorker } from 'src/global_types'
import {
  renderModals,
  useModal,
  useWiredData,
} from 'src/helpers'
import { listServiceWorkers } from 'src/services/service_workers'

import SettingsSection from 'src/components/settings_section'
import {
  default as Table,
  ErrorRow,
  LoadingRow,
} from 'src/components/table'
import { default as Button, ButtonGroup } from 'src/components/button'
import {
  AddEditServiceWorkerModal,
  DeleteServiceModal,
} from './modals'

type WorkerModal = { worker: ServiceWorker }

export default (props: {
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
}) => {
  const columns: Array<string> = cellOrder()

  const wiredServiceWorkers = useWiredData<Array<ServiceWorker>>(
    listServiceWorkers,
    (err) => <ErrorRow span={columns.length} error={err} />,
    () => <LoadingRow span={columns.length} />
  )
  const deleteModal = useModal<WorkerModal>(mProps => (
    <DeleteServiceModal {...mProps} />
  ), wiredServiceWorkers.reload)
  const editModal = useModal<WorkerModal>(mProps => (
    <AddEditServiceWorkerModal {...mProps} />
  ), wiredServiceWorkers.reload)


  React.useEffect(() => {
    props.onReload(wiredServiceWorkers.reload)
    return () => { props.offReload(wiredServiceWorkers.reload) }
  })

  return (
    <SettingsSection title="Service Worker List" width="wide">
      <Table columns={columns}>
        {wiredServiceWorkers.render(data => <>
          {
            data.map((worker) => (
              <tr key={worker.name}>
                {
                  cellOrder(worker, {
                    showDeleteModal: (worker) => deleteModal.show({ worker }),
                    showEditModal: (worker) => editModal.show({ worker })
                  })
                    .map((v, colIndex) => (
                      <td key={worker.name + ":" + columns[colIndex]}>{v}</td>
                    ))
                }
              </tr>
            ))
          }
        </>)}
      </Table>
      {renderModals(deleteModal, editModal)}
    </SettingsSection>
  )
}

type Actions = {
  showDeleteModal: (worker: ServiceWorker) => void
  showEditModal: (worker: ServiceWorker) => void
}

const emptyActions: Actions = {
  showDeleteModal: () => { },
  showEditModal: () => { },
}

function cellOrder(): Array<string>
function cellOrder(worker: ServiceWorker, actions: Actions): Array<React.ReactNode>
function cellOrder(worker?: ServiceWorker, actions?: Actions): Array<string | React.ReactNode> {
  const { showEditModal: showEdit, showDeleteModal: showDelete } = (actions ?? emptyActions)

  return [
    worker?.name ?? "Name",
    (worker?.deleted ?? "Deleted").toString(),
    worker == undefined ? "Actions" : (
      <ButtonGroup>
        <Button small onClick={() => showEdit(worker)}>Edit</Button>
        <Button danger small onClick={() => showDelete(worker)}>Delete</Button>
      </ButtonGroup>
    ),
  ]
}
