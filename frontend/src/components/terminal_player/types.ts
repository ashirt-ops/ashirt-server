
export type BasicTerminalEvent = {
  eventTime: number
  eventSource: string
  eventContent: string
  eventIndex: number
}

export type TerminalEvent = BasicTerminalEvent & {
  totalDelay: number
  frameDuration: number
}

export type ExpandedTerminalEvent = TerminalEvent & {
  bookmarks: Array<string>
}

export type TerminalRecordingHeader = {
  version: 2
  width: number
  height: number
  timestamp: number
  title: string
  env: {
    shell: string
    term: string
  }
}

export type TerminalRecordingData = {
  header: TerminalRecordingHeader
  events: Array<ExpandedTerminalEvent>
  error: string | null
}

export type RateChangeEventBody = {
  oldRate: number
  newRate: number
}

export type PositionChangeEventBody = {
  elapsedTime: number
  playbackPosition: number
  index: number
  terminalTime: number
}
