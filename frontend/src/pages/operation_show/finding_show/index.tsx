import * as React from 'react'
import FindingInfo from './finding_info'
import Timeline from 'src/components/timeline'
import classnames from 'classnames/bind'
import {ChangeEvidenceOfFindingModal, RemoveEvidenceFromFindingModal, EditFindingModal, DeleteFindingModal} from '../finding_modals'
import {EditEvidenceModal} from '../evidence_modals'
import {Evidence, Finding} from 'src/global_types'
import {useNavigate, useParams} from 'react-router-dom'
import {default as Button, ButtonGroup} from 'src/components/button'
import {getFinding} from 'src/services'
import {useWiredData, useModal, renderModals} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default () => {
  const { slug, uuid } = useParams<{ slug: string, uuid: string }>()
  // useParams puts everything in a partial, so our type above doesn't matter.
  const operationSlug = slug!
  const findingUuid = uuid!

  const navigate = useNavigate()
  const wiredFinding = useWiredData(React.useCallback(() => getFinding({
    operationSlug,
    findingUuid,
  }), [operationSlug, findingUuid]))
  const [lastEditedUuid, setLastEditedUuid] = React.useState("")

  const reloadToTop = () => {
    setLastEditedUuid("")
    wiredFinding.reload()
  }

  const addRemoveEvidenceModal = useModal<{finding: Finding, initialEvidence: Array<Evidence>}>(modalProps => (
    <ChangeEvidenceOfFindingModal {...modalProps} onChanged={reloadToTop} operationSlug={operationSlug} />
  ))
  const editFindingModal = useModal<{finding: Finding}>(modalProps => (
    <EditFindingModal {...modalProps} onEdited={reloadToTop} operationSlug={operationSlug} />
  ))
  const deleteFindingModal = useModal<{finding: Finding}>(modalProps => (
    <DeleteFindingModal {...modalProps} onDeleted={() => navigate(`/operations/${operationSlug}/findings`)} operationSlug={operationSlug} />
  ))
  const editEvidenceModal = useModal<{evidence: Evidence}>(modalProps => (
    <EditEvidenceModal {...modalProps} onEdited={ ()=>{
      setLastEditedUuid(modalProps.evidence.uuid)
      wiredFinding.reload()
    }} operationSlug={operationSlug} />
  ))
  const removeEvidenceFromFindingModal = useModal<{evidence: Evidence, finding: Finding}>(modalProps => (
    <RemoveEvidenceFromFindingModal {...modalProps} onRemoved={reloadToTop} operationSlug={operationSlug} />
  ))

  return <>
    {wiredFinding.render(({finding, evidence}) => (
      <div className={cx('root')}>
        <div className={cx('finding-info')}>
          <div className={cx('actions')}>
            <Button small className={cx('left')} icon={require('./back.svg')} onClick={() => navigate(-1)}>Back</Button>
            <ButtonGroup className={cx('right')}>
              <Button small onClick={() => addRemoveEvidenceModal.show({ finding, initialEvidence: evidence })}>Add/Remove Evidence</Button>
              <Button small onClick={() => editFindingModal.show({ finding })}>Edit</Button>
              <Button small onClick={() => deleteFindingModal.show({ finding })}>Delete</Button>
            </ButtonGroup>
          </div>
          <FindingInfo finding={finding} />
        </div>
        <div className={cx('timeline')}>
          <Timeline
            scrollToUuid={lastEditedUuid}
            evidence={evidence}
            actions={[
              {
                label: 'Remove From Finding',
                act: evidence => removeEvidenceFromFindingModal.show({ evidence, finding })
              },
              {
                label: 'Edit',
                act: evidence => editEvidenceModal.show({ evidence })
              },
            ]}
            onQueryUpdate={query => navigate(`/operations/${operationSlug}/evidence?q=${encodeURIComponent(query.trim())}`)}
            operationSlug={operationSlug}
            query=""
          />
        </div>
      </div>
    ))}
    {renderModals(addRemoveEvidenceModal, editFindingModal, deleteFindingModal, editEvidenceModal, removeEvidenceFromFindingModal)}
  </>
}
