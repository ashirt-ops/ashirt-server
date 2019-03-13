// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import MarkdownRenderer from 'src/components/markdown_renderer'
import TagList from 'src/components/tag_list'
import classnames from 'classnames/bind'
import { Finding } from 'src/global_types'
import FindingStatus from '../../finding_list/finding_status'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  className?: string,
  finding: Finding,
}) => (
    <div className={cx('root', props.className)}>
      <h1 className={cx('title')}>{props.finding.title}</h1>
      <h2 className={cx('statuslabel')}>Ticket: </h2><FindingStatus className={cx('status')} finding={props.finding} />
      <TagList tags={props.finding.tags} />
      <MarkdownRenderer className={cx('description')}>{props.finding.description}</MarkdownRenderer>
    </div>
  )
