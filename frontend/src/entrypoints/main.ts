// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { BrowserRouter } from 'react-router-dom'
import { backendDataSource } from 'src/services/data_sources/backend'
import { bootApp } from 'src/app_root'

bootApp(BrowserRouter, backendDataSource)
