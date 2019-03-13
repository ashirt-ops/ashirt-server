// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import {RouteComponentProps, generatePath} from 'react-router-dom'

export function subUrl(pageProps: RouteComponentProps, newParams: {[k: string]: string}): string {
  return generatePath(
    pageProps.match.path,
    {...pageProps.match.params, ...newParams},
  )
}
