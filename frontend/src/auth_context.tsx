// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import {UserOwnView} from 'src/global_types'

export type AuthContextType = {
  user: UserOwnView | null,
}

export default React.createContext<AuthContextType>({
  user: null,
});
