// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import {OperationStatus, operationStatusToLabel} from 'src/global_types'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  numEvidence?: number,
  numTags?: number,
  numUsers: number,
  status: OperationStatus,
  showDetailsModal: () => void,
}) => (
  <button className={cx('root', props.className)} onClick={() => props.showDetailsModal()}>
     <div
      className={cx('status', `status-${props.status}`)}
      title={`Operation status: ${operationStatusToLabel[props.status]}`}
    />
     <div
      className={cx('icon', 'tags')}
      title={`${props.numTags} tag${props.numTags === 1 ? ' belongs' : 's belong'} to this operation`}
      children={props.numTags}
    />
    <div
      className={cx('icon', 'evidence')}
      title={`${props.numEvidence} pieces of evidence${props.numEvidence === 1 ? ' belongs' : 's belong'} to this operation`}
      children={props.numEvidence}
    />
    <div
      className={cx('icon', 'users')}
      title={`${props.numUsers} user${props.numUsers === 1 ? ' belongs' : 's belong'} to this operation`}
      children={props.numUsers}
    />
  </button>
)
