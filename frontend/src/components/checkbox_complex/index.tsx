import { type ChangeEventHandler } from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const CheckboxComplex = (props: {
  className?: string
  value?: boolean
  label?: string
  title?: string
  disabled?: boolean
  onChange?: ChangeEventHandler<HTMLInputElement>
}) => (
  <label
    title={props.title}
    className={cx('root', props.className)}
    onClick={(e) => e.stopPropagation()}
  >
    <input
      type="checkbox"
      checked={props.value}
      disabled={props.disabled}
      title={props.title}
      onChange={props.onChange}
    />
    <div className={cx('indicator')} />
    {props.label}
  </label>
)
export default CheckboxComplex
