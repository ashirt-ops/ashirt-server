// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import DateRangePicker from 'src/components/date_range_picker'
import Input from 'src/components/input'
import Modal from 'src/components/modal'
import classnames from 'classnames/bind'
import { default as Button, ButtonGroup } from 'src/components/button'
import { getDateRangeFromQuery, addOrUpdateDateRangeInQuery, useModal, renderModals } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))


type FilterDetail = { field: string, description: React.ReactNode }

const FilterDescriptionRow = (props: {
  filter: FilterDetail
}) => (
  <tr className={cx('filter-row')}>
    <td className={cx('filter-field')}>{props.filter.field}</td>
    <td className={cx('filter-description')}>{props.filter.description}</td>
  </tr>
)

const SearchHelpModal = (props: {
  onRequestClose: () => void,
}) => {
  const parameters: Array<FilterDetail> = [
    {
      field: 'tag',
      description:
        <>
          <p className={cx('filter-description-p')}>
            Filters the result by requiring that the evidence or finding contain each of the
            specified tag fields.
          </p>
          <p>Multiple <span className={cx('code')}>tag</span> fields can be specified.</p>
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
          <p>Only one <span className={cx('code')}>operator</span> field can be specified.</p>
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
            Date Format: <span className={cx('code')}>yyyy-mm-dd,yyyy-mm-dd</span> where
            y, m, and d are year, month and day digits respectively.
            For example: <span className={cx('code')}>2020-01-01,2020-01-31</span> covers the entire
            month of January, 2020.
          </p>
          <p>Only one <span className={cx('code')}>range</span> field can be specified.</p>
          <p>Click on the calendar next to the Timeline Filter to help specify the date.</p>
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
          <p>Only one <span className={cx('code')}>uuid</span> field can be specified.</p>
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
          <p>Only one <span className={cx('code')}>with-evidence</span> field can be specified.</p>
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
            Possible values: <span className={cx('code')}>true</span>, <span className={cx('code')}>false</span>
          </p>
          <p>
            Provide <span className={cx('code')}>true</span> to require the evidence has been linked
            with a finding, or <span className={cx('code')}>false</span> to require evidence that has
            not been linked with a finding.
            <em className={cx('space-before')}>
              This will only have an effect in the Evidence Timeline.
            </em>
          </p>
          <p>Only one <span className={cx('code')}>linked</span> field can be specified.</p>
        </>
    },
    // {
    //   field: 'sort-direction',
    //   description:
    //     <>
    //       <p>
    //         Orders the filter in a particular direction. By default, wiith no filter provided,
    //         results are ordered by "last evidence first", or an effective reverse-chronological
    //         order.
    //       </p>
    //       <p>
    //         Possible values:
    //         <span className={cx('code', 'space-before')}>asc</span>,
    //         <span className={cx('code', 'space-before')}>ascending</span> or
    //         <span className={cx('code', 'space-before')}>chronological</span>
    //       </p>
    //       <p>
    //         Each of the above values will order the results in a "first-evidence-first", or
    //         chronological order.
    //       </p>
    //       <p>Only one <span className={cx('code')}>sort-direction</span> field can be specified.</p>
    //     </>
    // },
  ]

  return <Modal title="Search Help" onRequestClose={props.onRequestClose}>
    <div>
      <p>
        Timeline results can be influenced by using the Filter Timeline text box. Filters are applied
        by specifying the correct syntax, and are joined together to provide ever-narrower search
        results.
      </p>
      <p>
        In general, search queries are presented in the following form:
        <span className={cx('example')}>Free Text specific-field:"specific value"</span>
        Here, the "Free Text Section" phrase will search all evidence for any occurance
        of <span className={cx('code')}>Free</span> and <span className={cx('code')}>Text</span>.
        The <span className={cx('code')}>specific-field</span> will search just the specific-field
        attribute for the value <span className={cx('code')}>specific value</span>. Filters must be
        specified without a space on either side of the colon. Multiple filters can be provided in
        a single search. Also note that <span className={cx('code')}>specific value</span> was
        written in quotes. This is not
        required for any field, but provide the ability to search over phrases rather than words.
      </p>
      <p>
        As an example, consider the search:
        <span className={cx('example')}>"Darth Vader" creator:"George Lucas" tag:sci-fi</span>
        Performing this search over a set of movies might return several <em>Star Wars</em> movies but
        exclude <em>Indiana Jones</em>. Removing <span className={cx('code')}>Darth Vader</span> would
        expand the search to include <em>THX 1138</em>, while
        adding <span className={cx('code')}>"Jar Jar Binks"</span> would condense the results to
        just <em>Star Wars</em> episodes 1-3.
      </p>
      <p>The below table lists all of the currently available filters, and value limitations, if any:</p>
      <table>
        <tbody>
          {parameters.map(f => <FilterDescriptionRow filter={f} key={f.field} />)}
        </tbody>
      </table>
    </div>

  </Modal>
}

export default (props: {
  onRequestCreateFinding: () => void,
  onRequestCreateEvidence: () => void,
  onSearch: (query: string) => void,
  query: string,
}) => {
  const [queryInput, setQueryInput] = React.useState<string>("")
  const helpModal = useModal<void>(modalProps => <SearchHelpModal {...modalProps} />)
  React.useEffect(() => {
    setQueryInput(props.query)
  }, [props.query])

  const inputRef = React.useRef<HTMLInputElement>(null)

  return (
    <>
      <div className={cx('root')}>
        <Input
          ref={inputRef}
          className={cx('search')}
          value={queryInput}
          onChange={setQueryInput}
          placeholder="Filter Timeline"
          icon={require('./search.svg')}
          onKeyDown={e => {
            if (e.which == 13) {
              inputRef.current?.blur()
              props.onSearch(queryInput)
            }
          }}
        />
        <ButtonGroup>
          <DateRangePicker
            range={getDateRangeFromQuery(queryInput)}
            onSelectRange={r => {
              const newQuery = addOrUpdateDateRangeInQuery(queryInput, r)
              setQueryInput(newQuery)
              props.onSearch(newQuery)
            }}
          />
          <Button onClick={e => helpModal.show()} title="Search Help">?</Button>
        </ButtonGroup>

        <ButtonGroup>
          <Button onClick={props.onRequestCreateFinding}>Create Finding</Button>
          <Button onClick={props.onRequestCreateEvidence}>Create Evidence</Button>
        </ButtonGroup>
      </div>
      {renderModals(helpModal)}
    </>
  )
}
