// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { Evidence } from "src/global_types"

import { useFormField, useModal, renderModals } from 'src/helpers'
import { runServiceWorkerMatrix } from 'src/services'

import { BulletProps, ManagedServiceWorkerChooser } from 'src/components/bullet_chooser'
import Button from 'src/components/button'
import ChallengeModalForm from 'src/components/challenge_modal_form'
import EvidenceChooser from 'src/components/evidence_chooser'
import Modal from 'src/components/modal'
import SettingsSection from 'src/components/settings_section'
import WithLabel from 'src/components/with_label'


const cx = classnames.bind(require('./stylesheet'))

export const BatchRunWorker = (props: {
  operationSlug: string
}) => {
  const [selectedWorkers, setSelectedWorkers] = React.useState<Array<BulletProps>>([])
  const [selectedEvidence, setSelectedEvidence] = React.useState<Array<Evidence>>([])

  const chooseEvidenceModal = useModal<{}>(modalProps => (
    <ChooseEvidenceModal
      initialEvidence={selectedEvidence}
      operationSlug={props.operationSlug}
      onChanged={(list) => setSelectedEvidence(
        list.map(uuid => ({
          uuid,
          description: '',
          operator: {
            firstName: '',
            lastName: '',
            slug: ''
          },
          occurredAt: new Date(),
          tags: [],
          metadata: [],
          contentType: 'image',
          sendImageInfo: false
        }))
      )}
      {...modalProps}
    />
  ))

  const startWOrkersModal = useModal<{}>(modalProps => (
    <StartWorkerModal
      onSubmit={ async() => {
        if (selectedEvidence.length == 0) {
          throw new Error("Some services must be selected")
        }
        runServiceWorkerMatrix({
          operationSlug: props.operationSlug,
          workers: selectedWorkers.map(bp => bp.name),
          evidenceUuids: selectedEvidence.map(bp => bp.uuid),
        })
      }}
      workers={selectedWorkers}
      evidence={selectedEvidence}
      {...modalProps}
    />
  ))

  const startButtonEnabled = (
    selectedWorkers.length > 0 // some workers selected
    && (selectedEvidence.length > 0) // some evidence selected
  )

  const selectedItemsText = `${selectedEvidence.length} items selected`

  return (
    <SettingsSection title="Run Workers">
      <em className={cx('preamble')}>
        You can re-run workers on all, or a certain subset, of evidence for this operation.
        Note that this process may take awhile to complete and may not necessarily produce
        better data than before.
      </em>

      <div className={cx('control-container')}>
        <Area className={cx('worker-control')}>
          <ManagedServiceWorkerChooser
            label='Choose workers'
            operationSlug={props.operationSlug}
            value={selectedWorkers}
            onChange={setSelectedWorkers}
          />
        </Area>

        <Area>
          <WithLabel label='Select Evidence'>
            <div className={cx('multi-item-row')}>
              <Button
                className={cx('choose-button')}
                onClick={() => chooseEvidenceModal.show({})}
              >
                Browse
              </Button>
              <div className={cx('selected-label')}>{selectedItemsText}</div>
            </div>
          </WithLabel>
        </Area>
        <Area startOfRow colspan={2} className={cx('start-button')}>
          <Button
            primary
            disabled={!startButtonEnabled}
            title={startButtonEnabled ? "Start the workers" : "Choose some workers and evidence to start"}
            onClick={() => startWOrkersModal.show({})}
          >
            Start
          </Button>
        </Area>
      </div>

      {renderModals(chooseEvidenceModal, startWOrkersModal)}
    </SettingsSection>
  )
}
export default BatchRunWorker


const ChooseEvidenceModal = (props: {
  initialEvidence: Array<Evidence>,
  onRequestClose: () => void,
  onChanged: (uuid: Array<string>) => void,
  operationSlug: string,
}) => {
  const evidenceField = useFormField<Array<Evidence>>(props.initialEvidence)

  return (
    <Modal title="Search for evidence" onRequestClose={props.onRequestClose}>
      <EvidenceChooser operationSlug={props.operationSlug} {...evidenceField} includeSelectAll/>
      <Button primary className={cx('submit-button')} onClick={() => {
        props.onChanged(evidenceField.value.map(evi => evi.uuid))
        props.onRequestClose()
      }}>Select</Button>
    </Modal>
  )
}

export const StartWorkerModal = (props: {
  workers: Array<BulletProps>,
  evidence: Array<Evidence>
  onRequestClose: () => void,
  onSubmit: () => Promise<void>
}) => {
  const quantityText = props.evidence.length == 1
    ? "1 piece"
    : `${props.evidence.length} pieces`
  const warningText = (
    `This will start workers for ${quantityText} of evidence. ` +
    `Are you sure you want to continue?`
  )

  return (
    <ChallengeModalForm
      modalTitle="Start Services"
      warningText={warningText}
      submitText="Start"
      handleSubmit={props.onSubmit}
      onRequestClose={props.onRequestClose}
    />
  )
}

const Area = (props: {
  startOfRow?: true
  colspan?: 2
  children?: React.ReactNode
  className?: string
}) => {
  const colspan = props.colspan ?? 1
  return (
    <div className={cx(
      props.startOfRow ? 'grid-col-1' : null,
      colspan == 2 ? 'grid-span-2' : null,
      props.className
    )}>
      {props.children}
    </div>
  )
}
