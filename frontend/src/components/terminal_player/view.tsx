import * as React from 'react'
import classnames from 'classnames/bind'
import { PositionChangeEventBody, RateChangeEventBody, ExpandedTerminalEvent } from './types'
import { format } from 'date-fns'

import { ClickPopover } from 'src/components/popover'
import { default as Menu, MenuItem } from 'src/components/menu'
import { default as Button, ButtonGroup } from 'src/components/button'
import { CreateBookmarkModal } from './modals'
import { useElementRect } from 'src/helpers/use_element_rect'
import {
  default as TerminalPlayer,
  EventTypeFrameAdvance, EventTypeRateChange, EventTypeDesiredRateChange,
} from './player'

import "xterm/css/xterm.css"
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  content: string
  playerUUID: string
  onTerminalScriptUpdated: (b: Blob) => Promise<void>
}) => {
  const rootRef = React.useRef<HTMLDivElement | null>(null)
  const termRef = React.useRef<HTMLDivElement | null>(null)
  const termPlayer = React.useRef<TerminalPlayer | null>(null)

  const [termStyle, setTermStyle] = React.useState<React.CSSProperties>({})
  const [wrapperStyle, setWrapperStyle] = React.useState<React.CSSProperties>({})

  const rootRect = useElementRect(rootRef)
  React.useEffect(() => {
    if (rootRect == null || termRef.current == null) {
      return
    }
    const wFactor = rootRect.width / termRef.current.clientWidth
    // 70 here refers to the control padding area, defined in css.
    const hFactor = (rootRect.height - 70) / termRef.current.clientHeight
    let scale = Math.min(wFactor, hFactor)

    if (isNaN(scale) || scale == Infinity) {
      scale = 1
    }

    setTermStyle({
      transform: `translate(${(rootRect.width - (scale * termRef.current.clientWidth)) / 2}px) scale(${scale})`,
      transformOrigin: 'top left',
    })
    setWrapperStyle({ height: termRef.current.clientHeight * scale })
  }, [rootRect, termRef])

  React.useEffect(() => {
    if (termPlayer.current != null || termRef.current == null) {
      return
    }

    const player = new TerminalPlayer(props.content)
    termPlayer.current = player
    player.init(termRef.current)

    return () => {
      player.cleanup()
      termPlayer.current = null
    }
  }, [props.playerUUID, props.content])

  return <>
    <div className={cx('root')} ref={rootRef} onClick={e => e.stopPropagation()}>
      <div style={wrapperStyle}>
        <div ref={termRef} className={cx("terminal")} style={termStyle} />
      </div>
      {termPlayer.current != null &&
        <PlaybackControlLogic player={termPlayer.current} onTerminalScriptUpdated={props.onTerminalScriptUpdated} />
      }
    </div>
  </>
}

const PlaybackControlLogic = (props: {
  player: TerminalPlayer
  onTerminalScriptUpdated: (b: Blob) => Promise<void>
}) => {
  const [model, setModel] = React.useState<Array<string> | null>(null)
  const [bookmarkInsertIndex, setBookmarkInserIndex] = React.useState(0)
  const [playing, setPlaying] = React.useState(false)
  const [desiredRate, setDesiredRate] = React.useState(1)
  const [playbackTime, setPlaybackTime] = React.useState(0)
  const [completed, setCompleted] = React.useState(0)

  const termPlayer = props.player

  const handlerPairs: Array<[string, (a: PositionChangeEventBody | RateChangeEventBody) => void]> =
  React.useMemo(() => [
      [EventTypeFrameAdvance, (evt: PositionChangeEventBody) => { setPlaybackTime(evt.elapsedTime); setCompleted(evt.playbackPosition) }],
      [EventTypeRateChange, (evt: RateChangeEventBody) => setPlaying(evt.newRate != 0)],
      [EventTypeDesiredRateChange, (evt: RateChangeEventBody) => setDesiredRate(evt.newRate)]
    ], [])

  React.useEffect(() => {
    handlerPairs.forEach(p => termPlayer.on(...p))
    return () => { handlerPairs.forEach(p => termPlayer.removeListener(...p)) }
  }, [handlerPairs, termPlayer])

  return <>
    <PlaybackControls
      disabled={termPlayer.getError() != null}
      playButtonText={playing ? "⏸️" : "▶️"}
      desiredPlaybackRate={desiredRate}
      playbackTime={formatAsDuration(playbackTime)}
      supportedPlaybackRates={[.5, 1, 2, 4, 8, 16, 32, 64]}
      setRate={termPlayer.setRate}
      playPauseToggle={e => { e.preventDefault(); (playing ? termPlayer.pause : termPlayer.play)() }}
      resetStream={e => { e.preventDefault(); termPlayer.reset() }}
      openBookmarkModal={e => { e.preventDefault(); setModel(termPlayer.getBookmarkAtCursor()); setBookmarkInserIndex(termPlayer.getCurrentIndex()) }}
    >
      <ProgressBar
        percentCompleted={completed}
        bookmarks={termPlayer.getBookmarks()}
        streamDuration={termPlayer.getDuration()}
        jumpToEventIndex={termPlayer.jumpToEventIndex}
        jumpToPosition={termPlayer.jumpToPosition}
        closestEventTime={termPlayer.nearestEventTime}
      />
    </PlaybackControls>

    {model != null &&
      <CreateBookmarkModal
        onRequestClose={() => setModel(null)}
        updateEncoding={(desc: string) => {
          if (desc !== "") {
            termPlayer.addBookmark(bookmarkInsertIndex, desc)
          }
          else {
            termPlayer.removeBookmarkAtCursor()
          }
          return termPlayer.export()
        }}
        handleSubmit={props.onTerminalScriptUpdated}
        onCancel={() => { }}
        initialDescription={(() => {
          return model.join("\n")
        })()}
      />
    }
  </>
}

