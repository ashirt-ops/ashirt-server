// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  onClick: () => void,
}) => (
  <button
    className={cx('root')}
    onClick={props.onClick}
  >
    <div>+</div>
    New Operation
  </button>
)
