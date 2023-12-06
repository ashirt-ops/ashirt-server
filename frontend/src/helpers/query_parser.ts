import * as dateFns from 'date-fns'

function stringifyQuery(tokens: { [key: string]: Array<FilterValue> }): string {
  const query = []

  for (let key in tokens) {
    for (let token of tokens[key]) {
      let part = ''
      if (key !== '') part = `${key}:`
      if (token.modifier == 'not') {
        part += "!"
      }
      if (token.value.indexOf(' ') === -1) part += token.value
      else part += `"${token.value}"`
      query.push(part)
    }
  }

  return query.join(' ')
}

export type FilterModifier = "not" | undefined

export type FilterValue = {
  value: string
  modifier?: FilterModifier
}

export type FilterModified<T> = T & {
  modifier?: FilterModifier
}

export type ParsedQuery = { [key: string]: Array<FilterValue> }

export function parseQuery(query: string): ParsedQuery {
  const tokenized: { [key: string]: Array<FilterValue> } = {}
  let modifier: FilterModifier
  let currentToken = ''
  let inQuote = false
  let currentKey = ''

  const pushToken = () => {
    if (currentToken.length === 0) return
    if (tokenized[currentKey] == null) tokenized[currentKey] = []
    tokenized[currentKey].push({
      modifier,
      value: currentToken,
    })
    currentToken = ''
    currentKey = ''
    modifier = undefined
  }

  for (let i = 0; i < query.length; i++) {
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

      // modifiers
      case '!':
        if (currentKey != "" && currentToken == "") {
          modifier = "not"
          continue
        }
    }
    currentToken += query[i]
  }
  pushToken()
  return tokenized
}

export function addTagToQuery(query: string, tagToAdd: string): string {
  return addFieldToQuery(query, "tag", tagToAdd)
}

export function addOperatorToQuery(query: string, userSlugToAdd: string): string {
  return addFieldToQuery(query, "operator", userSlugToAdd)
}

function addFieldToQuery(query: string, field: string, value: string, negate?: boolean): string {
  const tokenized = parseQuery(query)
  tokenized[field] = (tokenized[field] ?? [])
  if (tokenized[field].findIndex(item => item.value === value) !== -1) {
    return query
  }
  tokenized[field].push({ value: value, modifier: negate ? "not" : undefined })
  return stringifyQuery(tokenized)
}

export function getDateRangeFromQuery(query: string): [Date, Date] | null {
  const tokenized = parseQuery(query)
  if (!tokenized.range) {
    return null
  }
  return parseDateRangeString(tokenized.range[0].value)
}

export function parseDateRangeString(dateRangeStr: string): [Date, Date] | null {
  const [from, to] = dateRangeStr.split(',').map(str => dateFns.parseISO(str))
  if (!from || !to || !dateFns.isValid(from) || !dateFns.isValid(to)) {
    return null
  }
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
    tokenized.range = [{ value: stringifiedRange.join(',') }]
  }
  return stringifyQuery(tokenized)
}
