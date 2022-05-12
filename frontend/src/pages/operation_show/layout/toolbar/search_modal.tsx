// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'

import Modal from 'src/components/modal'

const cx = classnames.bind(require('./modal.stylesheet'))

export const SearchHelpModal = (props: {
  onRequestClose: () => void,
}) => {
  const onKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      props.onRequestClose()
    }
  }

  React.useEffect(() => {
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  })

  return <Modal title="Search Help" onRequestClose={props.onRequestClose} >
    <div>
      <p>
        Timeline results can be influenced by using the Filter Timeline text box. Filters are applied
        by specifying the correct syntax, and are joined together to provide ever-narrower search
        results.
      </p>
      <p>
        In general, search queries are presented in the following form:

        <CodeExample>Free Text specific-field:"specific value"</CodeExample>

        Here, the "Free Text Section" phrase will search all evidence for any occurance of
        {" "}<CodeSnippet>Free</CodeSnippet> and <CodeSnippet>Text</CodeSnippet>. The
        {" "}<CodeSnippet>specific-field</CodeSnippet> will search just the specific-field attribute
        for the value <CodeSnippet>specific value</CodeSnippet>. Filters must be specified without a
        space on either side of the colon. Multiple filters can be provided in a single search. Also
        note that <CodeSnippet>specific value</CodeSnippet> was written in quotes. This is not
        required for any field, but provide the ability to search over phrases rather than words.
      </p>
      <p>
        As an example, consider the search:

        <CodeExample>"Darth Vader" creator:"George Lucas" tag:sci-fi</CodeExample>

        Performing this search over a set of movies might return several <em>Star Wars</em> movies but
        exclude <em>Indiana Jones</em>. Removing <CodeSnippet>Darth Vader</CodeSnippet> would
        expand the search to include <em>THX 1138</em>, while adding
        {" "}<CodeSnippet>"Jar Jar Binks"</CodeSnippet> would narrow the results to just
        {" "}<em>Star Wars</em> episodes 1-3.
      </p>
      <p>The below table lists all of the currently available filters, and value limitations, if any:</p>
      <table>
        <tbody>
          {HelpText.map(f => <FilterDescriptionRow filter={f} key={f.field} />)}
        </tbody>
      </table>
    </div>

  </Modal >
}

const CodeSnippet = (props: {
  children: React.ReactNode,
  className?: string | Array<string>
}) => <span className={cx('code', props.className)}>{props.children}</span>

const CodeExample = (props: {
  children: React.ReactNode
  className?: string | Array<string>
}) => <span className={cx('example', props.className)}>{props.children}</span>

type FilterDetail = { field: string, description: React.ReactNode }
const FilterDescriptionRow = (props: {
  filter: FilterDetail
}) => (
  <tr className={cx('filter-row')}>
    <td className={cx('filter-field')}>{props.filter.field}</td>
    <td className={cx('filter-description')}>{props.filter.description}</td>
  </tr>
)


const valuesAsCodeSnippets = (vals: Array<string>) => {
  return vals.map((v, i) => (
    <span key={i}>
      <CodeSnippet>{v}</CodeSnippet>
      {(i + 1) == vals.length ? '' : ', '}
    </span>
  ))
}

const HelpText: Array<FilterDetail> = [
  {
    field: 'tag',
    description:
      <>
        <p>
          Filters the result by requiring that the evidence or finding contain each of the
          specified tag fields.
        </p>
        <p>Multiple <CodeSnippet>tag</CodeSnippet> fields can be specified.</p>
        <p>To easily create this filter, click on the desired tags next to any evidence.</p>
      </>
  },
  {
    field: 'operator',
    description:
      <>
        <p>
          Filters the result by requiring that the evidence or finding was created by a particular
          user.
        </p>
        <p>Only one <CodeSnippet>operator</CodeSnippet> field can be specified.</p>
        <p>To easily create this filter, click on the desired username next to any evidence.</p>
      </>
  },
  {
    field: 'range',
    description:
      <>
        <p>
          Filters the result by requiring that the evidence to have occurred within a particular
          date range. In the findings timeline, this will require that all evidence for a finding
          be contained with the indicated date range. Only one range can be specified.
          Date Format: <CodeSnippet>yyyy-mm-dd,yyyy-mm-dd</CodeSnippet> where
          y, m, and d are year, month and day digits respectively.
          For example: <CodeSnippet>2020-01-01,2020-01-31</CodeSnippet> covers the entire
          month of January, 2020.
        </p>
        <p>Only one <CodeSnippet>range</CodeSnippet> field can be specified.</p>
        <p>Click on the calendar next to the Timeline Filter to help specify the date.</p>
      </>
  },
  {
    field: 'sort',
    description:
      <>
        <p>
          Orders the filter in a particular direction. By default, wiith no filter provided,
          results are ordered by "last evidence first", or an effective reverse-chronological
          order.
        </p>
        <p>
          Possible values:
          {" "}<CodeSnippet>asc</CodeSnippet>,
          {" "}<CodeSnippet>ascending</CodeSnippet> or
          {" "}<CodeSnippet>chronological</CodeSnippet>
        </p>
        <p>
          Each of the above values will order the results in a "first-evidence-first", or
          chronological order.
        </p>
        <p>Only one <CodeSnippet>sort</CodeSnippet> field can be specified.</p>
      </>
  },
  {
    field: 'linked',
    description:
      <>
        <p>
          Filters the result by finding evidence that either has, or has not been attached to a finding.
        </p>
        <p>
          Possible values: {" "}
          {valuesAsCodeSnippets(['true', 'false'])}
        </p>
        <p>
          Provide <CodeSnippet>true</CodeSnippet> to require the evidence has been linked
          with a finding, or <CodeSnippet>false</CodeSnippet> to require evidence that has
          not been linked with a finding.
          {" "}<em>This will only have an effect in the Evidence Timeline.</em>
        </p>
        <p>Only one <CodeSnippet>linked</CodeSnippet> field can be specified.</p>
      </>
  },
  {
    field: 'with-evidence',
    description:
      <>
        <p>
          Filters the result by requiring a fidning to contain a particular piece of evidence.
          <em>This will only have an effect in the Findings Timeline.</em>
        </p>
        <p>Only one <CodeSnippet>with-evidence</CodeSnippet> field can be specified.</p>
      </>
  },
  {
    field: 'type',
    description:
      <>
        <p>
          Filters the result by requiring that the evidence have the matching type as the one
          specified in the filter.
          {" "}<em>This will only have an effect in the Evidence Timeline.</em>
        </p>
        <p>
          Possible values: {" "}
          {
            valuesAsCodeSnippets([
              'image', 'codeblock', 'terminal-recording', 'http-request-cycle', 'event', 'none'
            ])
          }
        </p>
        <p>
        </p>
        <p>Only one <CodeSnippet>type</CodeSnippet> field can be specified.</p>
      </>
  },
  {
    field: 'uuid',
    description:
      <>
        <p>
          Filters the result by requiring that the evidence or finding have a particular ID.
          This is typically used to share evidence with other users. While it can be specified
          manually, the preferred method is to click the "Copy Permalink" button
          next to the desired evidence, and share the link as needed.
        </p>
        <p>Only one <CodeSnippet>uuid</CodeSnippet> field can be specified.</p>
      </>
  },

]
