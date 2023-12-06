const getValue = (key: string, defaultValue: string) => {
  const value = localStorage.getItem(key)

  return (value == null) ? defaultValue : value
}

const prefIncludeDeletedUsers = "admin-includeDeletedUsers"

export const getIncludeDeletedUsers = (): boolean => {
  return getValue(prefIncludeDeletedUsers, false.toString()) === "true"
}

export const setIncludeDeletedUsers = (value: boolean) => {
  localStorage.setItem(prefIncludeDeletedUsers, value.toString())
}
