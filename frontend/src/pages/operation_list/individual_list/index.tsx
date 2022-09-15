// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import NewOperationButton from '../new_operation_button'
import OperationCard from '../operation_card'
import { FilterText, Operation } from 'src/global_types'
import { UseModalOutput } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

type Header = "Other" | "Favorites" | null

const normalizedInclude = (baseString: string, term: string) => {
  return baseString.toLowerCase().includes(term.toLowerCase())
}

export default (props: {
  ops: Operation[],
  header: Header,
  newOperationModal: UseModalOutput<{}>,
  filterText: FilterText,
}) => {
  const header = props.header

  return (
      <div key={header}>
        {header && <h1 className={cx('opTitle')}>
            {header}
        </h1>}
        <div className={cx('operationList')}>
        {
          props.ops
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
                className={cx('card')}
              />
            )})
        }
        {header !== "Other" && <NewOperationButton onClick={() => props.newOperationModal.show({})} />}
      </div>
      <br/>
      </div>
    )
  }
