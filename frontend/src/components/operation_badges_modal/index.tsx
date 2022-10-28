// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { EvidenceCount, OperationStatus, TopContrib } from 'src/global_types'
import Modal from '../modal'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onRequestClose: () => void,
  topContribs: Array<TopContrib>,
  evidenceCount: EvidenceCount,
  status: OperationStatus,
  numTags: number;
}) => {

  const evidenceNameMap = {
    imageCount: 'Images',
    codeblockCount: 'Codeblocks',
    recordingCount: 'Recordings',
    eventCount: 'Events',
    harCount: 'HAR files',
  }
  type ObjectKey = keyof typeof evidenceNameMap;
  const evidencePresent = Object.values(props.evidenceCount).reduce((prev, curr) => prev + curr, 0) > 0

  const numTags = props.numTags

  return (
    <Modal smallerWidth={true} title='More Details' onRequestClose={props.onRequestClose}>
      <div className={cx('root')}>
          <div>
          {numTags && (
              <div className={cx('grouping')}>
                <h1 className={cx('large-text', 'group')}>{numTags}</h1>
                <p className={cx('supporting-text')}>TAGS</p>
              </div>
            )}
            {!!props.topContribs?.length && props?.topContribs?.map(contrib =>
                (<div className={cx('grouping')}>
                  <h1 className={cx('large-text')}>{contrib.slug}</h1>
                  <p className={cx('small-supporting-text')}>TOP CONTRIBUTOR</p>
                </div>)
            )}
            </div>
          <div>
          {evidencePresent && (
            Object.entries(props.evidenceCount).map(ebc => ebc[1] > 0 && (
                  <div className={cx('grouping')}>
                    <h1 className={cx('large-text', 'group')}>{`${ebc[1]} `}</h1>
                    <p className={cx('supporting-text')}>{evidenceNameMap[ebc[0] as ObjectKey].toUpperCase()}</p>
                  </div>
                )
              )
          )}
          </div>
      </div>
    </Modal>
  )
}
