// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import ErrorDisplay from '../../components/error_display'
import {RouteComponentProps} from 'react-router-dom'

export default (props: RouteComponentProps) => (
  <ErrorDisplay err={new Error(`404 - The path ${props.location.pathname} is invalid`)} />
)
