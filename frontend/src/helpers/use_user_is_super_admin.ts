// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import AuthContext from 'src/auth_context'

export function useUserIsSuperAdmin() {
  const user = React.useContext(AuthContext).user
  return user != null && user.admin
}
