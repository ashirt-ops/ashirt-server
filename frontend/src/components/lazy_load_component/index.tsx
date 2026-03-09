import { useState, useEffect, useRef, type ReactNode } from 'react'

export default function LazyLoadComponent(props: { children: ReactNode }) {
  const [isVisible, setIsVisible] = useState(false)
  const containerRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true)
          observer.disconnect()
        }
      },
      { root: null, rootMargin: '0px', threshold: 0.1 },
    )

    if (containerRef.current) {
      observer.observe(containerRef.current)
    }

    const el = containerRef.current
    return () => {
      if (el) {
        observer.unobserve(el)
      }
    }
  }, [])

  return (
    <div ref={containerRef} style={{ height: '100%' }}>
      {isVisible ? props.children : null}
    </div>
  )
}
