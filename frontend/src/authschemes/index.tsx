// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useAsyncComponent } from 'src/helpers'
import { SupportedAuthenticationScheme } from "src/global_types"

export type AuthFrontend = {
  Linker: React.FunctionComponent<{ onSuccess: () => void, authFlags?: Array<string> }>,
  Login: React.FunctionComponent<{ query: URLSearchParams, authFlags?: Array<string> }>,
  Settings: React.FunctionComponent<{ userKey: string, authFlags?: Array<string> }>,
}

// @ts-ignore - this is a webpack compile-time include of src/authschemes/*/index.ts
// all matched files will be code-split at compile time
// https://webpack.js.org/guides/dependency-management/#requirecontext
const ctx = require.context('src/authschemes', true, /\.\/[^/]+\/index.ts$/, 'lazy')

const bundledModuleNames: Array<string> = ctx.keys()
const loadAuthModule: (name: string) => Promise<{ default: AuthFrontend, configure: (schemeDetails: SupportedAuthenticationScheme) => AuthFrontend }> = ctx

// getAuthFrontend fetches a auth frontend module by name and returns the promise for the AuthFrontend for the given authSchemeType
// assuming its frontend is exported in `src/authschemes/<schemecode>/index.ts`
async function getAuthFrontend(authSchemeType: string, schemeDetails?: SupportedAuthenticationScheme): Promise<AuthFrontend> {
  const modulePath = `./${authSchemeType}/index.ts`
  if (bundledModuleNames.indexOf(modulePath) === -1) {
    throw Error(`Unable to load frontend auth module for "${authSchemeType}". Please make sure "${modulePath}" exists`)
  }

  const module = await loadAuthModule(modulePath)

  const rtn = schemeDetails
    ? module.configure(schemeDetails)
    : module.default

  return rtn
}

// useAuthFrontendComponent is a helper to include a specific AuthFrontend component for a given authScheme
// in an external react component.
// This is done in a typesafe way to preserve the typesignature of the component specified
//
// For example:
//
// const Login = useAuthFrontend('local', 'Login')
// return <Login query={{}} />
export function useAuthFrontendComponent<Key extends keyof AuthFrontend>(authSchemeType: string, key: Key, schemeDetails?: SupportedAuthenticationScheme): AuthFrontend[Key] {
  return useAsyncComponent(React.useCallback(() => (
    getAuthFrontend(authSchemeType, schemeDetails).then(module => module[key])
  ), [authSchemeType, key, schemeDetails]))
}
