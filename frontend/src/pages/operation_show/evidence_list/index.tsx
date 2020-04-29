// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { Evidence } from 'src/global_types'
import { RouteComponentProps } from 'react-router-dom'
import { ViewName } from '../types'
import { useDataSource, getEvidenceList } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'

import Layout from '../layout'
import Timeline from 'src/components/timeline'
import { EditEvidenceModal, DeleteEvidenceModal, ChangeFindingsOfEvidenceModal } from '../evidence_modals'

export default (props: RouteComponentProps<{slug: string}>) => {
  const ds = useDataSource()
  const {slug} = props.match.params
  const query: string = new URLSearchParams(props.location.search).get('q') || ''

  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList(ds, {
    operationSlug: slug,
    query: query,
  }), [ds, slug, query]))

  const editModal = useModal<{evidence: Evidence}>(modalProps => (
    <EditEvidenceModal {...modalProps} operationSlug={slug} onEdited={wiredEvidence.reload} />
  ))
  const deleteModal = useModal<{evidence: Evidence}>(modalProps => (
    <DeleteEvidenceModal {...modalProps} operationSlug={slug} onDeleted={wiredEvidence.reload} />
  ))
  const assignToFindingsModal = useModal<{evidence: Evidence}>(modalProps => (
    <ChangeFindingsOfEvidenceModal {...modalProps} operationSlug={slug} onChanged={() => {/* no need to reload here */}} />
  ))

  const navigate = (view: ViewName, query: string) => {
    let path = `/operations/${slug}/${view}`
    if (query != '') path += `?q=${query}`
    props.history.push(path)
  }

  return (
    <Layout
      onEvidenceCreated={wiredEvidence.reload}
      onNavigate={navigate}
      operationSlug={slug}
      query={query}
      view="evidence"
    >
      {wiredEvidence.render(evidence => (
        <Timeline
          evidence={evidence}
          actions={{
            'Edit': evidence => editModal.show({evidence}),
            'Delete': evidence => deleteModal.show({evidence}),
            'Assign Findings': evidence => assignToFindingsModal.show({evidence}),
          }}
          onQueryUpdate={query => navigate('evidence', query)}
          operationSlug={slug}
          query={query}
        />
      ))}

      {renderModals(editModal, deleteModal, assignToFindingsModal)}
    </Layout>
  )
}
