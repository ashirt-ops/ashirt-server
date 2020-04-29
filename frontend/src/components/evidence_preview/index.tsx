// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { SupportedEvidenceType, CodeBlock } from 'src/global_types'
import { useDataSource, getEvidenceAsCodeblock, getEvidenceAsString, updateEvidence } from 'src/services'
import { useWiredData } from 'src/helpers'

import TerminalPlayer from 'src/components/terminal_player'
import { CodeBlockViewer } from '../code_block'

const cx = classnames.bind(require('./stylesheet'))

function getComponent(evidenceType: SupportedEvidenceType) {
  switch (evidenceType) {
    case 'codeblock':
      return EvidenceCodeblock
    case 'image':
      return EvidenceImage
    case 'terminal-recording':
      return EvidenceTerminalRecording
    case 'none':
    default:
      return null
  }
}

export default (props: {
  operationSlug: string,
  evidenceUuid: string,
  contentType: SupportedEvidenceType,
  className?: string,
  fitToContainer?: boolean,
  onClick?: (event: React.MouseEvent<HTMLDivElement, MouseEvent>) => void,
}) => {
  const Component = getComponent(props.contentType)
  if (Component == null) return null

  const className = cx(
    'root',
    props.className,
    props.contentType,
    props.fitToContainer ? 'fit' : 'full',
    { clickable: props.onClick },
  )

  return (
    <div className={className} onClick={props.onClick}>
      <Component {...props} />
    </div>
  )
}

type EvidenceProps = {
  operationSlug: string,
  evidenceUuid: string,
}

const EvidenceCodeblock = (props: EvidenceProps) => {
  const ds = useDataSource()
  const wiredEvidence = useWiredData<CodeBlock>(React.useCallback(() => getEvidenceAsCodeblock(ds, {
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [ds, props.operationSlug, props.evidenceUuid]))
  React.useEffect(wiredEvidence.reload, [props.evidenceUuid])

  return wiredEvidence.render(evi => <CodeBlockViewer value={evi} />)
}

const EvidenceImage = (props: EvidenceProps) => {
  const fullUrl = `/web/operations/${props.operationSlug}/evidence/${props.evidenceUuid}/media`
  return <img src={fullUrl} />
}

const EvidenceTerminalRecording = (props: EvidenceProps) => {
  const ds = useDataSource()
  const wiredEvidence = useWiredData<string>(React.useCallback(() => getEvidenceAsString(ds, {
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [ds, props.operationSlug, props.evidenceUuid]))
  React.useEffect(wiredEvidence.reload, [props.evidenceUuid])

  const updateContent = (content: Blob): Promise<void> => updateEvidence(ds, {
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
    updatedContent: content,
  })

  return wiredEvidence.render(evi => <TerminalPlayer content={evi} playerUUID={props.evidenceUuid} onTerminalScriptUpdated={updateContent} />)
}
