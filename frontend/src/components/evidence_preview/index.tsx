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
    const url = "https://my-ashirt-data-1234.s3.us-west-2.amazonaws.com/e308c2de-31a0-4594-bb73-9a7806f045a9?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIAZXQTDXNEMGF66V2A%2F20231023%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20231023T164535Z&X-Amz-Expires=86400&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEAkaCXVzLXdlc3QtMiJIMEYCIQCFKt5paW5zWqfPoFvDfiNurU%2B%2Byk9O6E5ZAa%2FwGy6iqgIhAIbGIsmwqzOaP2Z46Dg4prgeiK2ddQ7Y8peAbKpNCySnKuMDCDIQABoMNjY4OTgxMjQ2NzkyIgyU779DVnK5iDPPqysqwAONNiWGhDOo3hjFedga6r6kWslJW1edWYDjW00%2BBtrnH0lRfHcf0ZROmSfJd%2FfZh%2Fec3YCTK47dNpsJz0CHqpQT7Vftp9FnAOQfbBL4lK67r5ZZEwsVYvybRbAFXaVKpY2rSm8W5%2F6QmfAPM4ZQne4k7EyC8v0f2k6XfGeubvmHXMOoLAhJcETuTZBZQLWFx3ZlNqeNyRNfu%2F7279w9Q%2FTqhwAe3l6pyJH%2FJVodGY7lUcDW3rIdHQWhI1rwC4n9AU%2B9afy%2B3IeUWUMCaNfACOT4wpQAAOuxvyEPzLPhWGGLy8pPAj80kgR3jtDM2v0XmxGyovsei1PEGz5zy5zzjHtZnC0X14j7%2B1UzacHGQtEmH9HdTIUIif7R8LWnFEIO46vNZ8np9ZPUycS%2BlaXwBv9EcMR0IRDGzyOLpjIE%2BdQAED6OzPOhqe9hAVYOKLtADEm6TD8puBuVn2aEQwTUX5EPXoF8dlq5m5regNrcgB9rPMS6xY%2BDrdifgGE0fnpX9%2FBAP0YbMd9fZLszaGf6iABg%2F5spP2HKCid0JL5SnOzCzBQOjTTpI9kbxLCEiMl%2Ba62XG0xGQXoTdjnc7DTlMDCtMMjF2qkGOqQBOPPZBqSwo7L7%2Fm5hlg%2FiUb45kVNmGgW4DmiKuA7zrnwTdP6MPrxxBy66y5clN5HvD7ZXb0fs8eon246r3wZ6AAOTrZtyyymKGdAB1Xk5bekznf5en8p%2FaIXoNAC0IJGxT3OMOyFJvHlZi1B1KE%2FcBLieDQLSYDlfz%2FH5l5bNYoeH4QRrYMOGyhzSDQhLLBdr4LzxIrzMpsi3CpiYHYi5amo2%2Bl0%3D&X-Amz-SignedHeaders=host&response-content-type=image%2Fjpeg&X-Amz-Signature=cd0bb641d91a7d0fb38f0adfc31405d1eb36628f2e320a83cf9b3369fc5ff471"
    return <img src={url} />
  } else {
    console.log("using s3 get evidenceUrl")
    const wiredImageInfo = useWiredData<ImageInfo>(React.useCallback(() => getEvidenceUrl({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidenceUuid,
    }), [props.operationSlug, props.evidenceUuid]))
  
    // TODO TN rename wiredimageinfor
    console.log("about to see wiredImageinfo")
    return wiredImageInfo.render(url => {
      // const parsedJSON = JSON.parse(url)
      console.log("typeof ___url", typeof url)
      console.log("typeof ___Parsedurl", typeof url)
      console.log("___url.yrk", url.url)
      // console.log(url == null, url == undefined, url == "")      
      // if (url.url != ""){
      //   console.log("___url JSON", url.url)
      // } else {
      //   url.url = "https://upload.wikimedia.org/wikipedia/commons/9/97/The_Earth_seen_from_Apollo_17.jpg" 
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
