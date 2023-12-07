import * as React from 'react'
import {UserOwnView} from 'src/global_types'

export type AuthContextType = {
  user: UserOwnView | null,
}

export default React.createContext<AuthContextType>({
  user: null,
});
