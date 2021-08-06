import Linker from './linker'
import Login from './login'
import Settings from './settings'
import { AuthFrontend } from 'src/authschemes'
import { SupportedAuthenticationScheme } from 'src/global_types'

const localAuthFrontend: AuthFrontend = {
  Linker: Linker,
  Login: Login,
  Settings: Settings,
}

export const configure = (_config: SupportedAuthenticationScheme): AuthFrontend => {
  return localAuthFrontend
}

export default localAuthFrontend
