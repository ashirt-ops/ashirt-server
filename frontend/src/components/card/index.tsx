import { type ReactNode } from 'react'
import classnames from 'classnames/bind'
import { Link } from 'react-router'
const cx = classnames.bind(require('./stylesheet'))

const Card = (props: { className?: string; linkTo?: string; children: ReactNode }) =>
  props.linkTo == null ? (
    <div className={cx('root', props.className)}>{props.children}</div>
  ) : (
    <Link to={props.linkTo} className={cx('root', props.className)}>
      {props.children}
    </Link>
  )
export default Card
