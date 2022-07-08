// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Sidebar from './sidebar'
import { Toolbar } from './toolbar'
import classnames from 'classnames/bind'
import { CreateEvidenceModal } from '../evidence_modals'
import { CreateFindingModal } from '../finding_modals'
import { ViewName } from 'src/global_types'
import { useModal, useWiredData, renderModals } from 'src/helpers'
import { NavToFunction } from 'src/helpers/navigate-to-query'
import { BuildReloadBus } from 'src/helpers/reload_bus'
import { getSavedQueries, getOperation } from 'src/services'
const cx = classnames.bind(require('./stylesheet'))

const noOp = () => { }

export type CreateButtonPosition = "sidebar-inline" | "sidebar-above" | "filter" | "none"

export default (props: {
  children: React.ReactNode,
  onEvidenceCreated?: () => void,
  onFindingCreated?: () => void,
  onNavigate: NavToFunction,
  operationSlug: string,
  query: string,
  view: ViewName,
}) => {
  const reloadBus = BuildReloadBus()

  const [expanded, setExpanded] = React.useState(false)
  const createEvidenceModal = useModal<{}>(modalProps => (
    <CreateEvidenceModal {...modalProps} onCreated={props.onEvidenceCreated || noOp} operationSlug={props.operationSlug} />
  ))
  const createFindingModal = useModal<{}>(modalProps => (
    <CreateFindingModal {...modalProps} onCreated={props.onFindingCreated || noOp} operationSlug={props.operationSlug} />
  ))

  const wiredData = useWiredData(React.useCallback(() => Promise.all([
    getSavedQueries({ operationSlug: props.operationSlug }),
    getOperation(props.operationSlug),
  ]), [props.operationSlug]))

  React.useEffect(() => {
    reloadBus.onReload(wiredData.reload)
    return () => { reloadBus.offReload(wiredData.reload) }
  })

  const showCreateButtons: CreateButtonPosition = 'filter'

  return wiredData.render(([queries, operation]) => {
    const currentQuery = queries.find(q => q.type == props.view && q.query == props.query)
    return (
      <div className={cx('root')}>
        <div className={cx(expanded ? 'expanded-toolbar' : 'toolbar')}>
          <Toolbar
            operationSlug={props.operationSlug}
            query={props.query}
            onSearch={query => props.onNavigate(props.view, query)}
            expandedView={expanded}
            setExpandedView={setExpanded}
            viewName={props.view}
            onRequestCreateFinding={() => createFindingModal.show({})}
            onRequestCreateEvidence={() => createEvidenceModal.show({})}
            showCreateButtons={showCreateButtons}
            requestQueriesReload={reloadBus.requestReload}
            queryName={currentQuery?.name}
          />
        </div>
        <div className={cx('sidebar')}>
          <Sidebar
            onRequestCreateFinding={() => createFindingModal.show({})}
            onRequestCreateEvidence={() => createEvidenceModal.show({})}
            currentQuery={props.query}
            currentView={props.view}
            onNavigate={props.onNavigate}
            showCreateButtons={showCreateButtons}
            requestQueriesReload={reloadBus.requestReload}
            operation={operation}
            queries={queries}
          />
        </div>
        <div className={cx('children')}>
          {props.children}
        </div>

        {renderModals(createEvidenceModal, createFindingModal)}
      </div>
    )
  }
  )
}
