import {
  Har, Log, Creator, Entry, Request, Response, Header, Content, PostData, Param,
  // Browser, Cache, CacheDetails, Cookie, Page, PageTiming, QueryString, Timings
} from 'har-format'

// Note that a lot of comments exist here because the har format was not finalized, and not everyone
// encodes them the same way. Everything below (mostly) conforms to the har types, but is ignored
// since it gets in the way.

export const isAHar = (o: any): o is Har => hasField("log", o, isHarLog)

const isHarLog = (o: any): o is Log => {
  return (
    hasArray("entries", o, isHarEntry) // used by viewer
    && hasField("creator", o, isHarCreator) // used by viewer

    // && hasString("version", o)
    // && hasMaybeField("browser", o, isHarBrowser)
    // && hasMaybeArray("pages", o, isHarPage)
    // && hasMaybeString("comment", o)
  )
}

const isHarCreator = (o: any): o is Creator => {
  return (
    hasString("name", o) // used by viewer
    && hasString("version", o) // used by viewer
    // && hasMaybeString("comment", o)
  )
}

// const isHarBrowser = (o: any): o is Browser => {
//   return (
//     hasString("name", o)
//     && hasString("version", o)
//     && hasMaybeString("comment", o)
//   )
// }

// const isHarPage = (o: any): o is Page => {
//   return (
//     hasString("startedDateTime", o)
//     && hasString("id", o)
//     && hasString("title", o)
//     && hasField("pageTimings", o, isPageTiming)
//     && hasMaybeString("comment", o)
//   )
// }

// const isPageTiming = (o: any): o is PageTiming => {
//   return (
//     hasMaybeNumber("onContentLoad", o)
//     && hasMaybeNumber("onLoad", o)
//     && hasMaybeString("comment", o)
//   )
// }

const isHarEntry = (o: any): o is Entry => {
  return (
    hasNumber("time", o) // used by viewer
    && hasField("request", o, isHarRequest) // used by viewer
    && hasField("response", o, isHarResponse) // used by viewer
    && hasMaybeString("serverIPAddress", o) // used by viewer
    // && hasMaybeString("pageref", o)
    // && hasString("startedDateTime", o)
    // && hasField("cache", o, isHarEntryCache)
    // && hasField("timings", o, isHarEntryTimings)
    // && hasMaybeString("connection", o)
    // && hasMaybeString("comment", o)
  )
}

const isHarRequest = (o: any): o is Request => {
  return (
    hasString("method", o) // used by viewer
    && hasString("url", o) // used by viewer
    && hasArray("headers", o, isHarHeader) // used by viewer
    && hasMaybeField("postData", o, isHarPostData) // used by viewer
    && hasString("httpVersion", o) // used by viewer
    // && hasArray("cookies", o, isHarCookie)
    // && hasArray("queryString", o, isHarQueryString)
    // && hasNumber("headersSize", o)
    // && hasNumber("bodySize", o)
    // && hasMaybeString("comment", o)
  )
}

const isHarResponse = (o: any): o is Response => {
  return (
    hasNumber("status", o) // used by viewer
    && hasArray("headers", o, isHarHeader) // used by viewer
    && hasField("content", o, isHarContent) // used by viewer
    // && hasString("statusText", o)
    // && hasString("httpVersion", o)
    // && hasArray("cookies", o, isHarCookie)
    // && hasString("redirectURL", o)
    // && hasNumber("headersSize", o)
    // && hasNumber("bodySize", o)
    // && hasMaybeString("comment", o)
  )
}

// const isHarCookie = (o: any): o is Cookie => {
//   return (
//     hasString("name", o)
//     && hasString("value", o)
//     && hasMaybeString("path", o)
//     && hasMaybeString("domain", o)
//     && hasMaybeField("expires", o, (t) => t == null || isAString(t))
//     && hasMaybeBoolean("httpOnly", o)
//     && hasMaybeBoolean("secure", o)
//     && hasMaybeString("comment", o)
//   )
// }

