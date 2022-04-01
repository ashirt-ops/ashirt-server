// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import {generatePath, useRouteMatch, useParams} from 'react-router-dom'

export function subUrl(newParams: {[k: string]: string}): string {
  const match = useRouteMatch()
  const params = useParams()
  return generatePath(
    match.path,
    {...params, ...newParams},
  )
}
