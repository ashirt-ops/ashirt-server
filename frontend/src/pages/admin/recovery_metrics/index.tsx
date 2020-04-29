// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { RecoveryMetrics } from 'src/global_types'
import { useDataSource, deleteExpiredRecoveryCodes, getRecoveryMetrics } from 'src/services'
import { useWiredData } from 'src/helpers'

import Button from 'src/components/button'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {}) => {
  const ds = useDataSource()
  const wiredRecoveryMetrics = useWiredData<RecoveryMetrics>(React.useCallback(() => (
    getRecoveryMetrics(ds)
  ), [ds]))
  const numFormatter = new Intl.NumberFormat()
  return (
    <>
      <h1>Recovery Code Metrics</h1>
      {wiredRecoveryMetrics.render(metrics => <dl className={cx('metric-list')}>
        <dt>Expired Codes</dt>
        <dd>{numFormatter.format(metrics.expiredCount)}</dd>
        <dt>Active Codes</dt>
        <dd>{numFormatter.format(metrics.activeCount)}</dd>
      </dl>
      )}
      <Button primary onClick={() => {
        deleteExpiredRecoveryCodes(ds)
          .then(wiredRecoveryMetrics.reload)
      }}>Remove Expired Codes</Button>
    </>
  )
}
