// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { useModal, renderModals } from 'src/helpers'
import classnames from 'classnames/bind'

import Button from 'src/components/button'
import Modal from 'src/components/modal'

const cx = classnames.bind(require('./stylesheet'))

// Provides a clickable area that spawns a modal designed for general help and keyboard
// shortcut information.
// usage:
//  <Help className={cx('help')} // Apply classnames to help style
//
//      // The preamble provides a space at the top of the help to provide a short paragraph on the page/component
//      preamble = "Review and Edit the accumulated evidence for this operation"
//
//      // shortcuts are typed via KeyboardShortcut. In general, each has a set of keys that produce the action,
//      // and a description of what that action is. Lowercase letter will be converted to uppercase characters (a-z only)
//      // and uppercase has "shift+" pre-pended to the combination
//      shortcuts = {
//        { keys: ["ArrowDown", "ArrowRight", "j"], description: "Move to the next evidence" }, // here, ArrowDown, ArrowRight, and j do the same thing
//        { keys: ["ArrowUp", "ArrowLeft", "k"], description: "Move to the previous evidence" },
//        { keys: ["g"], description: "Move to the top of the evidence list" }, // g is represented as `G`
//        { keys: ["G"], description: "Move to the bottom of the evidence list" }, // G is represented as `Shift + G`
//      }
//  />
export default (props: {
  preamble?: string,
  shortcuts?: Array<KeyboardShortcut>,
  className?: string,
}) => {
  const helpModal = useModal<void>(modalProps => <HelpModal preamble={props.preamble} shortcuts={props.shortcuts} {...modalProps} />)

  const onKeyDown = (e: KeyboardEvent) => {
    if (e.target == null || e.target !== document.body) return

    if (e.key === '?') helpModal.show()
  }

  React.useEffect(() => {
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  })

  return (
    <>
      <a className={cx(props.className)} onClick={() => helpModal.show()}>
        <img className={cx('help-icon')} src={require('./info-icon.svg')} />
      </a>
      {renderModals(helpModal)}
    </>
  )
}

export type KeyboardShortcut = {
  keys: Array<string>
  description: string
  uniqueKey?: string
}

const commonKeyboardShortcuts = [
  { keys: ["?"], description: "Open this help menu", uniqueKey: "stdhelp-?" },
  { keys: ["Escape"], description: "Close this window", uniqueKey: "stdhelp-Escape"},
]

export const HelpModal = (props: {
  preamble?: string,
  shortcuts?: Array<KeyboardShortcut>,
  onRequestClose: () => void,
}) => {
  const onKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') props.onRequestClose()
  }

  React.useEffect(() => {
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  })

  return (
    <Modal title="Help" onRequestClose={props.onRequestClose}>
      <div className={cx('preamble')}>{props.preamble}</div>
      {!props.shortcuts ? null :
        <div className={cx('shortcuts')}>
          <h1>Keyboard Shortcuts</h1>
          {props.shortcuts
            .concat(commonKeyboardShortcuts)
            .map(shortcut => <KeyboardShortcutKey key={shortcut.uniqueKey ? shortcut.uniqueKey : shortcut.keys[0]} shortcut={shortcut} />)
          }
        </div>
      }
      <Button primary onClick={props.onRequestClose}>Close</Button>
    </Modal>
  )
}

const KeyboardShortcutKey = (props: {
  shortcut: KeyboardShortcut
}) => {
  return (
    <div>
      <em className={cx('shortcut-description')}>{props.shortcut.description}: </em>
      {props.shortcut.keys
        .filter(key => key.length > 0)
        .map<React.ReactNode>((key, i) => {
          const singleLetter = key.length == 1
          const firstKeyCode = key.charCodeAt(0)
          const isUppercaseLetter = (65 <= firstKeyCode) && (firstKeyCode <= 90) // A - Z
          const isLowercaseLetter = (97 <= firstKeyCode) && (firstKeyCode <= 122) // a - z

          let content
          switch (true) {
            case singleLetter && isUppercaseLetter:
              content = <span><Key>Shift</Key>+<Key>{key}</Key></span>
              break
            case singleLetter && isLowercaseLetter:
              content = <Key>{key.toUpperCase()}</Key>
              break
            case key === ' ':
              content = <Key>Space</Key>
              break
            default:
              content = <Key>{key}</Key>
          }
          return <div key={i} className={cx('keygroup')}>{content}</div>
        })
        .reduce((prev, curr) => [prev, ', ', curr])
      }
    </div>
  )
}

const Key = (props: { children: React.ReactNode }) => <div className={cx('key')}>{props.children}</div>
