export function computeDelta<T>(before: Array<T>, after: Array<T>): [Array<T>, Array<T>] {
  const additions = after
    .filter(item => !before.includes(item))

  const subtractions = before
    .filter(item => !after.includes(item))

  return [additions, subtractions]
}
