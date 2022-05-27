// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'

import { ServiceWorker } from 'src/global_types'
import {
  renderModals,
  useModal,
  useWiredData,
} from 'src/helpers'
import { listServiceWorkers, testServiceWorker } from 'src/services'

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

  const [testDataState, dispatchTestData] = React.useReducer(testDataReducer, {})

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
                  cellOrder(worker, testDataState[worker.name] ?? initialTestData, {
                    showDeleteModal: (worker) => deleteModal.show({ worker }),
                    showEditModal: (worker) => editModal.show({ worker }),
                    testService: async (worker) => {
                      dispatchTestData({ type: 'start', worker: worker.name })
                      let passedTest = true
                      try {
                        const data = await testServiceWorker({ id: worker.id })
                        dispatchTestData({
                          type: 'finish',
                          worker: worker.name,
                          passedTest: data.live,
                          message: data.message,
                        })
                      }
                      catch (err) {
                        dispatchTestData({
                          type: 'finish',
                          worker: worker.name,
                          passedTest: false,
                          message: err,
                        })
                      }
                    }
                  }).map((v, colIndex) => (
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
  testService: (worker: ServiceWorker) => void
}

const emptyActions: Actions = {
  showDeleteModal: () => { },
  showEditModal: () => { },
  testService: () => { },
}

function cellOrder(): Array<string>
function cellOrder(worker: ServiceWorker, testData: TestData, actions: Actions): Array<React.ReactNode>
function cellOrder(worker?: ServiceWorker, testData?: TestData, actions?: Actions): Array<string | React.ReactNode> {
  const {
    showEditModal: showEdit,
    showDeleteModal: showDelete,
    testService,
  } = (actions ?? emptyActions)

  return [
    worker?.name ?? "Name",
    (worker?.deleted ?? "Deleted").toString(),
    testData == undefined ? "Status" : (<WorkerStatusIcon {...testData} />),
    worker == undefined ? "Actions" : (
      <>
        <ButtonGroup>
          <Button small onClick={() => showEdit(worker)}>Edit</Button>
          <Button small onClick={() => testService(worker)}>Test</Button>
          <Button danger small onClick={() => showDelete(worker)}>Delete</Button>
        </ButtonGroup>
      </>
    ),
  ]
}

type TestData = {
  isTesting: boolean,
  testResult: 'connected' | 'offline' | null
  testMessage: string
}

const WorkerStatusIcon = (props: TestData) => {
  const { isTesting, testResult } = props

  if (isTesting) {
    return (
      <div>...</div>
    )
  } else if (testResult != null) {
    const title = props.testMessage.trim() == ""
      ? undefined
      : props.testMessage.split(/>>(.*)/, 2).join("\n") // put error on new line if present
    return (
      <div title={title}>{testResult === 'connected' ? "Working" : "Offline"}</div>
    )
  }

  return null
}

const initialTestData: TestData = {
  isTesting: false,
  testResult: null,
  testMessage: "",
}

type TestDataState = Record<string, TestData>

const testDataReducer = (state: TestDataState, action: TestDataAction): TestDataState => {
  if (action.type == 'start') {
    return {
      ...state,
      [action.worker]: { isTesting: true, testResult: null, testMessage: "Testing..." }
    }
  }
  if (action.type == 'finish') {
    return {
      ...state,
      [action.worker]: {
        isTesting: false,
        testResult: action.passedTest ? 'connected' : 'offline',
        testMessage: action.message
      }
    }
  }
  return state
}

type TestDataAction =
  | TestDataActionStartTest
  | TestDataActionFinishTest

type TestDataActionStartTest = {
  type: 'start'
  worker: string
}
type TestDataActionFinishTest = {
  type: 'finish'
  passedTest: boolean
  worker: string
  message: string
}
