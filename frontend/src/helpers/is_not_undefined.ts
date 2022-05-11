// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

/**
 * isNotUndefined is a dual purpose check/type guard that will ensure that the given value
 * is not _exactly_ undefined. This is most useful for filter, which normally won't be able
 * to deduce that the result set does not contain undefined values
 * @param t 
 * @returns 
 */
export function isNotUndefined<T>(t: T | undefined): t is T {
  return t !== undefined
}
