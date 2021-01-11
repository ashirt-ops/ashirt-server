import { Request, Response } from 'har-format'

export const mimetypeToAceLang = (mimetype: string) => {
  if (mimetype.includes("text/javascript") || mimetype.includes('application/json')) {
    return 'javascript'
  }
  else if (mimetype.includes('text/html')) {
    return 'html'
  }
  else if (mimetype.includes('text/css')) {
    return 'css'
  }
  else if (mimetype.includes('text/xml')) {
    return 'xml'
  }
  return ''
}

export const requestToRaw = (req: Request) => {
  const parsedUrl = new URL(req.url)
  const reqSummary = req.method + " " + parsedUrl.pathname + parsedUrl.search + " " + req.httpVersion + "\n"

  return reqSummary + req.headers.map(h => `${h.name}: ${h.value} `).join("\n")
}

export const responseToRaw = (resp: Response) => {
  return resp.headers.map(h => `${h.name}: ${h.value} `).join("\n")
}
