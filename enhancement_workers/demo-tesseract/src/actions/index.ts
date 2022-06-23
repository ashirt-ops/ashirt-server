export * from './processAction'

export type ActionResponse<T> = {
  result:
  | "Complete"
  | "Deferred"
  | "Error"
  | "Bad Request"
  | "Unhandled"
  | "Cannot process"
  data: T
}
