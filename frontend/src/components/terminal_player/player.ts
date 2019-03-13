import { Terminal } from "xterm"
import {
  TerminalEvent, TerminalRecordingHeader,
  TerminalRecordingData, ExpandedTerminalEvent,
  PositionChangeEventBody
} from './types'

import {EventEmitter} from 'events'

export const EventTypeFrameAdvance = 'frame advance'
export const EventTypeHeadJump = 'head jump'
export const EventTypeRateChange = 'rate change'
export const EventTypeDesiredRateChange = 'desired rate change'

export default class extends EventEmitter {
  terminal: Terminal
  running: boolean
  activeTimeout: NodeJS.Timeout
  currentIndex: number
  content: TerminalRecordingData
  playbackRate: number

  desiredPlaybackRate: number
  currentPlaybackRate: number

  duration: number
  minRate = 0.5
  maxRate = 64
  maxFrameDelay = 1600

  // advanceFrameOnPlay is present for the scenario where user jumps to a frame while paused.
  // We draw to that frame selection, but don't advance frames unless playing.
  // This flag allows us to be associated with the frame we just played, instead of the frame we are just about to play
  advanceFrameOnPlay = false

  constructor(
    content: string
  ) {
    super()
    this.content = parseTerminalRecording(content)
    this.terminal = new Terminal({
      cols: Math.max(this.content.header.width, 80),
      rows: Math.max(this.content.header.height, 30),
      disableStdin: true,
      fontFamily: "monospace",
      bellStyle: "none",
      scrollback: 0, // scrollback presents some odd UI, so disabling it for now.
    })
    this.currentIndex = 0
    this.playbackRate = 0
    this.duration = this.content.events[this.content.events.length - 1].totalDelay
    this.currentPlaybackRate = 0
    this.setDesiredRate(1)

    const parseError = this.getError()
    if (parseError != null) {
      this.terminal.writeln("Unable to play content. Error encountered:")
      this.terminal.writeln(parseError)
    }
  }

  init(el: HTMLDivElement) {
    this.running = true
    this.terminal.open(el)
  }

  cleanup() {
    this.running = false
    clearTimeout(this.activeTimeout)
    this.terminal.dispose()
  }

  private eventsUntilIndex(i: number): string {
    return this.content.events.slice(0, i).reduce((acc, cur) => acc + cur.eventContent , "")
  }

  private makeFrameAdvanceEvent(eventIndex: number = this.currentIndex): PositionChangeEventBody {
    const thisEvent = this.getEvent(eventIndex)
    return {
      elapsedTime: thisEvent.totalDelay,
      playbackPosition: thisEvent.totalDelay / this.getDuration(),
      index: eventIndex,
      terminalTime: this.eventTimeToRealTime(eventIndex)
    }
  }

  private eventTimeToRealTime(eventIndex: number): number {
    const thisEvent = this.getEvent(eventIndex)
    return this.content.header.timestamp + Math.trunc(thisEvent.eventTime / 1000)
  }

  private writeEvent(eventIndex: number, flush = false) {
    if (flush) {
      this.terminal.reset()
      if (eventIndex > 0) {
        this.writeToTerm(this.eventsUntilIndex(eventIndex))
      }
    }
    this.writeToTerm(this.content.events[eventIndex].eventContent)
    this.emit(EventTypeFrameAdvance, this.makeFrameAdvanceEvent(eventIndex))
  }

  private writeToTerm = (msg: string) => {
    try {
      this.terminal.write(msg)
    } catch { }
  }

  private waitForFrame() {
    const events = this.content.events
    const delay = this.currentIndex == 0 ? 0 : events[this.currentIndex - 1].frameDuration / this.getRate()
    this.activeTimeout = setTimeout(() => {
      this.writeEvent(this.currentIndex)
      this.currentIndex++
      if (this.currentIndex < events.length - 1) {
        this.waitForFrame()
      }
    }, delay)
  }

