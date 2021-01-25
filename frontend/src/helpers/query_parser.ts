// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as dateFns from 'date-fns'

function stringifyQuery(tokens: {[key: string]: Array<string>}): string {
  const query = []

  for (let key in tokens) {
    for (let token of tokens[key]) {
      let part = ''
      if (key !== '') part = `${key}:`
      if (token.indexOf(' ') === -1) part += token
      else part += `"${token}"`
      query.push(part)
    }
  }

  return query.join(' ')
}

export function parseQuery(query: string): {[key: string]: Array<string>} {
  const tokenized: {[key: string]: Array<string>} = {}
  let currentToken = ''
  let inQuote = false
  let currentKey = ''

  const pushToken = () => {
    if (currentToken.length === 0) return
    if (tokenized[currentKey] == null) tokenized[currentKey] = []
    tokenized[currentKey].push(currentToken)
    currentToken = ''
    currentKey = ''
  }

  for (let i = 0; i < query.length; i++ ) {
    switch (query[i]) {
      case ' ':
        if (inQuote) break
        pushToken()
        continue
      case '"':
        inQuote = !inQuote
        continue
      case ':':
        if (currentKey !== '') break
        currentKey = currentToken
        currentToken = ''
        continue
    }
    currentToken += query[i]
  }
  pushToken()
  return tokenized
}

export function addTagToQuery(query: string, tagToAdd: string): string {
  const tokenized = parseQuery(query)
  if (!tokenized.tag) tokenized.tag = []
  else if (tokenized.tag.indexOf(tagToAdd) !== -1) return query
  tokenized.tag.push(tagToAdd)
  return stringifyQuery(tokenized)
}

export function addOperatorToQuery(query: string, userSlugToAdd: string): string {
  const tokenized = parseQuery(query)
  tokenized.operator = [userSlugToAdd]
  return stringifyQuery(tokenized)
}

export function getDateRangeFromQuery(query: string): [Date, Date] | null {
  const tokenized = parseQuery(query)
  if (!tokenized.range) return null
  const [from, to] = tokenized.range[0].split(',').map(str => dateFns.parseISO(str))
  if (!from || !to || !dateFns.isValid(from) || !dateFns.isValid(to)) return null
  return [dateFns.startOfDay(from), dateFns.endOfDay(to)]
}

export function addOrUpdateDateRangeInQuery(query: string, range: [Date, Date] | null): string {
  const tokenized = parseQuery(query)
  if (range == null) {
    tokenized.range = []
  } else {
    const stringifiedRange = range.map((d, i) => {
      const useShorthand = dateFns.isEqual(
        d,
        i === 0 ? dateFns.startOfDay(d) : dateFns.endOfDay(d),
      )
      return useShorthand ? dateFns.format(d, 'yyyy-MM-dd') : d.toISOString()
    })
    tokenized.range = [stringifiedRange.join(',')]
  }
  return stringifyQuery(tokenized)
}
