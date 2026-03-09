class HttpError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

type QueryObj = Record<string, unknown>

export default async function xhr(
  method: string,
  path: string,
  data?: object | null,
  query?: QueryObj,
) {
  return request((res) => res.json(), method, path, data, query)
}

export async function xhrText(
  method: string,
  path: string,
  data?: object | null,
  query?: QueryObj,
) {
  return request((res) => res.text(), method, path, data, query)
}

function getErrorMessage(body: unknown): string | undefined {
  if (body != null && typeof body === 'object' && 'error' in body) {
    return (body as { error: string }).error
  }
}

async function request<T>(
  decode: (res: Response) => Promise<T>,
  method: string,
  path: string,
  data?: object | null,
  query?: QueryObj,
): Promise<T> {
  path = '/web' + path
  if (query != null) path += `?${new URLSearchParams(query as Record<string, string>).toString()}`
  let res
  if (method === 'GET') {
    res = await fetch(path, { method })
  } else {
    const body = JSON.stringify(data)
    const headers = {
      'Content-Type': 'application/json',
    }
    res = await fetch(path, { method, body, headers })
  }

  const responseJson = await decode(res)
  const errorMessage = getErrorMessage(responseJson)
  if (res.status < 200 || res.status >= 300 || errorMessage) {
    throw new HttpError(res.status, errorMessage ?? '')
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
