// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Chooser from 'src/components/chooser'
import EvidencePreview from 'src/components/evidence_preview'
import Lightbox from 'src/components/lightbox'
import classnames from 'classnames/bind'
import { Evidence } from 'src/global_types'

import { getEvidenceList } from 'src/services'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  disabled?: boolean,
  onChange: (e: Array<Evidence>) => void,
  operationSlug: string,
  value: Array<Evidence>,
  includeSelectAll?: true
}) => {
  const fetchEvidence = React.useCallback((query: string) => getEvidenceList({ operationSlug: props.operationSlug, query }), [props.operationSlug])
  return (
    <Chooser
      {...props}
      placeholder="Filter Evidence"
      fetch={ fetchEvidence }
      renderRow={evi => <EvidenceRow operationSlug={props.operationSlug} evidence={evi} />}
    />
  )
}

const EvidenceRow = (props: {
  evidence: Evidence,
  operationSlug: string,
}) => {
  const [lightboxOpen, setLightboxOpen] = React.useState(false)

  return <>
    <div className={cx('media')} >
      <EvidencePreview
        operationSlug={props.operationSlug}
        evidenceUuid={props.evidence.uuid}
        contentType={props.evidence.contentType}
        onClick={(e) => { e.stopPropagation(); setLightboxOpen(true)} }
        fitToContainer
        viewHint="small"
        interactionHint="inactive"
      />
    </div>
    <div className={cx('description')} onClick={e => e.stopPropagation()}>{props.evidence.description}</div>
    <div onClick={e => e.stopPropagation()}>
      <Lightbox isOpen={lightboxOpen} onRequestClose={() => setLightboxOpen(false)}>
        <EvidencePreview
          operationSlug={props.operationSlug}
          evidenceUuid={props.evidence.uuid}
          contentType={props.evidence.contentType}
          viewHint="large"
          interactionHint="active"
        />
      </Lightbox>
    </div>
  </>
}
