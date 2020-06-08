// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { HashRouter } from 'react-router-dom'
import { bootApp } from 'src/app_root'
import { fetchJsonp } from 'src/helpers/fetch_jsonp'
import { makeArchiveDataSource, OperationArchiveData } from 'src/services/data_sources/archive'
import { mapValues } from 'lodash'

fetchJsonp('archiveJsonp', 'data.json')
  .then((data: any) => {
    // Golang serializes empty arrays to null. This restores them to their proper type
    // before passing to the archive data source
    return mapValues(data, value => (value == null) ? [] : value)
  })
  .then((data: OperationArchiveData) => { bootApp(HashRouter, makeArchiveDataSource(data)) })
  .catch(err => alert(err.message))
