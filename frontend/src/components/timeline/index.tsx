import * as React from 'react'
import EvidencePreview from 'src/components/evidence_preview'
import Lightbox from 'src/components/lightbox'
import MarkdownRenderer from 'src/components/markdown_renderer'
import TagList from 'src/components/tag_list'
import classnames from 'classnames/bind'
import Help from 'src/components/help'
import { ClickPopover } from 'src/components/popover'
import { Tag, Evidence } from 'src/global_types'
import { addTagToQuery, addOperatorToQuery } from 'src/helpers'
import { default as Button, ButtonGroup } from 'src/components/button'
import { CopyTextButton } from 'src/components/text_copiers'
import { format } from 'date-fns'
import { default as Menu, MenuItem } from 'src/components/menu'
import EvidencesContextProvider from 'src/contexts/evidences_context'

const cx = classnames.bind(require('./stylesheet'))

type Action = {
  label: string
  act: (evi: Evidence) => void
  canAct?: (evi: Evidence) => { disabled: boolean, title?: string }
}
type Actions = Array<Action>

export default (props: {
  actions: Actions,
  extraActions?: Actions,
  evidence: Array<Evidence>,
  onQueryUpdate: (q: string) => void,
  operationSlug: string,
  query: string,
  scrollToUuid?: string,
}) => {
  const rootRef = React.useRef<HTMLDivElement | null>(null)
  const lightboxRef = React.useRef<HTMLDivElement | null>(null)

  const [activeChildIndex, setActiveChildIndex] = React.useState<number>(0)
  const [quicklookVisible, setQuicklookVisible] = React.useState<boolean>(false)
  // const [currImageData, setCurrImageData] = React.useState<UrlData| null>(null)

  const onKeyDown = (e: KeyboardEvent) => {
    // Only handle keystrokes if nothing is focused (target is body)
    // or if the focused element belongs to this component (child of root or lightbox)
    if (e.target == null) return
    if (e.target !== document.body && !elementInRef(e.target as HTMLElement, [rootRef, lightboxRef])) return

    const children = refDivChildren(rootRef)
    let newActiveChildIndex = activeChildIndex
    switch (e.key) {
      case 'ArrowDown': case 'ArrowRight': case 'j':
        newActiveChildIndex = Math.min(activeChildIndex + 1, children.length - 1)
        break
      case 'ArrowUp': case 'ArrowLeft': case 'k':
        newActiveChildIndex = Math.max(activeChildIndex - 1, 0)
        break
      case 'g':
        newActiveChildIndex = 0
        break
      case 'G':
        newActiveChildIndex = children.length - 1
        break
      case ' ':
        setQuicklookVisible(!quicklookVisible)
        break
      case 'Enter':
        setQuicklookVisible(true)
        break
      case 'Escape':
        setQuicklookVisible(false)
        break
      default:
        return
    }
    e.preventDefault()
    scrollRef(rootRef, children[newActiveChildIndex].offsetTop)
    setActiveChildIndex(newActiveChildIndex)
  }

  React.useEffect(() => {
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  })

  const activeEvidence = props.evidence[activeChildIndex]
  if (activeEvidence == null) return null

  return (
    <EvidencesContextProvider>
      <div className={cx('root')} ref={rootRef}>
        {props.evidence.map((evi, idx) => {
          const active = activeChildIndex === idx
          return (
              <TimelineRow
                {...props}
                focusUuid={props.scrollToUuid}
                active={active}
                evidence={evi}
                key={evi.uuid}
                onPreviewClick={() => { setActiveChildIndex(idx); setQuicklookVisible(true) }}
                onClick={() => setActiveChildIndex(idx)}
              />
          )
        })}
        <Help className={cx('help')}
          preamble="Review and Edit the accumulated evidence for this operation"
          shortcuts={KeyboardShortcuts}
        />
      </div>
      <Lightbox canUseFitToggle={activeEvidence.contentType == "image"}
        isOpen={quicklookVisible} onRequestClose={() => setQuicklookVisible(false)}>
        <div ref={lightboxRef}>
            <EvidencePreview
              operationSlug={props.operationSlug}
              evidenceUuid={activeEvidence.uuid}
              contentType={activeEvidence.contentType}
              useS3Url={activeEvidence.sendUrl}
              viewHint="large"
              interactionHint="active"
            />
        </div>
      </Lightbox>
    </EvidencesContextProvider>
  )
}

