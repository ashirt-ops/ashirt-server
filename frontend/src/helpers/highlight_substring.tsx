// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { escapeRegExp } from 'lodash'

/**
 * highlightSubstring breaks a given string into words that match the given regex, joined with
 * the rest of the string. This should preserve case.
 *
 * Note: this probably isn't close to the speediest solution. Use caution (and minLength option) when
 * using this with a large piece of text.
 *
 * @example
 * const result = highlightSubstring("The quick brown fox jumps over the lazy dog.", /the/gi, "highlight")
 * assert( result, [
 *   <span className="highlight">The</span>,
 *   <span> quick brown fox jumps over </span>,
 *   <span className="highlight">the</span>,
 *   <span> lazy dog.</span>,
 * ])
 *
 * @param s The string with a substring to highlight
 * @param regex What part of the string to match. Must be a global match (/.../g)
 * @param className What class name to apply to the highlighted word
 * @returns An array of spans. Spans will either be plain, or with the given classname.
 */
export const highlightSubstring = (s: string, regexAsStr: string, className: any,
  options?: {
    regexFlags?: string
    minLength?: number
  }
): Array<React.ReactNode> => {
  const rtn: Array<React.ReactNode> = []
  if (s === "" || regexAsStr.length < (options?.minLength ?? 1)) {
    return [<span>{s}</span>]
  }

  const matches = [...s.matchAll(new RegExp(escapeRegExp(regexAsStr), "g" + (options?.regexFlags ?? "")))]

  const endOfWord = (match: RegExpMatchArray) => (match.index ?? 0) + match[0].length
  const highlight = (v: string) => <span className={className}>{v}</span>

  if (matches.length) {
    if ((matches[0].index ?? 0) > 0) {
      rtn.push(<span>{s.substring(0, matches[0].index)}</span>)
    }

    for (let i = 0; i < matches.length; i++) {
      const item = matches[i]
      const next = matches[i + 1]
      const [value] = item
      rtn.push(highlight(value))
      if (next) {
        const end = endOfWord(item)
        const startOfNextWord = next.index ?? end
        if (end != startOfNextWord) {
          rtn.push(<span>{s.substring(end, startOfNextWord)}</span>)
        }
      }
    }
    const lastEntry = (matches[matches.length - 1])
    rtn.push(<span>{s.substring(endOfWord(lastEntry))}</span>)
  }
  else {
    rtn.push(<span>{s}</span>)
  }

  return rtn
}
