// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import classnames from 'classnames/bind'
import { createPortal } from 'react-dom'
import { useWindowSize, useElementRect } from 'src/helpers'

const cx = classnames.bind(require('./stylesheet'))

function useEventListenerUnlessRef(ref: React.MutableRefObject<HTMLElement|null>, eventName: string, fn: () => void) {
  const callUnlessThisEl = (e: Event) => {
    for (let target = e.target; target !== ref.current; target = (target as HTMLElement).parentElement) {
      if (target == null) return fn()
    }
  }
  React.useEffect(() => {
    document.addEventListener(eventName, callUnlessThisEl)
    return () => document.removeEventListener(eventName, callUnlessThisEl)
  })
}

// This component displays content in a react portal on the page directly under the children element.
// If there is not enough room to display the content under the children element but there is enough room
// above, it displays content above child.
// It calls onRequestClose if the user clicks outside of content
//
// This component is used to easily build dropdowns & menus
const Popover = (props: {
  children: React.ReactElement,
  closeOnContentClick?: boolean,
  content: React.ReactNode,
  isOpen: boolean,
  className?: string,
  onClick?: () => void,
  onRequestClose?: () => void,
}) => {
  const targetRef = React.useRef<HTMLDivElement|null>(null)
  const contentRef = React.useRef<HTMLDivElement|null>(null)

  const windowSize = useWindowSize()
  const targetRect = useElementRect(props.isOpen ? targetRef.current : null)
  const contentRect = useElementRect(props.isOpen ? contentRef.current : null)

  useEventListenerUnlessRef(contentRef, 'mousedown', () => { if (props.onRequestClose) props.onRequestClose() })

  const onContentClick = () => {
    if (props.closeOnContentClick && props.onRequestClose) {
      setTimeout(props.onRequestClose)
    }
  }

  // Compute proper x/y location if it were to go over page boundaries
  const contentStyle: React.CSSProperties = {position: 'fixed'}
  if (props.isOpen && targetRect && contentRect) {
    if (targetRect.bottom + contentRect.height > windowSize.height && targetRect.top - contentRect.height > 0) {
      contentStyle.bottom = windowSize.height - targetRect.top
    } else {
      contentStyle.top = targetRect.bottom
    }
    if (targetRect.left + contentRect.width > windowSize.width && targetRect.right - contentRect.width > 0) {
      contentStyle.right = windowSize.width - targetRect.right
    } else {
      contentStyle.left = targetRect.left
    }
  }

  return (
    <>
      <div className={cx(props.className)} ref={targetRef} onClick={props.onClick}>{props.children}</div>
      {props.isOpen && createPortal((
        <div style={contentStyle} ref={contentRef} onClick={onContentClick}>{props.content}</div>
      ), document.body)}
    </>
  )
}

// A version of popover that opens automatically when `children` receives a click event
export const ClickPopover = (props: {
  content: React.ReactNode,
  children: React.ReactElement,
  closeOnContentClick?: boolean,
  className?: string,
}) => {
  const [isOpen, setIsOpen] = React.useState(false)

  return (
    <Popover
      className={cx(props.className)}
      children={props.children}
      closeOnContentClick={props.closeOnContentClick}
      content={props.content}
      isOpen={isOpen}
      onClick={() => setIsOpen(true)}
      onRequestClose={() => setIsOpen(false)}
    />
  )
}

export default Popover
