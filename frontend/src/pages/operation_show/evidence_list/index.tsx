// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Layout from '../layout'
import Timeline from 'src/components/timeline'
import {
  EditEvidenceModal,
  DeleteEvidenceModal,
  ChangeFindingsOfEvidenceModal,
  MoveEvidenceModal,
  EvidenceMetadataModal,
} from '../evidence_modals'
import { Evidence } from 'src/global_types'
import { useNavigate, useLocation, useParams } from 'react-router-dom'
import { getEvidenceList } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'
import { mkNavTo } from 'src/helpers/navigate-to-query'

import { saveAs } from 'file-saver';
var JSZip = require("jszip");

export default () => {
  const { slug } = useParams<{ slug: string }>()
  const operationSlug = slug! // useParams puts everything in a partial, so our type above doesn't matter.
  const location = useLocation()
  const navigate = useNavigate()

  const query: string = new URLSearchParams(location.search).get('q') || ''
  const [lastEditedUuid, setLastEditedUuid] = React.useState("")

  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList({
    operationSlug,
    query,
  }), [operationSlug, query]))

  const reloadToTop = () => {
    setLastEditedUuid("")
    wiredEvidence.reload()
  }

  const editModal = useModal<{ evidence: Evidence }>(modalProps => (
    <EditEvidenceModal {...modalProps} operationSlug={operationSlug} onEdited={() => {
      setLastEditedUuid(modalProps.evidence.uuid)
      wiredEvidence.reload()
    }} />
  ))
  const viewModal = useModal<{ evidence: Evidence }>(modalProps => (
    <EvidenceMetadataModal {...modalProps} operationSlug={operationSlug} onUpdated={wiredEvidence.reload} />
  ))
  const deleteModal = useModal<{ evidence: Evidence }>(modalProps => (
    <DeleteEvidenceModal {...modalProps} operationSlug={operationSlug} onDeleted={reloadToTop} />
  ))
  const assignToFindingsModal = useModal<{ evidence: Evidence }>(modalProps => (
    <ChangeFindingsOfEvidenceModal {...modalProps} operationSlug={operationSlug} onChanged={() => {/* no need to reload here */ }} />
  ))

  const moveModal = useModal<{ evidence: Evidence }>(modalProps => (
    <MoveEvidenceModal {...modalProps} operationSlug={operationSlug} onEvidenceMoved={() => { }} />
  ))

  const navTo = mkNavTo({
    navTo: navigate,
    slug: operationSlug
  })

  const zipIt = async () => {
    var zip = new JSZip();
    zip.file("Hello.txt", "Hello World\n");
    console.log("zipping?")
    var img = zip.folder("images");
    // const imgData = await fetch('https://pbs.twimg.com/media/FuYBCUHWIAA4LkB?format=jpg&name=4096x4096')
    // img.file("smile.gif", imgData, {base64: true});
    const content = await zip.generateAsync({type:"blob"})
    saveAs(content, "example.zip");
    console.log("finsihed zipping?", content)
    // .then(function(content) {
    //     // see FileSaver.js
        
    // });
  }
  return (
    <Layout
      onEvidenceCreated={reloadToTop}
      onNavigate={navTo}
      operationSlug={operationSlug}
      query={query}
      view="evidence"
    >
      <button onClick={zipIt}>Zip it</button>
      {wiredEvidence.render(evidence => (
        <Timeline
          scrollToUuid={lastEditedUuid}
          evidence={evidence}
          actions={[
            { label: "Edit", act: evidence => editModal.show({ evidence }) },
            { label: "Assign Findings", act: evidence => assignToFindingsModal.show({ evidence }) },
          ]}
          extraActions={[
            { label: 'Move', act: evidence => moveModal.show({ evidence }) },
            { label: 'Delete', act: evidence => deleteModal.show({ evidence }) },
            { label: "Metadata", act: evidence => viewModal.show({ evidence }) },
          ]}
          onQueryUpdate={query => navTo('evidence', query)}
          operationSlug={operationSlug}
          query={query}
        />
      ))}

      {renderModals(editModal, deleteModal, assignToFindingsModal, moveModal, viewModal)}
    </Layout>
  )
}