  private updatePlaybackRate(newRate: number) {
    const oldRate = this.currentPlaybackRate
    this.currentPlaybackRate = newRate
    this.emit(EventTypeRateChange, {oldRate, newRate})
  }

  private getEvent(index: number): ExpandedTerminalEvent {
    const events = this.content.events
    return events[clamp(index, 0, events.length - 1)]
  }

  private findClosestEventViaPosition = (position: number) => {
    position = clamp(position, 0, 1)
    const desiredTimeOffset = this.getDuration() * position
    return findClosestEvent(desiredTimeOffset, this.content.events, 0, this.content.events.length - 1)
  }


  //public interface methods

  //playback controls

  play = () => {
    if (!this.isPlaying()) {
      if (this.advanceFrameOnPlay) {
        this.currentIndex++
        this.advanceFrameOnPlay = false
      }
      this.updatePlaybackRate(this.desiredPlaybackRate)
      this.waitForFrame()
    }
  }

  pause = () => {
    if (this.isPlaying()) {
      clearTimeout(this.activeTimeout)
    }
    this.updatePlaybackRate(0)
  }

  reset = () => {
    this.pause()
    this.jumpToIndex(0)
  }

  //position will be enforced to be between 0 and 1
  jumpToPosition = (position: number) => {
    const newIndex = this.findClosestEventViaPosition(position)
    this.jumpToIndex(newIndex)
  }

  jumpToEventIndex = (evtIndex: number) => {
    this.jumpToIndex(this.content.events.findIndex(evt => evt.eventIndex == evtIndex))
  }

  jumpToIndex = (index: number) => {
    clearTimeout(this.activeTimeout)
    this.currentIndex = index
    this.writeEvent(this.currentIndex, true)
    this.emit(EventTypeHeadJump, this.makeFrameAdvanceEvent())
    if (this.isPlaying()) {
      this.currentIndex++
      this.waitForFrame()
    }
    else {
      this.advanceFrameOnPlay = true
    }
  }

  elapsedTime = ():number => this.content.events[this.currentIndex].totalDelay

  setRate = (newRate: number) => this.setDesiredRate(clamp(newRate, this.minRate, this.maxRate))
  setDesiredRate = (newRate: number) => {
    const oldRate = this.desiredPlaybackRate
    this.desiredPlaybackRate = newRate
    this.emit(EventTypeDesiredRateChange, {oldRate, newRate})
    if (this.isPlaying()) {
      this.updatePlaybackRate(this.desiredPlaybackRate)
    }
  }

  faster = () => this.setDesiredRate(Math.min(this.getRate() * 2, this.maxRate))
  slower = () => this.setDesiredRate(Math.max(this.getRate() / 2, this.minRate))

  // query state
  isPlaying = () => this.getRate() != 0
  getDuration = () => this.duration
  getStartDate = () => this.content.header.timestamp
  getRate = () => this.currentPlaybackRate
  getDesiredRate = () => this.desiredPlaybackRate
  getError = () => this.content.error
  getEventAtCursor = () => this.content.events[this.currentIndex]
  getCurrentIndex = () => this.currentIndex

  nearestEventTime = (position: number): number => this.eventTimeToRealTime(this.findClosestEventViaPosition(position))

  // Playback Script Controls

  addBookmarkAtCursor = (description: string): void => this.addBookmark(this.currentIndex, description)
  removeBookmarkAtCursor = (): void => this.removeBookmark(this.currentIndex)
  removeBookmark = (index: number): void => {
    const evt = this.getEvent(index)
    evt.bookmarks = []
  }

  // addBookmark updates the current event with the provided description. Note that this overwrites
  // whatever bookmark text was present, so acts both as an addBookmark and editBookmark function
  addBookmark = (index: number, description: string): void => {
    const evt = this.getEvent(index)
    evt.bookmarks = description.split("\n")
  }

  getBookmark = (index: number): Array<string> => this.getEvent(index).bookmarks
  getBookmarkAtCursor = (): Array<string> => this.getBookmark(this.currentIndex)