const TimelineRow = (props: {
  active: boolean,
  actions: Actions,
  extraActions?: Actions,
  evidence: Evidence,
  onQueryUpdate: (q: string) => void,
  operationSlug: string,
  query: string,
  focusUuid?: string,
  onPreviewClick: () => void,
  onClick: () => void
}) => {
  const self = React.useRef<null | HTMLDivElement>(null)

  React.useEffect(() => {
    if (self.current != null && props.evidence.uuid == props.focusUuid) {
      self.current.scrollIntoView()
    }
  }, [self, props.focusUuid, props.evidence.uuid])

  const onTagClick = (t: Tag) => {
    if (!t.id) {
      return
    }
    props.onQueryUpdate(addTagToQuery(props.query, t.name))
  }

  const onOperatorClick = () => {
    props.onQueryUpdate(addOperatorToQuery(props.query, props.evidence.operator.slug))
  }

  const permalink = `${window.location.origin}/operations/${props.operationSlug}/evidence/${props.evidence.uuid}`

  return (
    <div ref={self} className={cx('timeline-row', { active: props.active })} onClick={props.onClick}>
      <div className={cx('left')}>
        <EvidencePreview
          fitToContainer
          onClick={props.onPreviewClick}
          operationSlug={props.operationSlug}
          evidenceUuid={props.evidence.uuid}
          contentType={props.evidence.contentType}
          useS3Url={props.evidence.sendUrl}
          viewHint="medium"
          interactionHint="inactive"
        />
      </div>
      <div className={cx('right')}>
        <div>{format(props.evidence.adjustedAt ?? props.evidence.occurredAt, "MMMM do, yyyy 'at' HH:mm:ss")}</div>
        <a href="#" onClick={onOperatorClick}>
          {props.evidence.operator.firstName} {props.evidence.operator.lastName}
        </a>
        <TagList
          tags={
            props.evidence.adjustedAt
            ? [
                ...props.evidence.tags,
                { id: 0, name: 'Adjusted Timestamp', colorName: 'yellow' }
              ]
            : props.evidence.tags
          }
          onTagClick={onTagClick}
        />
        <ButtonGroup>
          {
            props.actions.map(action => (
              <Button
                small
                key={action.label}
                onClick={() => action.act(props.evidence)}
                {...action.canAct?.(props.evidence)}
              >
                {action.label}
              </Button>
            ))
          }
          <CopyTextButton small textToCopy={permalink}>Copy Permalink</CopyTextButton>
          {renderExtraActions(props.evidence, props.extraActions)}
        </ButtonGroup>
        <MarkdownRenderer className={cx('description')}>{props.evidence.description}</MarkdownRenderer>
      </div>
    </div>
  )
}

const renderExtraActions = (evidence: Evidence, extraActions?: Actions) => {
  if (!extraActions) {
    return null
  }

  const menuItems = extraActions.map(action => (
    <MenuItem
      key={action.label}
      onClick={() => action.act(evidence)}
      {...action.canAct?.(evidence)}
    >
      {action.label}
    </MenuItem>
  ))

  return (
    <ClickPopover className={cx('popover')} closeOnContentClick content={<Menu>{menuItems}</Menu>}>
      <Button small className={cx('arrow')} />
    </ClickPopover>
  )
}

// Returns true if the element el is a child of any of the supplied react refs
function elementInRef(el: HTMLElement, refs: Array<React.MutableRefObject<HTMLElement | null>>): boolean {
  const targetEls = refs.map(el => el.current)
  while (!targetEls.includes(el)) {
    if (el === document.body) return false
    if (el.parentElement == null) return false
    el = el.parentElement
  }
  return true
}

// Returns an array of direct children that are divs for the passed react ref
function refDivChildren(ref: React.MutableRefObject<HTMLDivElement | null>): Array<HTMLDivElement> {
  if (ref.current == null) return []
  // @ts-ignore - typescript is unable to determine that children is an array of HTMLDivElements
  return Array.from(ref.current.children).filter(el => el instanceof HTMLDivElement)
}

// Scroll a react ref to the given position in pixels
function scrollRef(ref: React.MutableRefObject<HTMLDivElement | null>, scrollTop: number) {
  if (ref.current == null) return
  ref.current.scrollTop = scrollTop
}

export const KeyboardShortcuts = [
  { keys: ["ArrowDown", "ArrowRight", "j"], description: "Move to the next evidence" },
  { keys: ["ArrowUp", "ArrowLeft", "k"], description: "Move to the previous evidence" },
  { keys: ["g"], description: "Move to the top of the evidence list" },
  { keys: ["G"], description: "Move to the bottom of the evidence list" },
  { keys: ["Enter"], description: "Open evidence large view" },
  { keys: ["Escape"], description: "Close evidence large view" },
  { keys: [" "], description: "Toggle evidence large view" },
  { keys: ["z", "Z"], description: "Toggle Best Fit vs Standard views" },
]
