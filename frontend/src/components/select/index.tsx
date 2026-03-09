import { type ReactNode, Children, isValidElement } from 'react'
import classnames from 'classnames/bind'
import WithLabel from 'src/components/with_label'
const cx = classnames.bind(require('./stylesheet'))

const Select = (props: {
  children: ReactNode
  className?: string
  disabled?: boolean
  label?: string
  onChange: (v: string) => void
  value: string
}) => (
  <WithLabel className={cx(props.className)} label={props.label}>
    <div className={cx('wrapper')}>
      <select
        value={props.value}
        onChange={(e) => props.onChange(e.target.value)}
        disabled={props.disabled}
      >
        {props.children}
      </select>
      <div className={cx('visual', { disabled: props.disabled })}>
        {getNameForValue(props.value, props.children)}
      </div>
    </div>
  </WithLabel>
)
export default Select

function getNameForValue(curValue: string, children: ReactNode): string {
  let name = curValue
  Children.forEach(children, (child) => {
    if (!isValidElement(child)) {
      return
    }

    const { value, children: childChildren } = child.props as {
      value?: string
      children?: string
    }

    const effectiveValue = value ?? childChildren

    if (effectiveValue === curValue) {
      name = typeof childChildren === 'string' ? childChildren : curValue
    }
  })

  return name
}