const PlaybackControls = (props: {
  disabled: boolean
  playButtonText: string
  desiredPlaybackRate: number
  playbackTime: string
  supportedPlaybackRates: Array<number>
  setRate: (newRate: number) => void
  playPauseToggle: (e: React.MouseEvent<Element, MouseEvent>) => void
  resetStream: (e: React.MouseEvent<Element, MouseEvent>) => void
  openBookmarkModal: (e: React.MouseEvent<Element, MouseEvent>) => void
  children: React.ReactNode,
}) => <>
    <div className={cx('controls')}>
      <ButtonGroup className={cx("playback-controls", "left-cluster")}>
        <Button disabled={props.disabled} onClick={props.resetStream}>⏮</Button>
        <Button disabled={props.disabled} className={cx("play-button")} onClick={props.playPauseToggle}>{props.playButtonText}</Button>
      </ButtonGroup>
      <div className={cx("timestamp")} >{props.playbackTime}</div>
      {props.children}

      <ButtonPopover
        className={'rate-menu'}
        disabled={props.disabled}
        selectedLabel={props.desiredPlaybackRate + ' ✕'}
        labels={props.supportedPlaybackRates.map(x => x + ' ✕')
        }
        valueSelected={(label) => label === "Normal" ? 1 : props.setRate(parseFloat(label))}
      />

      <ButtonGroup className={cx("playback-controls", "right-cluster")}>
        <Button disabled={props.disabled} icon={require('./book-mark.svg')} onClick={props.openBookmarkModal} />
      </ButtonGroup>
    </div>
  </>

const ProgressBar = (props: {
  percentCompleted: number
  bookmarks: Array<ExpandedTerminalEvent>
  streamDuration: number
  jumpToEventIndex: (idx: number) => void
  jumpToPosition: (pct: number) => void
  closestEventTime: (pct: number) => number
}) => {
  const progressRef = React.useRef<HTMLDivElement | null>(null)
  const [hoverText, setHoverText] = React.useState("")

  const percentIntoProgressBar = (clientX: number): number => {
    if (progressRef.current) {
      const containerRect = progressRef.current.getBoundingClientRect()
      return (clientX - containerRect.left) / containerRect.width
    }
    return 0
  }

  const progressBarEvents = {
    onClick: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
      props.jumpToPosition(percentIntoProgressBar(e.clientX))
    },
    onMouseMove: (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
      const percentComplete = percentIntoProgressBar(e.clientX)
      const eventTime = props.closestEventTime(percentComplete)
      setHoverText(`Jump To: ${format(eventTime * 1000, "MMM dd, yyyy '@' HH:mm")}`)
    },
    title: hoverText,
  }

  return <>
    <div ref={progressRef} className={cx("progressbar-container")} {...progressBarEvents} >
      <div className={cx("progressbar")} style={{ width: `${props.percentCompleted * 100}%` }} />
      {props.bookmarks.map((mark) => {
        const position = mark.totalDelay / props.streamDuration
        return (
          <img
            src={require('./marker.svg')}
            key={position}
            title={`* ${mark.bookmarks.join("\n* ")}`}
            className={cx("bookmark")}
            onClick={(e) => {
              e.stopPropagation()
              props.jumpToEventIndex(mark.eventIndex)
            }}
            style={{ left: `${position * 100}%` }}
          />)
      }
      )}
    </div>
  </>
}

const formatAsDuration = (ms: number): string => {
  const fullSeconds = Math.trunc(ms / 1000)
  const fullMinutes = Math.trunc(fullSeconds / 60)
  const hours = Math.trunc(fullMinutes / 60)
  const lpad = (s: number) => (s < 10 ? "0" : "") + s

  return `${lpad(hours)}:${lpad(fullMinutes % 60)}:${lpad(fullSeconds % 60)}`
}

const ButtonPopover = (props: {
  labels: Array<string>
  disabled?: boolean
  selectedLabel: string
  className?: string
  valueSelected: (label: string) => void
}) => {
  return <>
    <ClickPopover closeOnContentClick content={
      <Menu >
        {props.labels.map(label =>
          <MenuItem key={label} onClick={(e) => { e.preventDefault(); props.valueSelected(label) }}>
            {`${label} ${props.selectedLabel === label ? " ✔" : ""}`}
          </MenuItem>
        )}
      </Menu>
    }>
      <Button disabled={props.disabled} className={cx(props.className)} onClick={e => e.preventDefault()}>{props.selectedLabel}</Button>
    </ClickPopover>
  </>
}
