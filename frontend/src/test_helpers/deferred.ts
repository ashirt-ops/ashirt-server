// A Deferred is a promise that can be resolved/rejected externally
export type Deferred<T> = {
  resolve: (value: T) => void,
  reject: (err: Error) => void,
  promise: Promise<T>,
}

export function createDeferred<T>(): Deferred<T> {
  let resolve: (value: T) => void = () => {}
  let reject: (err: Error) => void = () => {}
  return {
    resolve(value: T) { setTimeout(() => resolve(value)) },
    reject(err: Error) { setTimeout(() => reject(err)) },
    promise: new Promise<T>((res, rej) => { resolve = res; reject = rej }),
  }
}
