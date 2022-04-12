// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Layout from '../layout'
import Timeline from 'src/components/timeline'
import {
  EditEvidenceModal,
  DeleteEvidenceModal,
  ChangeFindingsOfEvidenceModal,
  MoveEvidenceModal,
  ViewEvidenceMetadataModal,
  AddEvidenceMetadataModal,
  EvidenceMetadataModal,
} from '../evidence_modals'
import { Evidence, ViewName } from 'src/global_types'
import { RouteComponentProps } from 'react-router-dom'
import { getEvidenceList } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'

export default (props: RouteComponentProps<{ slug: string }>) => {
  const { slug } = props.match.params
  const query: string = new URLSearchParams(props.location.search).get('q') || ''
  const [lastEditedUuid, setLastEditedUuid] = React.useState("")

  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList({
    operationSlug: slug,
    query: query,
  }), [slug, query]))

  const reloadToTop = () => {
    setLastEditedUuid("")
    wiredEvidence.reload()
  }

  const editModal = useModal<{ evidence: Evidence }>(modalProps => (
    <EditEvidenceModal {...modalProps} operationSlug={slug} onEdited={() => {
      setLastEditedUuid(modalProps.evidence.uuid)
      wiredEvidence.reload()
    }} />
  ))
  const viewModal = useModal<{ evidence: Evidence }>(modalProps => (
    <EvidenceMetadataModal {...modalProps} operationSlug={slug} onUpdated={wiredEvidence.reload} />
  ))
  const createMetadataModal = useModal<{ evidence: Evidence }>(modalProps => (
    <AddEvidenceMetadataModal {...modalProps} operationSlug={slug} onCreated={wiredEvidence.reload} />
  ))
  const deleteModal = useModal<{ evidence: Evidence }>(modalProps => (
    <DeleteEvidenceModal {...modalProps} operationSlug={slug} onDeleted={reloadToTop} />
  ))
  const assignToFindingsModal = useModal<{ evidence: Evidence }>(modalProps => (
    <ChangeFindingsOfEvidenceModal {...modalProps} operationSlug={slug} onChanged={() => {/* no need to reload here */ }} />
  ))

  const moveModal = useModal<{ evidence: Evidence }>(modalProps => (
    <MoveEvidenceModal {...modalProps} operationSlug={slug} onEvidenceMoved={() => { }} />
  ))

  const navigate = (view: ViewName, query: string) => {
    let path = `/operations/${slug}/${view}`
    if (query != '') path += `?q=${encodeURIComponent(query.trim())}`
    props.history.push(path)
  }

  return (
    <Layout
      onEvidenceCreated={reloadToTop}
      onNavigate={navigate}
      operationSlug={slug}
      query={query}
      view="evidence"
    >
      {wiredEvidence.render(evidence => (
        <Timeline
          scrollToUuid={lastEditedUuid}
          evidence={evidence}
          actions={[
            {
              label: "Edit",
              act: evidence => editModal.show({ evidence }),
            },
            {
              label: "Metadata",
              act: evidence => viewModal.show({ evidence }),
              canAct: (evidence) => evidence.metadata.length > 0
                ? { disabled: false }
                : { disabled: true, title: "No metadata available" },
            },
            {
              label: "Assign Findings",
              act: evidence => assignToFindingsModal.show({ evidence }),
            },
          ]}
          extraActions={[
            { label: 'Move', act: evidence => moveModal.show({ evidence }) },
            { label: 'Delete', act: evidence => deleteModal.show({ evidence }) },
            { label: 'Create Metadata', act: evidence => createMetadataModal.show({ evidence }) }
          ]}
          onQueryUpdate={query => navigate('evidence', query)}
          operationSlug={slug}
          query={query}
        />
      ))}

      {renderModals(editModal, deleteModal, assignToFindingsModal, moveModal, viewModal, createMetadataModal)}
    </Layout>
  )
}
