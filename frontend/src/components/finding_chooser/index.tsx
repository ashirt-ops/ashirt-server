import * as React from 'react'
import Chooser from 'src/components/chooser'
import Lightbox from 'src/components/lightbox'
import MarkdownRenderer from 'src/components/markdown_renderer'
import TagList from 'src/components/tag_list'
import classnames from 'classnames/bind'
import { Finding } from 'src/global_types'
import { getFindings } from 'src/services'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  disabled?: boolean,
  onChange: (e: Array<Finding>) => void,
  operationSlug: string,
  value: Array<Finding>,
}) => {
  const fetchFindings = React.useCallback((query: string) => getFindings({ operationSlug: props.operationSlug, query }), [props.operationSlug])
  return (
    <Chooser
      {...props}
      placeholder="Filter Findings"
      fetch={fetchFindings}
      renderRow={finding => <FindingRow finding={finding} />}
    />
  )
}

const FindingRow = (props: {
  finding: Finding,
}) => {
  const [lightboxOpen, setLightboxOpen] = React.useState(false)

  return (
    <div>
      <div className={cx('title')}>
        <a href="#" onClick={e => { e.stopPropagation(); e.preventDefault(); setLightboxOpen(true) }}>{props.finding.title}</a>
        <span>{props.finding.category}</span>
      </div>
      <TagList tags={props.finding.tags} />
      <div onClick={e => e.stopPropagation()}>
        <Lightbox isOpen={lightboxOpen} onRequestClose={() => setLightboxOpen(false)}>
          <div className={cx('lightbox')}>
            <h1>{props.finding.title}</h1>
            <MarkdownRenderer>{props.finding.description}</MarkdownRenderer>
          </div>
        </Lightbox>
      </div>
    </div>
  )
}
