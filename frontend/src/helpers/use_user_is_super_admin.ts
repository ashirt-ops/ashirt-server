import { useContext } from 'react'
import AuthContext from 'src/auth_context'

export function useUserIsSuperAdmin() {
  const user = useContext(AuthContext).user
  return user != null && user.admin
}
