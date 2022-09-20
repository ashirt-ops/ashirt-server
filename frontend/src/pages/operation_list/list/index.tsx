// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import NewOperationButton from '../new_operation_button'
import OperationCard from '../operation_card'
import { FilterText, Operation } from 'src/global_types'
import { UseModalOutput } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

const normalizedInclude = (baseString: string, term: string) => {
  return baseString.toLowerCase().includes(term.toLowerCase())
}

export default (props: {
  ops: Operation[],
  newOperationModal: UseModalOutput<{}>,
  filterText: FilterText,
  onFavoriteToggled: (slug: string, isFavorite: boolean) => void
}) => {
  const favoriteOps = props.ops?.filter(op => op.favorite)
  const otherOps = props.ops?.filter(op => !op.favorite)

  return (
    <div>
      <div className={cx('operationList')}>
        {
          favoriteOps
            .filter(op => normalizedInclude(op.name, props.filterText.value))
            .map(op => {
              return (
                <OperationCard
                  slug={op.slug}
                  status={op.status}
                  numUsers={op.numUsers}
                  key={op.slug}
                  name={op.name}
                  favorite={op.favorite}
                  numTags={op.numTags}
                  numEvidence={op.numEvidence}
                  onFavoriteClick={() => props.onFavoriteToggled(op.slug, !(op.favorite))}
                  className={cx('card')}
                />
              )
            })
        }
        {favoriteOps?.length && <NewOperationButton onClick={() => props.newOperationModal.show({})} />}
      </div>
      <div className={cx('operationList')}>
        {
          otherOps
            .filter(op => normalizedInclude(op.name, props.filterText.value))
            .map(op => {
              return (
                <OperationCard
                  slug={op.slug}
                  status={op.status}
                  numUsers={op.numUsers}
                  key={op.slug}
                  name={op.name}
                  favorite={op.favorite}
                  numTags={op.numTags}
                  numEvidence={op.numEvidence}
                  onFavoriteClick={() => props.onFavoriteToggled(op.slug, !(op.favorite))}
                  className={cx('card')}
                />
              )
            })
        }
        {!favoriteOps?.length && <NewOperationButton onClick={() => props.newOperationModal.show({})} />}
      </div>
      <br />
    </div>
  )
}
