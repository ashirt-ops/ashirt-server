import { randomUUID as uuidv4 } from 'crypto'
import { pino, Logger } from 'pino'

let baseLogger: Logger | null = null

function setupLogging() {
  const logger = pino()
  baseLogger = logger
}

export function getAppLogger(): Logger {
  if (baseLogger == null) {
    setupLogging()
  }
  return baseLogger!
}

export function appInfoLog(message: string, extra?: Record<string, unknown>) {
  const logExtra = extra === undefined ? {} : extra
  getAppLogger().info(message)
}

export function appWarnLog(message: string, extra?: Record<string, unknown>) {
  const logExtra = extra === undefined ? {} : extra
  getAppLogger().warn(message)
}


export function getRequestLogger(): Logger {
  const logger = getAppLogger().child({ context: uuidv4() })
  return logger
}
