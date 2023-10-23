// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { CodeBlockViewer } from '../code_block'
import { HarViewer, isAHar } from '../http_cycle_viewer'
import { SupportedEvidenceType, CodeBlock, EvidenceViewHint, InteractionHint, ImageInfo } from 'src/global_types'
import { getEvidenceAsCodeblock, getEvidenceAsString, getEvidenceUrl, updateEvidence } from 'src/services/evidence'
import { useWiredData } from 'src/helpers'
import ErrorDisplay from 'src/components/error_display'

import TerminalPlayer from 'src/components/terminal_player'

const cx = classnames.bind(require('./stylesheet'))

function getComponent(evidenceType: SupportedEvidenceType) {
  switch (evidenceType) {
    case 'codeblock':
      return EvidenceCodeblock
    case 'image':
      return EvidenceImage
    case 'terminal-recording':
      return EvidenceTerminalRecording
    case 'http-request-cycle':
      return EvidenceHttpCycle
    case 'event':
      return EvidenceEvent
    case 'none':
    default:
      return null
  }
}

export default (props: {
  operationSlug: string,
  evidenceUuid: string,
  contentType: SupportedEvidenceType,
  viewHint?: EvidenceViewHint,
  interactionHint?: InteractionHint,
  className?: string,
  fitToContainer?: boolean,
  streamImage: boolean,
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
  viewHint?: EvidenceViewHint,
  interactionHint?: InteractionHint,
  streamImage: boolean
}

const EvidenceCodeblock = (props: EvidenceProps) => {
  const wiredEvidence = useWiredData<CodeBlock>(React.useCallback(() => getEvidenceAsCodeblock({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [props.operationSlug, props.evidenceUuid]))

  return wiredEvidence.render(evi => {
    console.log("evi in codeblock", evi)
  return <CodeBlockViewer value={evi} />})
}

// TODO TN - only send extra data if image?
const EvidenceImage = (props: EvidenceProps) => {
  if (props.streamImage) {
    console.log("non-s3 devlelopment")
    const fullUrl = `/web/operations/${props.operationSlug}/evidence/${props.evidenceUuid}/media`
    return <img src={fullUrl} />
  } else {
    console.log("using s3 get evidenceUrl")
    const wiredImageInfo = useWiredData<ImageInfo>(React.useCallback(() => getEvidenceUrl({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidenceUuid,
    }), [props.operationSlug, props.evidenceUuid]))
  
    // TODO TN rename wiredimageinfor
    console.log("about to see wiredImageinfo")
    return wiredImageInfo.render(url => {
      console.log("___url", url)
      // console.log(url == null, url == undefined, url == "")
      // if (url != ""){
      //   console.log("___url JSON", JSON.parse(url))
      // } else {
      //   url = "https://upload.wikimedia.org/wikipedia/commons/9/97/The_Earth_seen_from_Apollo_17.jpg" 
      // }
    return <img src={url.url} />
  })
  }
}

const EvidenceEvent = (_props: EvidenceProps) => {
  return <div className={cx('event')}></div>
}

const EvidenceTerminalRecording = (props: EvidenceProps) => {
  const wiredEvidence = useWiredData<string>(React.useCallback(() => getEvidenceAsString({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [props.operationSlug, props.evidenceUuid]))

  const updateContent = (content: Blob): Promise<void> => updateEvidence({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
    updatedContent: content,
  })

  return wiredEvidence.render(evi => { 
    console.log("evi term recording", evi); 
    return <TerminalPlayer content={evi} playerUUID={props.evidenceUuid} onTerminalScriptUpdated={updateContent} />})
}

const EvidenceHttpCycle = (props: EvidenceProps) => {
  const wiredEvidence = useWiredData<string>(React.useCallback(() => getEvidenceAsString({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [props.operationSlug, props.evidenceUuid]))

  return wiredEvidence.render(evi => {
    try {
      // console.log("evi http cycle", evi)
      const log = JSON.parse(evi)
      // console.log("log", log)
      if (isAHar(log)) {
        const isActive = props.interactionHint == 'inactive' ? {disableKeyHandler : true} : {}
        return <HarViewer log={log} viewHint={props.viewHint} {...isActive} />
      }
      return <ErrorDisplay title="Corrupted HAR file" err={new Error("unsupported format")} />
    }
    catch (err) {
      return <ErrorDisplay title="Corrupted HAR file" err={err}/>
    }
  })
}
