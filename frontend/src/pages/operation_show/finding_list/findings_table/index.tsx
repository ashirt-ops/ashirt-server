import * as React from 'react'
import TagList from 'src/components/tag_list'
import classnames from 'classnames/bind'
import Table from 'src/components/table'
import { Finding } from 'src/global_types'
import { Link } from 'react-router'
import { default as Button, ButtonGroup } from 'src/components/button'
import FindingStatus from '../finding_status'
import { format } from 'date-fns'
const cx = classnames.bind(require('./stylesheet'))

export default (props: {
  findings: Array<Finding>,
  onDelete: (f: Finding) => void,
  onEdit: (f: Finding) => void,
  operationSlug: string,
}) => (
    <Table className={cx('table')} columns={['Title', 'Category', 'Ticket', '# Evidence', 'Date Range', 'Tags']}>
      {props.findings.map(finding => (
        <tr key={finding.uuid}>
          <td className={cx('title-cell')}>
            <Link to={`/operations/${props.operationSlug}/findings/${finding.uuid}`} className={cx('title')}>
              {finding.title}
            </Link>
            <ButtonGroup className={cx('actions')}>
              <Button small onClick={() => props.onEdit(finding)}>Edit</Button>
              <Button small onClick={() => props.onDelete(finding)}>Delete</Button>
            </ButtonGroup>
          </td>
          <td>{finding.category}</td>
          <td><FindingStatus finding={finding} /></td>
          <td>{finding.numEvidence}</td>
          <td>{finding.occurredFrom && finding.occurredTo ? (
            `${format(finding.occurredFrom, 'MM/dd/yy')}-${format(finding.occurredTo, 'MM/dd/yy')}`
          ) : (
              'N/A'
            )}</td>
          <td className={cx('tags-cell')}>
            <TagList tags={finding.tags} />
          </td>
        </tr>
      ))}
    </Table>
  )
