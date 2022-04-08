export function isError(error: unknown): error is Error {
  return (
    !!error &&
    typeof error === "object" &&
    typeof (error as Error).message === "string"
  )
}
