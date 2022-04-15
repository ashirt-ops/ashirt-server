// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ErrorDisplay from '../../components/error_display'
import { useLocation } from 'react-router-dom'

export default () => {
  const { pathname } = useLocation()
  return (
    <ErrorDisplay err={new Error(`404 - The path ${pathname} is invalid`)} />
  )
}

