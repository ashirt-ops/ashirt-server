import * as React from 'react'

// Returns a ClientRect ({top, left, bottom, right, width, height}) for an optional
// HTMLElement that updates automatically on window resize/scroll
//
// This is useful for portals that want to position themselves over a target area
export function useElementRect(el: React.MutableRefObject<HTMLElement|null>|null): ClientRect|null {
  const [rect, setRect] = React.useState<ClientRect|null>(null)

  React.useLayoutEffect(() => {
    if (el == null) return
    const updateRect = () => {
      window.requestAnimationFrame(() => {
        if (el != null && el.current != null) setRect(el.current.getBoundingClientRect())
      })
    }
    updateRect()
    window.addEventListener('resize', updateRect)
    window.addEventListener('scroll', updateRect, true)
    return () => {
      window.removeEventListener('resize', updateRect)
      window.removeEventListener('scroll', updateRect, true)
    }
  }, [el])

  return rect
}
