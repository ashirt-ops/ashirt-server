import { type ReactNode, useState, useEffect } from 'react'
import classnames from 'classnames/bind'
import { createPortal } from 'react-dom'
const cx = classnames.bind(require('./stylesheet'))

export default function Lightbox(props: {
  children: ReactNode
  isOpen: boolean
  canUseFitToggle?: boolean
  onRequestClose: () => void
}) {
  const [exists, setExists] = useState<boolean>(false)
  const [animating, setAnimating] = useState<boolean>(true)
  const [full, setFull] = useState<boolean>(false)

  useEffect(() => {
    if (props.isOpen) {
      setExists(true)
      setAnimating(true)
      setTimeout(() => setAnimating(false), 20)
    } else {
      setAnimating(true)
      setTimeout(() => setExists(false), 200)
    }
  }, [props.isOpen])

  const onKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') props.onRequestClose()
    if (e.key === 'z' || e.key === 'Z') setFull(!full)
  }

  useEffect(() => {
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  })

  if (!exists) return null

  const toggleClasses = props.canUseFitToggle
    ? [full ? 'full' : 'fit', props.canUseFitToggle ? 'can-fit' : '']
    : []

  return createPortal(
    <div className={cx('root', animating ? 'animating' : 'open')} onClick={props.onRequestClose}>
      <div className={cx('content', ...toggleClasses)} onClick={(e) => e.stopPropagation()}>
        {props.children}
      </div>
    </div>,
    document.body,
  )
}