  getBookmarks = () => { // TODO: should we maintain a cached version of this, and update that on bookmark change events?
    return this.content.events.filter((evt) => evt.bookmarks.length > 0)
  }

  export = (): string => {
    const json = [JSON.stringify(this.content.header)]
    const jsonifyEvent = (evt: TerminalEvent) => JSON.stringify([evt.eventTime / 1000, 'o', evt.eventContent])
    const jsonifyBookmark = (evt: ExpandedTerminalEvent) => evt.bookmarks.map(desc => JSON.stringify([evt.eventTime / 1000, 'b', desc]))

    this.content.events.forEach((evt) => {
      json.push(jsonifyEvent(evt))
      json.push(...jsonifyBookmark(evt))
    })

    return json.join("\n")
  }
}

// Finds the closest event in a given list (via a binary search based on total delay).
const findClosestEvent = (needle: number, list: Array<TerminalEvent>, lowerBound: number, upperBound: number): number => {
  if (upperBound < lowerBound) {
    [upperBound, lowerBound] = [lowerBound, upperBound]
  }

  //  when very close, narrow in on technically closest event
  if (upperBound - lowerBound < 3) {
    while (upperBound != lowerBound) {
      const upperDiff = list[upperBound].totalDelay - needle
      const lowerDiff = needle - list[lowerBound].totalDelay
      if (lowerDiff < upperDiff) {
        upperBound--
      }
      else {
        lowerBound++
      }
    }
    return upperBound // should now be the same as lowerBound
  }

  const effectiveLength = upperBound - lowerBound
  const guessIndex = lowerBound + Math.trunc(effectiveLength / 2)
  const guess = list[guessIndex]
  if (guess.totalDelay == needle) { // extremely unlikely
    return guessIndex
  }

  let [newUpper, newLower] = [upperBound, lowerBound]

  if (guess.totalDelay < needle) {
    newLower = guessIndex + 1
  }
  else {
    newUpper = guessIndex - 1
  }

  return findClosestEvent(needle, list, newLower, newUpper)
}

const emptyHeader: TerminalRecordingHeader = {
  version: 2,
  width: 0,
  height: 0,
  timestamp: 0,
  title: "",
  env: {
    shell: "",
    term: ""
  }
}

const parseTerminalRecording = (content: string, maxDelay = 1600): TerminalRecordingData => {
  const response: TerminalRecordingData = {
    header: emptyHeader,
    events: [],
    error: null
  }

  try {
    const [header, ...rawEvents] = content
      .split("\n")
      .filter(line => line.trim() !== "")
      .map(line => JSON.parse(line))

    response.header = header
    const basicEvents = rawEvents.map((evt, idx) => ({
      eventTime: evt[0] * 1000, // convert to milliseconds
      eventSource: evt[1],
      eventContent: evt[2],
      eventIndex: idx
    }))

    for (const evt of basicEvents) {
      if (evt.eventSource == 'o') {
        response.events.push({
          ...evt,
          frameDuration: 0,
          totalDelay: 0,
          bookmarks: []
        })
      }
      if (evt.eventSource == 'b') {
        if (response.events.length > 0) {
          response.events[response.events.length - 1].bookmarks.push(evt.eventContent)
        }
      }
    }
    for (const [idx, evt] of response.events.entries()) {
      if (idx == 0) continue

      let lastEvent = response.events[idx - 1]
      lastEvent.frameDuration = Math.min(evt.eventTime - lastEvent.eventTime, maxDelay)
      evt.totalDelay = lastEvent.totalDelay + lastEvent.frameDuration
    }
  }
  catch (e) {
    response.error = e
  }

  return response
}

// clamp constrains n to a number between min and max, inclusive
// clamp(-1, 0, 4) == 0
// clamp(5, 0, 4) == 4
// clamp(1, 0, 4) == 1
// clamp(2.5, 0, 4) == 2.5
const clamp = (n: number, min: number, max: number) => Math.max(Math.min(max, n), min)
