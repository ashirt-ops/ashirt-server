import * as React from 'react'
import classnames from 'classnames/bind'

import Button from 'src/components/button'

const cx = classnames.bind(require('./stylesheet'))

const Pager = (props: {
  pageNumber: number
  maxPageNumber?: number
  children: React.ReactNode
  className?: string
  nextButtonText?: string
  prevButtonText?: string
  onPageUp: () => void
  onPageDown: () => void
}) => (
    <div className={cx('root', props.className)}>
      <Button disabled={props.pageNumber == 1} onClick={props.onPageDown}>{props.prevButtonText || 'previous'}</Button>
      {props.children}
      <Button disabled={props.maxPageNumber == undefined ? false : props.maxPageNumber <= props.pageNumber} onClick={props.onPageUp}>{props.nextButtonText || 'next'}</Button>
    </div>
  )

export const StandardPager = (props: {
  page: number
  maxPages?: number
  className?: string
  onPageChange: (newPage: number) => void
}) => {
  const setRelativePage = (offset: number) => {
    const nextPage = Math.max(1, props.page + offset)
    props.onPageChange(nextPage)
  }
  const incPage = () => setRelativePage(1)
  const decPage = () => setRelativePage(-1)


  return (
    <Pager
      maxPageNumber={props.maxPages}
      className={props.className}
      pageNumber={props.page}
      onPageUp={incPage}
      onPageDown={decPage}>
      <div className={cx('pageNumber')}>{props.page}</div>
    </Pager>
  )
}

export default Pager
