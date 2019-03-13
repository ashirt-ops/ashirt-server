// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  children: React.ReactNode,
  className?: string,
  title: string,
  width?: "wide"|"normal"|"narrow",
}) => (
  <section className={cx('root', props.width || "normal", props.className)}>
    <h1 className={cx('title')}>{props.title}</h1>
    {props.children}
  </section>
)
