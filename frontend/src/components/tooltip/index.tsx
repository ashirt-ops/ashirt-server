import { type ReactElement, useRef, useState, useEffect } from 'react'
import classnames from 'classnames/bind'
import { createPortal } from 'react-dom'

import { useElementRect } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

const RawTooltip = (props: { content: string; position: 'above' | 'below' | 'left' | 'right' }) => (
  <div className={cx('root', props.position)}>{props.content}</div>
)

// Tooltip (the default export) is a controlled component where the caller must explicitly
// set if the tooltip is open or not (useful for displaying a tip that goes away after a set
// time rather than based on hover
const Tooltip = (props: {
  children: ReactElement
  content: string
  isOpen: boolean
  onMouseOut?: () => void
  onMouseOver?: () => void
  position?: 'above' | 'below' | 'left' | 'right'
}) => {
  const targetRef = useRef<HTMLDivElement | null>(null)
  const [exists, setExists] = useState(false)
  const [animating, setAnimating] = useState(false)
  const rect = useElementRect(exists ? targetRef : null)

  useEffect(() => {
    let timeout: number
    if (props.isOpen) {
      setExists(true)
      setAnimating(true)
      timeout = window.setTimeout(() => {
        setAnimating(false)
      }, 20)
    } else {
      setAnimating(true)
      timeout = window.setTimeout(() => {
        setExists(false)
      }, 200)
    }
    return () => {
      clearTimeout(timeout)
    }
  }, [props.isOpen])

  return (
    <>
      <div
        className={cx('target')}
        ref={targetRef}
        children={props.children}
        onMouseOver={props.onMouseOver}
        onMouseOut={props.onMouseOut}
      />
      {exists &&
        rect &&
        createPortal(
          <div
            className={cx('positioner', { animating })}
            style={{ top: rect.top, left: rect.left, width: rect.width, height: rect.height }}
          >
            <RawTooltip content={props.content} position={props.position || 'above'} />
          </div>,
          document.body,
        )}
    </>
  )
}

// Hovertooltip is a tooltip that displays when the passed children emit an
// onMouseOver event and hides when the children emit an onMouseOut event
export const HoverTooltip = (props: {
  children: ReactElement
  content: string
  position?: 'above' | 'below' | 'left' | 'right'
}) => {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <Tooltip
      children={props.children}
      content={props.content}
      isOpen={isOpen}
      onMouseOut={() => setIsOpen(false)}
      onMouseOver={() => setIsOpen(true)}
      position={props.position}
    />
  )
}

export default Tooltip
