import Linker from './linker'
import Login from './login'
import Settings from './settings'
import {AuthFrontend} from 'src/authschemes'

const localAuthFrontend: AuthFrontend = {
  Linker: Linker,
  Login: Login,
  Settings: Settings,
}

export default localAuthFrontend
