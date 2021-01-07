
// trimURL shorts a full url to it's base components. This particular method focuses on showing
// the domain, starting path, and ending path (omitting middle directories and query)
// this works best with Github-like urls.
// Example:
// https://github.com/microsoft/vscode/blob/master/README.md => github.com/microsoft/vscode/.../README.md
export const trimURL = (original: string, maxNonUrlLength: number = 50, keepBeginingParts: number = 2): {
  isAUrl: boolean
  trimmedValue: string
} => {
  const urlRegex = /^([^:]+):\/\/([^/]+)(?:\/([^?]+)\??(.*)?)/i
  const urlComponents = original.match(urlRegex)

  if (!urlComponents) {
    return {
      isAUrl: false,
      trimmedValue: original.length <= maxNonUrlLength ? original : original.substr(0, maxNonUrlLength - 3) + "...",
    }
  }

  const [_, _proto, domain, path, _query] = urlComponents
  let pathParts = path.split("/")
  let lastOne = pathParts[pathParts.length - 1]
  const hashIndex = lastOne.indexOf('#')
  lastOne = lastOne.substr(0, hashIndex > -1 ? hashIndex : lastOne.length)
  pathParts[pathParts.length - 1] = lastOne

  const totalKeepParts = keepBeginingParts + 1

  const label = domain + (pathParts.length <= totalKeepParts
    ? `/${pathParts.join('/')}`
    : `/${pathParts.slice(0, keepBeginingParts).join('/')}/.../${lastOne}`)

  return {
    isAUrl: true,
    trimmedValue: label,
  }
}
