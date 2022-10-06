// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import Card from 'src/components/card'
import OperationBadges from 'src/components/operation_badges'
import OperationBadgesModal from 'src/components/operation_badges_modal'
import classnames from 'classnames/bind'
import { Link } from 'react-router-dom'
import { EvidenceCount, OperationStatus, TopContrib } from 'src/global_types'
import Button from 'src/components/button'
import { renderModals, useModal } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  slug: string,
  name: string,
  numEvidence: number,
  numTags: number,
  numUsers: number,
  status: OperationStatus,
  favorite: boolean,
  onFavoriteClick: () => void,
  topContribs: Array<TopContrib>,
  evidenceCount: EvidenceCount,
}) => {
  const { favorite } = props
  const moreDetailsModal = useModal<{}>(modalProps => (
    <OperationBadgesModal {...modalProps} topContribs={props.topContribs} evidenceCount={props.evidenceCount} status={props.status} />
  ))
  const handleDetailsModal = () => moreDetailsModal?.show({})
  return (
    <Card className={cx('root', props.className)}>
      <Link className={cx('name')} to={`/operations/${props.slug}/evidence`}>{props.name}</Link>
      <OperationBadges className={cx('badges')} numUsers={props.numUsers} status={props.status} numTags={props.numTags} numEvidence={props.numEvidence} showDetailsModal={handleDetailsModal} />
      <Link className={cx('edit')} to={`/operations/${props.slug}/edit`} title="Edit this operation" />
      <Link className={cx('overview')} to={`/operations/${props.slug}/overview`} title="Evidence Overview" />
      <Button
        className={cx('favorite-button', favorite && 'filled')}
        onClick={props.onFavoriteClick}>
      </Button>
      {renderModals(moreDetailsModal)}
    </Card>
  )
}
