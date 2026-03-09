import classnames from 'classnames/bind'

import { type RecoveryMetrics } from 'src/global_types'
import { useWiredData } from 'src/helpers'
import { deleteExpiredRecoveryCodes, getRecoveryMetrics } from 'src/services'

import Button from 'src/components/button'

const cx = classnames.bind(require('./stylesheet'))

export default function RecoveryMetrics(props: {}) {
  const wiredRecoveryMetrics = useWiredData<RecoveryMetrics>(getRecoveryMetrics)
  const numFormatter = new Intl.NumberFormat()
  return (
    <>
      <h1>Recovery Code Metrics</h1>
      {wiredRecoveryMetrics.render((metrics) => (
        <dl className={cx('metric-list')}>
          <dt>Expired Codes</dt>
          <dd>{numFormatter.format(metrics.expiredCount)}</dd>
          <dt>Active Codes</dt>
          <dd>{numFormatter.format(metrics.activeCount)}</dd>
        </dl>
      ))}
      <Button
        primary
        onClick={() => {
          deleteExpiredRecoveryCodes()
          wiredRecoveryMetrics.reload()
        }}
      >
        Remove Expired Codes
      </Button>
    </>
  )
}
