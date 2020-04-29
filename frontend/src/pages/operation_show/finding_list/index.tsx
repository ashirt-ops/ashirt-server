// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { Finding } from 'src/global_types'
import { RouteComponentProps } from 'react-router-dom'
import { ViewName } from '../types'
import { useDataSource, getFindings } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'

import FindingsTable from './findings_table'
import Layout from '../layout'
import { DeleteFindingModal, EditFindingModal } from '../finding_modals'

export default (props: RouteComponentProps<{slug: string}>) => {
  const ds = useDataSource()
  const {slug} = props.match.params
  const query: string = new URLSearchParams(props.location.search).get('q') || ''

  const wiredFindings = useWiredData(React.useCallback(() => getFindings(ds, {
    operationSlug: slug,
    query: query,
  }), [ds, slug, query]))
  React.useEffect(wiredFindings.reload, [slug, query])

  const navigate = (view: ViewName, query: string) => {
    let path = `/operations/${slug}/${view}`
    if (query != '') path += `?q=${query}`
    props.history.push(path)
  }

  const editFindingModal = useModal<{finding: Finding}>(modalProps => (
    <EditFindingModal {...modalProps} onEdited={wiredFindings.reload} operationSlug={slug} />
  ))
  const deleteFindingModal = useModal<{finding: Finding}>(modalProps => (
    <DeleteFindingModal {...modalProps} onDeleted={wiredFindings.reload} operationSlug={slug} />
  ))

  return (
    <Layout
      onFindingCreated={wiredFindings.reload}
      onNavigate={navigate}
      operationSlug={slug}
      query={query}
      view="findings"
    >
      {wiredFindings.render(findings => (
        <div style={{padding: 20}}>
          <FindingsTable
            findings={findings}
            onDelete={finding => deleteFindingModal.show({finding})}
            onEdit={finding => editFindingModal.show({finding})}
            operationSlug={slug}
          />
        </div>
      ))}

      {renderModals(editFindingModal, deleteFindingModal)}
    </Layout>
  )
}
