import Linker from './linker'
import Login from './login'
import {AuthFrontend} from 'src/authschemes'

const oktaAuthFrontend: AuthFrontend = {
  Linker: Linker,
  Login: Login,
  Settings: () => null,
}

export default oktaAuthFrontend
