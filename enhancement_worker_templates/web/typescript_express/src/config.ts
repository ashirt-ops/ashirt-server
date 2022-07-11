import dotenv from 'dotenv'

let data: Config | null = null
export type Config = {
  port: string
  enableDev: boolean
  backendUrl: string
  accessKey: string
  secretKeyB64: string
}

/**
 * config is a memoized function to read and parse the environment variables
 * @returns the parsed env as a Config type
 */
export function config() {
  if (data === null) {
    dotenv.config()
    data = loadConfig()
  }

  return data
}

/**
 * loadConfig loads the configuration from env. If any value is missing, a default value is used
 * instead
 * @returns The parsed config, in Record from
 */
const loadConfig = (): Config => {
  dotenv.config()
  return {
    port: process.env.PORT ?? "8001",
    enableDev: (process.env.ENABLE_DEV ?? "false").toLowerCase() == 'true',
    backendUrl: process.env.ASHIRT_BACKEND_URL ?? "",
    accessKey: process.env.ASHIRT_ACCESS_KEY ?? "",
    secretKeyB64: process.env.ASHIRT_SECRET_KEY ?? ""
  }
}
