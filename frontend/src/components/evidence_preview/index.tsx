import * as React from 'react'
import classnames from 'classnames/bind'
import { CodeBlockViewer } from '../code_block'
import { HarViewer, isAHar } from '../http_cycle_viewer'
import { SupportedEvidenceType, CodeBlock, EvidenceViewHint, InteractionHint, UrlData } from 'src/global_types'
import { getEvidenceAsCodeblock, getEvidenceAsString, getEvidenceAsUrlData, updateEvidence } from 'src/services/evidence'
import { useWiredData } from 'src/helpers'
import ErrorDisplay from 'src/components/error_display'
import LazyLoadComponent from 'src/components/lazy_load_component'


import TerminalPlayer from 'src/components/terminal_player'

const cx = classnames.bind(require('./stylesheet'))

function getComponent(evidenceType: SupportedEvidenceType) {
  switch (evidenceType) {
    case 'codeblock':
      return EvidenceCodeblock
    case 'image':
      return MemoizedEvidenceImage
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
  useS3Url: boolean,
  onClick?: (event: React.MouseEvent<HTMLDivElement, MouseEvent>) => void,
  imgDataSetter?: (urlData: UrlData | null) => void,
  preSavedS3UrlData?: UrlData,
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
      <LazyLoadComponent><Component {...props} /></LazyLoadComponent>
    </div>
  )
}

type EvidenceProps = {
  operationSlug: string,
  evidenceUuid: string,
  viewHint?: EvidenceViewHint,
  interactionHint?: InteractionHint,
  useS3Url: boolean
  imgDataSetter?: (urlData: UrlData) => void,
  preSavedS3UrlData?: UrlData,
}

const EvidenceCodeblock = (props: EvidenceProps) => {
  const wiredEvidence = useWiredData<CodeBlock>(React.useCallback(() => getEvidenceAsCodeblock({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [props.operationSlug, props.evidenceUuid]))

  return wiredEvidence.render(evi => <CodeBlockViewer value={evi} />)
}

const EvidenceImage = (props: EvidenceProps) => {
  console.log("EvidenceImage props", props?.preSavedS3UrlData?.url)
  let url = `/web/operations/${props.operationSlug}/evidence/${props.evidenceUuid}/media`
  if (props.useS3Url && props.preSavedS3UrlData && new Date(props.preSavedS3UrlData.expirationTime) > new Date()){
    url = props.preSavedS3UrlData.url
  } else if (props.useS3Url) {
    const wiredUrl = useWiredData<UrlData>(React.useCallback(() => getEvidenceAsUrlData({
      operationSlug: props.operationSlug,
      evidenceUuid: props.evidenceUuid,
    }), [props.operationSlug, props.evidenceUuid]))
    wiredUrl.expose(s3url => {
      props.imgDataSetter && props.imgDataSetter(s3url)
      url = s3url.url 
    })
  }
  return <img src={url} />
}

const sameURL = (prevProps: EvidenceProps, nextProps: EvidenceProps) => {
  console.log("sameURL props", prevProps?.preSavedS3UrlData?.url, nextProps?.preSavedS3UrlData?.url)
  return prevProps?.preSavedS3UrlData?.url === nextProps?.preSavedS3UrlData?.url;
};

const MemoizedEvidenceImage = React.memo(EvidenceImage, sameURL);

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

  return wiredEvidence.render(evi => <TerminalPlayer content={evi} playerUUID={props.evidenceUuid} onTerminalScriptUpdated={updateContent} />)
}

const EvidenceHttpCycle = (props: EvidenceProps) => {
  const wiredEvidence = useWiredData<string>(React.useCallback(() => getEvidenceAsString({
    operationSlug: props.operationSlug,
    evidenceUuid: props.evidenceUuid,
  }), [props.operationSlug, props.evidenceUuid]))

  return wiredEvidence.render(evi => {
    try {
      const log = JSON.parse(evi)
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
