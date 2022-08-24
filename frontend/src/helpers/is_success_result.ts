
import { ErrorResult, Result, SuccessfulResult } from 'src/global_types'

export const isSuccessfulResult = <T>(v: Result<T>): v is SuccessfulResult<T> => {
  return ('success' in v)
}

export const isErrorResult = <T>(v: Result<T>): v is ErrorResult => {
  return ('err' in v)
}

type ResultState = 
  | "success"
  | "unresolved"
  | "error"

export const getResultState = <T>(v: Result<T> | null): ResultState => {
  if (v === null) {
    return 'unresolved'
  }
  if (isSuccessfulResult(v)) {
    return 'success'
  }
  return 'error'
}
