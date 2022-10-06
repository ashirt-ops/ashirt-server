// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ActionMenu from './action_menu'
import OperationBadges from 'src/components/operation_badges'
import OperationBadgesModal from 'src/components/operation_badges_modal'
import classnames from 'classnames/bind'
import { Link } from 'react-router-dom'
import { DeleteQueryModal } from './query_modal'
import { Operation, SavedQuery, SavedQueryType, ViewName} from 'src/global_types'
import { default as ListMenu, ListItem, ListItemWithSaveButton, ListItemWithMenu} from 'src/components/list_menu'
import { useModal, renderModals} from 'src/helpers'
import { default as Button, ButtonGroup } from 'src/components/button'
import { CreateButtonPosition } from '..'
import { NavToFunction } from 'src/helpers/navigate-to-query'
import { SaveQueryModal } from 'src/components/filter_fields/filter-field-grid'
import { setFavorite } from 'src/services'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  currentQuery: string,
  currentView: ViewName,
  onNavigate: NavToFunction,
  onRequestCreateFinding: () => void,
  onRequestCreateEvidence: () => void,
  showCreateButtons: CreateButtonPosition
  queries: Array<SavedQuery>
  operation: Operation
  requestQueriesReload?: () => void
}) => {
  const {operation, queries} = props
  const [isFavorite, setIsFavorite] = React.useState(operation.favorite)

  React.useEffect(() => {
    setFavorite(operation.slug, isFavorite)
  }, [operation.slug, isFavorite])

  const moreDetailsModal = useModal<{}>(modalProps => (
    <OperationBadgesModal {...modalProps} topContribs={operation?.topContribs} evidenceCount={operation?.evidenceCount} status={operation?.status} />
  ))

  const handleDetailsModal = () => moreDetailsModal?.show({})

  return (
    <div className={cx('root')}>
      <header>
        <h1 title={operation.name}>{operation.name}</h1>
        <Link className={cx('edit')} to={`/operations/${operation.slug}/edit`} title="Edit this operation" />
        <Link className={cx('overview')} to={`/operations/${operation.slug}/overview`} title="View evidence overview" />
        <Button
          className={cx('favorite-button', isFavorite && 'filled')}
          onClick={() => setIsFavorite(!isFavorite)}
        >
        </Button>
        <OperationBadges {...operation} showDetailsModal={handleDetailsModal} />
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
        onSavedQueryChange={() => props.requestQueriesReload?.()}
        operationSlug={operation.slug}
        {...props}
      />
      <QueryList
        addNew={props.onRequestCreateFinding}
        name="Findings"
        type="findings"
        onSelectQuery={props.onNavigate.bind(null, 'findings')}
        savedQueries={queries.filter(q => q.type === 'findings')}
        onSavedQueryChange={() => props.requestQueriesReload?.() }
        operationSlug={operation.slug}
        {...props}
      />
      {renderModals(moreDetailsModal)}
    </div>
  )
}

const QueryList = (props: {
  addNew: () => void
  currentQuery: string,
  currentView: ViewName,
  name: string,
  onSavedQueryChange: (queryName?: string) => void,
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

  const onCreated = (queryName: string) => {
    props.onSavedQueryChange(queryName)
  }

  const onDeleted = (before: SavedQuery) => {
    if (before.query === props.currentQuery) {
        // Navigate to "All" if the current selected query was deleted
        props.onSelectQuery('')
    }
    props.onSavedQueryChange()
  }

  const saveQueryModal = useModal<{}>(modalProps => (
    <SaveQueryModal
      query={props.currentQuery}
      onSaved={(queryName: string) => {
        onCreated(queryName)
      }}
      operationSlug={props.operationSlug}
      view={props.currentView}
      {...modalProps}
    />
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
          onSave={() => saveQueryModal.show({})}
        />
      )}

      {props.savedQueries.map(savedQuery => (
        <ListItemWithMenu
          key={savedQuery.id}
          name={savedQuery.name}
          selected={isThisView && props.currentQuery === savedQuery.query}
          onSelect={() => props.onSelectQuery(savedQuery.query)}
          menu={(
            <ActionMenu
              name={savedQuery.name}
              query={savedQuery.query}
              onDelete={() => deleteQueryModal.show({savedQuery})}
            />
          )}
        />
      ))}
    </ListMenu>

    {renderModals(saveQueryModal, deleteQueryModal)}
  </>
}
