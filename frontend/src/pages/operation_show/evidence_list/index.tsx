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
import { Codeblock, DenormalizedTag, Evidence, ExportedEvidence, Media } from 'src/global_types'
import { useNavigate, useLocation, useParams } from 'react-router-dom'
import { getEvidenceList } from 'src/services'
import { useWiredData, useModal, renderModals } from 'src/helpers'
import { mkNavTo } from 'src/helpers/navigate-to-query'

import { saveAs } from 'file-saver';
import _ from 'lodash'
import Modal from 'src/components/modal'
import { languageToFileExtension } from 'src/helpers/languages'
const JSZip = require("jszip");

export default () => {
  const { slug } = useParams<{ slug: string }>()
  const operationSlug = slug! // useParams puts everything in a partial, so our type above doesn't matter.
  const location = useLocation()
  const navigate = useNavigate()
  const [evidence, setEvidence] = React.useState<Evidence[]>([])

  const query: string = new URLSearchParams(location.search).get('q') || ''
  const [lastEditedUuid, setLastEditedUuid] = React.useState("")
  const [showModal, setShowModal] = React.useState(false)

  const wiredEvidence = useWiredData(React.useCallback(() => getEvidenceList({
    operationSlug,
    query,
  }), [operationSlug, query]))

  React.useEffect(() => {
    wiredEvidence.expose(data => data.length && setEvidence(data))
  }, [wiredEvidence])

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

  const contentToFileExtension = {
    "image": "jpeg",
    "terminal-recording": "cast",
    "codeblock": languageToFileExtension,
    "event": "txt",
    "none": "txt",
    "http-request-cycle": "har",
  }

  const getMedia = async (evidence: ExportedEvidence[]): Promise<[void | Media[], ExportedEvidence[]]> => {
    const media = await Promise.all(evidence.map(async (e) => {
      const rawMedia = await fetch(`/web/operations/${operationSlug}/evidence/${e.uuid}/media`)
      if (rawMedia.status !== 200) { throw new Error("Error downloading media") }

      // remove tag IDs from the evidence
      e.tags.forEach(t => {
        delete (t as DenormalizedTag).id;
        delete (t as DenormalizedTag).colorName;
      })

      if (e.contentType === "codeblock") {
        const data: Codeblock = await rawMedia.json();
        e.fileName = `${e.uuid}.${contentToFileExtension[e.contentType][data?.contentSubtype]}`;
        e.sourceFileName= data?.metadata?.source
        return {
          uuid: e.uuid,
          contentType: e.contentType, 
          contentSubtype: data.contentSubtype,
          blob: new Blob([data.content], {type: `text/${contentToFileExtension[e.contentType][data.contentSubtype]}`})
        }
      } else {
        const blob = await rawMedia.blob()
        e.fileName = `${e.uuid}.${contentToFileExtension[e.contentType]}`;
        return {
          uuid: e.uuid,
          contentType: e.contentType, 
          blob
        }
      }
    }))
    .catch(() => {
      setShowModal(true)
    });
    return [media, evidence]
  }

  const exportEvidence = async () => {
    var zip = new JSZip();
    var evidenceFolder = zip.folder("evidence");
    const evidenceCopy = _.cloneDeep(evidence)
    const [media, modEvidenceCopy] = await getMedia(evidenceCopy)

    zip.file("evidence.json", JSON.stringify(modEvidenceCopy))

    if (media){
      media.forEach((mb) => {
        const fileName = mb.contentType === "codeblock" 
          ? `${mb.uuid}.${contentToFileExtension[mb.contentType][mb.contentSubtype!]}` 
          : `${mb.uuid}.${contentToFileExtension[mb.contentType]}`;
        evidenceFolder.file(fileName, mb.blob, {base64: true});
      })    
      const zipFile = await zip.generateAsync({type:"blob"})
      saveAs(zipFile, `evidence-${operationSlug}-${new Date().toISOString()}.zip`);
    }
  }

  return (
    <Layout
      onEvidenceCreated={reloadToTop}
      onNavigate={navTo}
      operationSlug={operationSlug}
      query={query}
      view="evidence"
      exportEvidence={exportEvidence}
    >
      {showModal && <Modal smallerWidth={true} title='Evidence Download Error' onRequestClose={() => setShowModal(false)}>
        <p>Error downloading evidence - please try again later</p>
      </Modal>}
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
