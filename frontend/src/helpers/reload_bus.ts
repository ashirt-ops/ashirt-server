// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { EventEmitter } from 'events'

export const reloadEvent = "reload"
export const reloadDoneEvent = "reload-done"

export const BuildReloadBus = () => {
  const bus = new EventEmitter()

  return {
    requestReload: () => {
      bus.emit(reloadEvent, null)
    },
    onReload: (listener: () => void) => {
      bus.on(reloadEvent, listener)
    },
    offReload: (listener: () => void) => {
      bus.removeListener(reloadEvent, listener)
    },

    reloadDone: () => {
      bus.emit(reloadDoneEvent, null)
    },
    onReloadDone: (listener: () => void) => {
      bus.on(reloadDoneEvent, listener)
    },
    offReloadDone: (listener: () => void) => {
      bus.removeListener(reloadDoneEvent, listener)
    },

    clean: () => {
      bus.removeAllListeners()
    }
  }
}

type Runnable = () => void
export type ReloadListenerFunc = (listener: () => void) => void

export type BusSupportedService = {
  requestReload: Runnable
  onReload: ReloadListenerFunc
  offReload: ReloadListenerFunc
  reloadDone: Runnable
  onReloadDone: ReloadListenerFunc
  offReloadDone: ReloadListenerFunc
  clean: Runnable
}
