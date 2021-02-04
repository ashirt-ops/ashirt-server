// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Sidebar from './sidebar'
import Toolbar from './toolbar'
import classnames from 'classnames/bind'
import {CreateEvidenceModal} from '../evidence_modals'
import {CreateFindingModal} from '../finding_modals'
import {ViewName} from 'src/global_types'
import {useModal, renderModals} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

const noOp = () => {}

export default (props: {
  children: React.ReactNode,
  onEvidenceCreated?: () => void,
  onFindingCreated?: () => void,
  onNavigate: (view: ViewName, query: string) => void,
  operationSlug: string,
  query: string,
  view: ViewName,
}) => {
  const createEvidenceModal = useModal<void>(modalProps => (
    <CreateEvidenceModal {...modalProps} onCreated={props.onEvidenceCreated || noOp} operationSlug={props.operationSlug} />
  ))
  const createFindingModal = useModal<void>(modalProps => (
    <CreateFindingModal {...modalProps} onCreated={props.onFindingCreated || noOp} operationSlug={props.operationSlug} />
  ))

  return (
    <div className={cx('root')}>
      <div className={cx('toolbar')}>
        <Toolbar
          onRequestCreateFinding={() => createFindingModal.show()}
          onRequestCreateEvidence={() => createEvidenceModal.show()}
          onSearch={query => props.onNavigate(props.view, query)}
          query={props.query}
          operationSlug={props.operationSlug}
          viewName={props.view}
        />
      </div>
      <div className={cx('sidebar')}>
        <Sidebar
          currentQuery={props.query}
          currentView={props.view}
          onNavigate={props.onNavigate}
          operationSlug={props.operationSlug}
        />
      </div>
      <div className={cx('children')}>
        {props.children}
      </div>

      {renderModals(createEvidenceModal, createFindingModal)}
    </div>
  )
}
