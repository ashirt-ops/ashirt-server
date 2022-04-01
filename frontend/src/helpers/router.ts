// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import {generatePath, useLocation, useParams, useResolvedPath, useRoutes} from 'react-router-dom'

export function subUrl(newParams: {[k: string]: string}): string {
  const {pathname} = useLocation()
  const p = useResolvedPath(pathname)
  console.log("!!!!!!!!!!!!!!!", p)
  const params = useParams()
  return generatePath(
    pathname,
    {...params, ...newParams},
  )
}
