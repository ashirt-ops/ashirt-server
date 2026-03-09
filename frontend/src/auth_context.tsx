import { createContext } from 'react'
import { type UserOwnView } from 'src/global_types'

export type AuthContextType = {
  user: UserOwnView | null
}

export default createContext<AuthContextType>({
  user: null,
})
