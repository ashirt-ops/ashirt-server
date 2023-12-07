import * as React from 'react'
import classnames from 'classnames/bind'
import { EvidenceCount, TopContrib } from 'src/global_types'
import Modal from '../modal'

const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onRequestClose: () => void,
  topContribs: Array<TopContrib>,
  evidenceCount: EvidenceCount,
  numTags: number;
}) => {

  const evidenceNameMap = {
    imageCount: 'Image',
    codeblockCount: 'Codeblock',
    recordingCount: 'Recording',
    eventCount: 'Event',
    harCount: 'HAR file',
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
            Object.entries(props.evidenceCount).map(ebc => {
              const count = ebc[1]

              const label = evidenceNameMap[ebc[0] as ObjectKey]
              const modLabel = count > 1 ? `${label}s` : label
              const upperCaseLabel = modLabel.toUpperCase()

              return count > 0 && (
                <div className={cx('grouping')}>
                  <h1 className={cx('large-text', 'group')}>{`${ebc[1]} `}</h1>
                  <p className={cx('supporting-text')}>{upperCaseLabel}</p>
                </div>
              )
            })
          )}
          </div>
      </div>
    </Modal>
  )
}
