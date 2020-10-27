// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import * as ace from 'ace-builds'
import AceEditor from 'react-ace'
import classnames from 'classnames/bind'
import { useElementRect } from 'src/helpers'
const cx = classnames.bind(require('./stylesheet'))

// This is a wrapper component around react-ace to fix some bugs
// and make it easier to work with. It handles automatic loading of modes
// with webpack chunks and it matches the size of the editor to the size
// of the parent container.
export default (props: {
  mode: string,
  onChange?: (code: string) => void,
  readOnly?: boolean,
  value: string,
}) => {
  const rootRef = React.useRef<HTMLDivElement|null>(null)
  const editorSize = useSizeOfParentContainer(rootRef)
  const mode = useLoadAceModeWithWebpack(props.mode)
  useStopPropagationOfSearchKeydowns(rootRef)

  return (
    <div className={cx('root')} ref={rootRef}>
      <AceEditor
        {...props}
        {...editorSize}
        mode={mode}
        theme="ashirt"
        editorProps={{ $blockScrolling: true }}
      />
    </div>
  )
}

// Use webpack codesplitting to break all mode-* files in ace-builds/src-noconflict
// into webpack chunks and load them as requested. This function will return the mode
// string only after it has loaded the mode via webpack chunk
function useLoadAceModeWithWebpack(requestedMode: string): string {
  const [loadedMode, setLoadedMode] = React.useState('plain_text')

  React.useEffect(() => {
    if (requestedMode === '') {
      setLoadedMode('plain_text')
      return
    }
    import(`ace-builds/src-noconflict/mode-${requestedMode}`)
      .then(() => setLoadedMode(requestedMode))
      .catch(err => console.error(`Unable to load mode: ${requestedMode}`))
  }, [requestedMode])

  return loadedMode
}

// By matching the size of the parent component we can use the editor like
// any other html component without having to worry about setting
// maxlines/minlines
//
// We do this by setting the pixel size of <AceEditor /> to the pixel size of its
// root container. Using values "100%" for both cause strange bugs where the editor
// may not display if a min-height/min-width is specified in a parent rather than
// an absolute height/width
function useSizeOfParentContainer(parentRef: React.MutableRefObject<HTMLDivElement|null>): {width: string, height: string} {
  const [size, setSize] = React.useState({ width: '100%', height: '100%' })

  const parentRect = useElementRect(parentRef)
  React.useEffect(() => {
    if(parentRect == null) {
      return
    }
    setSize({
      width: (`${parentRect.width}px`),
      height: (`${parentRect.height}px`),
    })
  }, [parentRect])
  return size
}

// Prevents keydown events in the search field from propagating up to react since there
// are a few places where we listen to keydown events higher up in the dom e.g. timeline.
// Since search field is self-contained within the ace-editor it seems reasonable that
// a user typing in the search field shouldn't emit keydown events up the dom.
function useStopPropagationOfSearchKeydowns(parentRef: React.MutableRefObject<HTMLDivElement|null>) {
  const onKeyDown = (e: KeyboardEvent) => {
    if (e.target && (e.target as HTMLElement).className === 'ace_search_field') {
      e.stopPropagation()
    }
  }
  React.useEffect(() => {
    const curParentRef = parentRef.current
    if (!curParentRef) return
    curParentRef.addEventListener('keydown', onKeyDown)
    return () => { curParentRef.removeEventListener('keydown', onKeyDown) }
  })
}

// Theme is defined in acetheme.styl but this prevents ace from
// trying to load a null theme which will be blocked by CSP
//
// @ts-ignore - typescript doesn't know about ace.define
ace.define("ace/theme/ashirt", [], function(_, exports) {
  exports.isDark = true
  exports.cssClass = 'ace-ashirt-theme'
})

// Import plain_text to prevent ace from attempting to download this causing CSP errors
//
// Other modes are loaded on the fly by calling `import()` in useEffect in the component
// itself. By using import() webpack will break each mode out into its own chunk,
// reducing the total bundle size of each chunk
import 'ace-builds/src-noconflict/mode-plain_text'
import 'ace-builds/src-noconflict/ext-searchbox'
