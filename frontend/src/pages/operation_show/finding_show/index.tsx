// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import FindingInfo from './finding_info'
import Timeline from 'src/components/timeline'
import classnames from 'classnames/bind'
import {ChangeEvidenceOfFindingModal, RemoveEvidenceFromFindingModal, EditFindingModal, DeleteFindingModal} from '../finding_modals'
import {EditEvidenceModal} from '../evidence_modals'
import {Evidence, Finding} from 'src/global_types'
import {RouteComponentProps} from 'react-router-dom'
import {default as Button, ButtonGroup} from 'src/components/button'
import {getFinding} from 'src/services'
import {useWiredData, useModal, renderModals} from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: RouteComponentProps<{slug: string, uuid: string}>) => {
  const {slug, uuid} = props.match.params
  const wiredFinding = useWiredData(React.useCallback(() => getFinding({
    operationSlug: slug,
    findingUuid: uuid,
  }), [slug, uuid]))

  const addRemoveEvidenceModal = useModal<{finding: Finding, initialEvidence: Array<Evidence>}>(modalProps => (
    <ChangeEvidenceOfFindingModal {...modalProps} onChanged={wiredFinding.reload} operationSlug={slug} />
  ))
  const editFindingModal = useModal<{finding: Finding}>(modalProps => (
    <EditFindingModal {...modalProps} onEdited={wiredFinding.reload} operationSlug={slug} />
  ))
  const deleteFindingModal = useModal<{finding: Finding}>(modalProps => (
    <DeleteFindingModal {...modalProps} onDeleted={() => props.history.push(`/operations/${slug}/findings`)} operationSlug={slug} />
  ))
  const editEvidenceModal = useModal<{evidence: Evidence}>(modalProps => (
    <EditEvidenceModal {...modalProps} onEdited={wiredFinding.reload} operationSlug={slug} />
  ))
  const removeEvidenceFromFindingModal = useModal<{evidence: Evidence, finding: Finding}>(modalProps => (
    <RemoveEvidenceFromFindingModal {...modalProps} onRemoved={wiredFinding.reload} operationSlug={slug} />
  ))

  return <>
    {wiredFinding.render(({finding, evidence}) => (
      <div className={cx('root')}>
        <div className={cx('finding-info')}>
          <div className={cx('actions')}>
            <Button small className={cx('left')} icon={require('./back.svg')} onClick={() => props.history.goBack()}>Back</Button>
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
            evidence={evidence}
            actions={{
              'Remove From Finding': evidence => removeEvidenceFromFindingModal.show({evidence, finding}),
              'Edit': evidence => editEvidenceModal.show({evidence}),
            }}
            onQueryUpdate={query => props.history.push(`/operations/${slug}/evidence?q=${encodeURIComponent(query)}`)}
            operationSlug={slug}
            query=""
          />
        </div>
      </div>
    ))}
    {renderModals(addRemoveEvidenceModal, editFindingModal, deleteFindingModal, editEvidenceModal, removeEvidenceFromFindingModal)}
  </>
}
