import * as React from 'react'
import { Finding } from 'src/global_types'
import classnames from 'classnames/bind'
import { trimURL } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  finding: Finding
  className?: string
}) => {
  let content
  if (!props.finding.readyToReport) {
    content = <em>Pending</em>
  } else if (!props.finding.ticketLink) {
    content = <em>Ready to Report</em>
  } else {
    content = <a href={props.finding.ticketLink} target="_blank">{trimURL(props.finding.ticketLink).trimmedValue}</a>
  }
  return <div className={cx('root', props.className)}>{content}</div>
}
