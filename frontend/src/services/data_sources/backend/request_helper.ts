import queryString from "query-string"

class HttpError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

type QueryObj = Record<string, any>

export default async function xhr(method: string, path: string, data?: Object | null, query?: QueryObj) {
  return request(res => res.json(), method, path, data, query)
}

export async function xhrText(method: string, path: string, data?: Object | null, query?: QueryObj) {
  return request(res => res.text(), method, path, data, query)
}

async function request(decode: (res: Response) => Promise<any>, method: string, path: string, data?: Object | null, query?: QueryObj) {
  path = '/web' + path
  if (query != null) path += `?${queryString.stringify(query)}`
  let res;
  if (method === 'GET') {
    res = await fetch(path, { method })
  } else {
    const body = JSON.stringify(data);
    const headers = {
      'Content-Type': 'application/json',
    }
    res = await fetch(path, { method, body, headers })
  }

  const responseJson = await decode(res)
  if (res.status < 200 || res.status >= 300 || (responseJson && responseJson.error)) {
    throw new HttpError(res.status, responseJson.error)
  }
  return responseJson
}

export async function reqMultipart(method: string, path: string, body: FormData) {
  path = '/web' + path
  const res = await fetch(path, { method, body })
  const responseJson = await res.json()
  if (res.status < 200 || res.status >= 300 || (responseJson && responseJson.error)) {
    throw new HttpError(res.status, responseJson.error)
  }
  return responseJson
}
