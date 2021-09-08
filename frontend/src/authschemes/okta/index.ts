import Linker from './linker'
import Login from './login'
import { AuthFrontend } from 'src/authschemes'
import { SupportedAuthenticationScheme } from 'src/global_types'

const oktaAuthFrontend: AuthFrontend = {
  Linker: Linker,
  Login: Login,
  Settings: () => null,
}

export const configure = (_config: SupportedAuthenticationScheme): AuthFrontend => {
  return oktaAuthFrontend
}

export default oktaAuthFrontend
