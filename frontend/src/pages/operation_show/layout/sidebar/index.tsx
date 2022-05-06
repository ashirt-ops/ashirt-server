// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ActionMenu from './action_menu'
import OperationBadges from 'src/components/operation_badges'
import classnames from 'classnames/bind'
import {Link} from 'react-router-dom'
import {NewQueryModal, EditQueryModal, DeleteQueryModal} from './query_modal'
import { SavedQuery, SavedQueryType, ViewName} from 'src/global_types'
import {default as ListMenu, ListItem, ListItemWithSaveButton, ListItemWithMenu} from 'src/components/list_menu'
import {getSavedQueries, getOperation} from 'src/services'
import {useWiredData, useModal, renderModals} from 'src/helpers'
import { default as Button, ButtonGroup } from 'src/components/button'
import { CreateButtonPosition } from '..'
import { NavToFunction } from 'src/helpers/navigate-to-query'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  currentQuery: string,
  currentView: ViewName,
  onNavigate: NavToFunction,
  onRequestCreateFinding: () => void,
  onRequestCreateEvidence: () => void,
  onReload: (listener: () => void) => void
  offReload: (listener: () => void) => void
  operationSlug: string,
  showCreateButtons: CreateButtonPosition
}) => {
  const wiredQueries = useWiredData(React.useCallback(() => Promise.all([
    getSavedQueries({operationSlug: props.operationSlug}),
    getOperation(props.operationSlug),
  ]), [props.operationSlug]))

  React.useEffect(() => {
    props.onReload(wiredQueries.reload)
    return () => { props.offReload(wiredQueries.reload) }
  })

  return wiredQueries.render(([queries, operation]) => (
    <div className={cx('root')}>
      <header>
        <h1 title={operation.name}>{operation.name}</h1>
        <Link className={cx('edit')} to={`/operations/${props.operationSlug}/edit`} title="Edit this operation" />
        <Link className={cx('overview')} to={`/operations/${props.operationSlug}/overview`} title="View evidence overview" />
        <OperationBadges {...operation} />
      </header>
      {props.showCreateButtons == 'sidebar-above' && (
        <ButtonGroup className={cx('create-evi-finding-group')}>
          <Button size="medium" onClick={props.onRequestCreateEvidence}>Create Evidence</Button>
          <Button size="medium" onClick={props.onRequestCreateFinding}>Create Finding</Button>
        </ButtonGroup>
      )}
      <QueryList
        addNew={props.onRequestCreateEvidence}
        name="Evidence"
        type="evidence"
        onSelectQuery={props.onNavigate.bind(null, 'evidence')}
        savedQueries={queries.filter(q => q.type === 'evidence')}
        onSavedQueryChange={wiredQueries.reload}
        {...props}
      />
      <QueryList
        addNew={props.onRequestCreateFinding}
        name="Findings"
        type="findings"
        onSelectQuery={props.onNavigate.bind(null, 'findings')}
        savedQueries={queries.filter(q => q.type === 'findings')}
        onSavedQueryChange={wiredQueries.reload}
        {...props}
      />
    </div>
  ))
}

const QueryList = (props: {
  addNew: () => void
  currentQuery: string,
  currentView: ViewName,
  name: string,
  onSavedQueryChange: () => void,
  onSelectQuery: (query: string, queryName?: string) => void,
  operationSlug: string,
  savedQueries: Array<SavedQuery>,
  type: SavedQueryType,
  showCreateButtons: CreateButtonPosition
}) => {
  const isThisView = props.currentView === props.type
  const currentQueryIsNew = (
    isThisView &&
    props.currentQuery !== '' &&
    !props.savedQueries.find(q => props.currentQuery === q.query)
  )

  const onCreated = () => {
    props.onSavedQueryChange()
  }
  const onEdited = (before: SavedQuery, after: SavedQuery) => {
    if (before.query === props.currentQuery && before.query !== after.query) {
      // Navigate to new query if the current selected query was edited
      props.onSelectQuery(after.query)
    }
    props.onSavedQueryChange()
  }
  const onDeleted = (before: SavedQuery) => {
    if (before.query === props.currentQuery) {
        // Navigate to "All" if the current selected query was deleted
        props.onSelectQuery('')
    }
    props.onSavedQueryChange()
  }

  const newQueryModal = useModal<void>(modalProps => (
    <NewQueryModal {...modalProps} operationSlug={props.operationSlug} query={props.currentQuery} type={props.type} onCreated={onCreated} />
  ))
  const editQueryModal = useModal<{savedQuery: SavedQuery}>(modalProps => (
    <EditQueryModal {...modalProps} operationSlug={props.operationSlug} onEdited={onEdited} view={props.type}/>
  ))
  const deleteQueryModal = useModal<{savedQuery: SavedQuery}>(modalProps => (
    <DeleteQueryModal {...modalProps} operationSlug={props.operationSlug} onDeleted={onDeleted} />
  ))

  return <>
    <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'baseline'}}>
      <h2>{props.name}</h2>
      {props.showCreateButtons === 'sidebar-inline' && (
        <Button small onClick={props.addNew}>Add New</Button>
      )}
    </div>
    <ListMenu>
      <ListItem
        name={`All ${props.name}`}
        selected={isThisView && props.currentQuery === ''}
        onSelect={() => props.onSelectQuery('')}
      />

      {currentQueryIsNew && (
        <ListItemWithSaveButton
          name={props.currentQuery}
          selected // If this is displayed it is always selected
          onSelect={() => {}}
          onSave={() => newQueryModal.show()}
        />
      )}

      {props.savedQueries.map(savedQuery => (
        <ListItemWithMenu
          key={savedQuery.id}
          name={savedQuery.name}
          selected={isThisView && props.currentQuery === savedQuery.query}
          onSelect={() => props.onSelectQuery(savedQuery.query, savedQuery.name)}
          menu={(
            <ActionMenu
              name={savedQuery.name}
              query={savedQuery.query}
              onEdit={() => editQueryModal.show({savedQuery})}
              onDelete={() => deleteQueryModal.show({savedQuery})}
            />
          )}
        />
      ))}
    </ListMenu>

    {renderModals(newQueryModal, editQueryModal, deleteQueryModal)}
  </>
}
