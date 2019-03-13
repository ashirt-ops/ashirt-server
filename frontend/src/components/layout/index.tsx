// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import NavBar from './nav_bar'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children: React.ReactNode,
}) => (
  <main className={cx('root')}>
    <NavBar />
    <div className={cx('children')}>
      {props.children}
    </div>
  </main>
)