const isHarHeader = (o: any): o is Header => {
  return (
    hasString("name", o)
    && hasString("value", o)
    && hasMaybeString("comment", o)
  )
}

// const isHarQueryString = (o: any): o is QueryString => {
//   return (
//     hasString("name", o) // used by viewer
//     && hasString("value", o) // used by viewer
//     && hasMaybeString("comment", o)
//   )
// }

const isHarPostData = (o: any): o is PostData => {
  return (
    hasString("mimeType", o)
    && hasMaybeString("text", o)
    && hasMaybeArray("param", o, isHarParam)
    // && hasMaybeString("comment", o)
  )
}

const isHarParam = (o: any): o is Param => {
  return (
    hasString("name", o)
    && hasMaybeString("value", o)
    // && hasMaybeString("fileName", o)
    // && hasMaybeString("contentType", o)
    // && hasMaybeString("comment", o)
  )
}

const isHarContent = (o: any): o is Content => {
  return (
    hasNumber("size", o)  // used by viewer
    && hasMaybeString("text", o)  // used by viewer
    && hasString("mimeType", o)  // used by viewer
    // && hasMaybeNumber("compression", o)
    // && hasMaybeString("encoding", o)
    // && hasMaybeString("comment", o)
  )
}

// const isHarEntryCache = (o: any): o is Cache => {
//   return (
//     hasMaybeField("beforeRequest", o, isHarCacheDetails)
//     && hasMaybeField("afterRequest", o, isHarCacheDetails)
//     && hasMaybeString("comment", o)
//   )
// }

// const isHarCacheDetails = (o: any): o is (CacheDetails | null) => {
//   return o === null || (
//     hasMaybeString("expires", o)
//     && hasString("lastAccess", o)
//     && hasString("eTag", o)
//     && hasNumber("hitCount", o)
//     && hasMaybeString("comment", o)
//   )
// }

// const isHarEntryTimings = (o: any): o is Timings => {
//   return (
//     hasMaybeNumber("blocked", o)
//     && hasMaybeNumber("dns", o)
//     && hasMaybeNumber("connect", o)
//     && hasMaybeNumber("send", o)
//     && hasNumber("wait", o)
//     && hasNumber("receive", o)
//     && hasMaybeNumber("ssl", o)
//     && hasMaybeString("comment", o)
//   )
// }

const hasField = (field: string, o: any, ofType: (b: any) => boolean) => {
  return field in o && ofType(o[field])
}

const hasNullField = (field: string, o: any) => {
  return field in o && o[field] == null
}

const hasMaybeField = (field: string, o: any, ofType: (b: any) => boolean, notNull: boolean = false): boolean => {
  return (!(field in o) || hasField(field, o, ofType)) || (!notNull && hasNullField(field, o))
}

const hasArray = (field: string, o: any, isChildOfType: (b: any) => boolean): boolean => {
  return (
    field in o
    && Array.isArray(o[field])
    && o[field].map((item: any) => isChildOfType(item)).reduce((acc: boolean, cur: boolean) => acc && cur, true)
  )
}

const hasMaybeArray = (field: string, o: any, ofType: (b: any) => boolean): boolean => {
  return (!(field in o) || hasArray(field, o, ofType))
}

const isAString = (o: any): o is string => typeof o === 'string'
const isANumber = (o: any): o is number => typeof o === 'number'
// const isABoolean = (o: any): o is boolean => typeof o === 'boolean'

const hasString = (field: string, o: any) => hasField(field, o, isAString)
const hasMaybeString = (field: string, o: any) => hasMaybeField(field, o, isAString)
const hasNumber = (field: string, o: any) => hasField(field, o, isANumber)
// const hasMaybeNumber = (field: string, o: any) => hasMaybeField(field, o, isANumber)
// const hasMaybeBoolean = (field: string, o: any) => hasMaybeField(field, o, isABoolean)
