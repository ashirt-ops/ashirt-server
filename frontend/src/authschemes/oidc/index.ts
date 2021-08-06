import { makeLinker } from './linker'
import { makeLogin } from './login'
import { AuthFrontend } from 'src/authschemes'
import { SupportedAuthenticationScheme } from 'src/global_types'

export type OIDCInstanceConfig = {
  code: string
  name: string
}

export const configure = (config: SupportedAuthenticationScheme): AuthFrontend => {
  const oidcConfig: OIDCInstanceConfig = {
    name: config.schemeName,
    code: config.schemeCode,
  }
  return {
    Linker: makeLinker(oidcConfig),
    Login: makeLogin(oidcConfig),
    Settings: () => null
  }
}

const defaultConfig: OIDCInstanceConfig = {
  code: "oidc",
  name: "Unconfigured OIDC" // you should never see this
}

const oidcAuthFrontend: AuthFrontend = {
  Linker: makeLinker(defaultConfig),
  Login: makeLogin(defaultConfig),
  Settings: () => null,
}

export default oidcAuthFrontend
