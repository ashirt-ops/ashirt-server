// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import { Mutex } from 'async-mutex'

const mutex = new Mutex

// fetchJsonp loads data via jsonp (used for the archive viewer) by appending a global callback function to window
// and loading the data as a javascript file. The data must be in the form `jsonpFuncName({ dataToLoad })`.
//
// Since this function requires binding to a global callback on window that may not be unique per call, this function
// uses a mutex to call the "thread unsafe" variant of this function
export async function fetchJsonp(jsonpFuncName: string, url: string): Promise<any> {
  const release  = await mutex.acquire()
  try {
    return await fetchJsonpUnsafe(jsonpFuncName, url)
  } finally {
    release()
  }
}

const DATA_INITIAL_VALUE = {}

function fetchJsonpUnsafe(jsonpFuncName: string, url: string): Promise<any> {
  // @ts-ignore
  if (window[jsonpFuncName] != null) return Promise.reject(`window.${jsonpFuncName} is already set`)

  return new Promise((resolve, reject) => {
    let data: any = DATA_INITIAL_VALUE
    // @ts-ignore
    window[jsonpFuncName] = (d: any) => { data = d }

    const script = document.createElement('script')
    script.src = url
    script.onerror = () => reject(Error(`Failed to load script ${url}`))
    script.onload = () => {
      if (data === DATA_INITIAL_VALUE) {
        reject(Error(`Loaded script ${url} did not call jsonp function ${jsonpFuncName}`))
      } else {
        resolve(data)
      }
    }

    document.body.appendChild(script)
  }).finally(() => {
    // @ts-ignore
    delete window[jsonpFuncName]
  })
}
