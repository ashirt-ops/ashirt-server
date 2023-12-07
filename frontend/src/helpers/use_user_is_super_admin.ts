import * as React from 'react'
import AuthContext from 'src/auth_context'

export function useUserIsSuperAdmin() {
  const user = React.useContext(AuthContext).user
  return user != null && user.admin
}
