import { type ReactNode, type MouseEvent, type KeyboardEvent, useRef, useEffect } from 'react'
import classnames from 'classnames/bind'
const cx = classnames.bind(require('./stylesheet'))

const Menu = (props: { maxHeight?: number; children: ReactNode }) => (
  <div className={cx('root')} style={{ maxHeight: props.maxHeight }}>
    {props.children}
  </div>
)

export const MenuItem = (props: {
  children: ReactNode
  icon?: string
  onClick?: (e: MouseEvent<Element>) => void
  selected?: boolean
  disabled?: boolean
  danger?: boolean
  onKeyUp?: (e: KeyboardEvent) => void
  onKeyDown?: (e: KeyboardEvent) => void
}) => {
  const ref = useRef<HTMLButtonElement | null>(null)

  useEffect(() => {
    if (props.selected && ref.current != null) {
      ref.current.scrollIntoView({ block: 'nearest' })
    }
  }, [props.selected])

  return (
    <button
      disabled={props.disabled}
      className={cx('menu-item', {
        selected: props.selected,
        clickable: props.onClick && !props.disabled,
        disabled: props.disabled,
        danger: props.danger,
      })}
      onClick={props.onClick}
      onKeyUp={props.onKeyUp}
      onKeyDown={props.onKeyDown}
      ref={ref}
    >
      {props.icon && (
        <div className={cx('icon')} style={{ backgroundImage: `url(${props.icon})` }} />
      )}
      {props.children}
    </button>
  )
}

export const MenuSeparator = () => <hr className={cx('separator')} />
export default Menu
