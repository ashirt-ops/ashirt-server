import { type ReactNode, useRef, useEffect } from 'react'
import classnames from 'classnames/bind'
import { createPortal } from 'react-dom'
import { useFocusFirstFocusableChild } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

export default function Modal(props: {
  children: ReactNode
  onRequestClose: () => void
  title: string
  smallerWidth?: boolean
}) {
  const rootRef = useRef(null)
  useFocusFirstFocusableChild(rootRef)

  useEffect(() => {
    const main = document.querySelector('main')
    if (main == null) return
    main.style.filter = 'blur(5px)'
    return () => {
      main.style.removeProperty('filter')
    }
  })

  return createPortal(
    <div className={cx('root')} onMouseDown={props.onRequestClose} ref={rootRef}>
      <div
        className={cx('modal', props.smallerWidth ? 'smaller-width' : '')}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <h1 className={cx('title')}>{props.title}</h1>
        <div className={cx('content')}>{props.children}</div>
      </div>
    </div>,
    document.body,
  )
}
