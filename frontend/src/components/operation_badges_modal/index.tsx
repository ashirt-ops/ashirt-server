// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { EvidenceCount, TopContrib } from 'src/global_types'
import Modal from '../modal'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onRequestClose: () => void,
  topContribs: Array<TopContrib>,
  evidenceCount: EvidenceCount,
}) => {

  const evidenceNameMap = {
    imageCount: 'Images',
    codeblockCount: 'Codeblocks',
    recordingCount: 'Recordings',
    eventCount: 'Events',
    harCount: 'HAR files',
  }

  type ObjectKey = keyof typeof evidenceNameMap;

  const evidencePresent = Object.values(props.evidenceCount).reduce((p, c) => p + c, 0) > 0 

  return (
    <Modal title="More Details" onRequestClose={props.onRequestClose}>
      <div className={cx("root")}>
          <div>
          {!!props.topContribs?.length && (
            <>
              <h1 className={cx('modal-heading')}>Top Contributor{props?.topContribs?.length > 1 && "s"}</h1>
              {props?.topContribs?.map(contrib => (
              <div className={cx("inner-div")} key={`${contrib.slug}`}>
                <p className={cx("row-item")}>{contrib.slug}: </p>
                <p className={cx("row-item", "right")}>{contrib.count}</p>
              </div>)
              )}
            </>
          )}
          </div>
          <div className={cx("column")}>
            {evidencePresent && (
              <>
                <h1 className={cx('modal-heading')}>Evidence by Category</h1>
                {Object.entries(props.evidenceCount).map(ebc => ebc[1] > 0 && (
                  <div key={`${ebc[0]}`} className={cx("inner-div")}>
                    <p className={cx("row-item")} >{evidenceNameMap[ebc[0] as ObjectKey]}: </p>
                    <p className={cx("row-item", "right")}>{ebc[1]}</p>
                  </div>)
                )}
            </>)}
          </div>
      </div>
    </Modal>
  )
}
