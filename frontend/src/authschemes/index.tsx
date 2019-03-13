// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { ParsedUrlQuery } from 'querystring'
import { useAsyncComponent } from 'src/helpers'

export type AuthFrontend = {
  Linker: React.FunctionComponent<{ onSuccess: () => void }>,
  Login: React.FunctionComponent<{ query: ParsedUrlQuery }>,
  Settings: React.FunctionComponent<{ userKey: string }>,
}

// @ts-ignore - this is a webpack compile-time include of src/authschemes/*/index.ts
// all matched files will be code-split at compile time
// https://webpack.js.org/guides/dependency-management/#requirecontext
const ctx = require.context('src/authschemes', true, /\.\/[^/]+\/index.ts$/, 'lazy')

const bundledModuleNames: Array<string> = ctx.keys()
const loadAuthModule: (name: string) => Promise<{default: AuthFrontend}> = ctx

// getAuthFrontend fetches a auth frontend module by name and returns the promise for the AuthFrontend for the given authschemecode
// assuming its frontend is exported in `src/authschemes/<schemecode>/index.ts`
async function getAuthFrontend(authSchemeCode: string): Promise<AuthFrontend> {
  const modulePath = `./${authSchemeCode}/index.ts`
  if (bundledModuleNames.indexOf(modulePath) === -1) {
    throw Error(`Unable to load frontend auth module for "${authSchemeCode}". Please make sure "${modulePath}" exists`)
  }

  const module = await loadAuthModule(modulePath)
  return module.default
}

// useAuthFrontendComponent is a helper to include a specific AuthFrontend component for a given authScheme
// in an external react component.
// This is done in a typesafe way to preserve the typesignature of the component specified
//
// For example:
//
// const Login = useAuthFrontend('local', 'Login')
// return <Login query={{}} />
export function useAuthFrontendComponent<Key extends keyof AuthFrontend>(authSchemeCode: string, key: Key): AuthFrontend[Key] {
  return useAsyncComponent(React.useCallback(() => (
    getAuthFrontend(authSchemeCode).then(module => module[key])
  ), [authSchemeCode, key]))
}
