// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import Linker from './linker'
import Login from './login'
import Settings from './settings'
import { AuthFrontend } from 'src/authschemes'
import { SupportedAuthenticationScheme } from 'src/global_types'

const webAuthnFrontend: AuthFrontend = {
  Linker: Linker,
  Login: Login,
  Settings: Settings,
}

export const configure = (_config: SupportedAuthenticationScheme): AuthFrontend => {
  return webAuthnFrontend
}

export default webAuthnFrontend
