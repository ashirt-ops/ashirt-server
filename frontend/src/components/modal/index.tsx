import { type ReactNode, useId, useState, useLayoutEffect } from 'react'
import classnames from 'classnames/bind'
import { createPortal } from 'react-dom'
import { PopoverPortalContext } from 'src/components/popover/portal_context'
const cx = classnames.bind(require('./stylesheet'))

export default function Modal(props: {
  children: ReactNode
  onRequestClose: () => void
  title: string
  smallerWidth?: boolean
}) {
  const titleId = useId()
  const [dialog, setDialog] = useState<HTMLDialogElement | null>(null)

  useLayoutEffect(() => {
    if (!dialog) return
    dialog.showModal()
    return () => dialog.close()
  }, [dialog])

  return createPortal(
    <dialog
      className={cx('root')}
      aria-labelledby={titleId}
      aria-modal="true"
      ref={setDialog}
      onCancel={(e) => {
        e.preventDefault()
        props.onRequestClose()
      }}
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) {
          // Do not let the backdrop steal focus after close() restores the opener.
          e.preventDefault()
          props.onRequestClose()
        }
      }}
    >
      <PopoverPortalContext.Provider value={dialog}>
        <div
          className={cx('modal', props.smallerWidth ? 'smaller-width' : '')}
          onMouseDown={(e) => e.stopPropagation()}
        >
          <h1 id={titleId} className={cx('title')}>
            {props.title}
          </h1>
          <div className={cx('content')}>{props.children}</div>
        </div>
      </PopoverPortalContext.Provider>
    </dialog>,
    document.body,
  )
}
