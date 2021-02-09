import { Har, Log, Creator, Browser, Page, PageTiming, Entry, Request, Response, Cookie, Header, QueryString, Content, PostData, Param, Cache, CacheDetails, Timings } from 'har-format'

export const isAHar = (o: any): o is Har => hasField("log", o, isHarLog)

const isHarLog = (o: any): o is Log => {
  console.log("isVersion?", hasString("version", o))
  console.log("creator?", hasField("creator", o, isHarCreator))
  console.log("browser?", hasMaybeField("browser", o, isHarBrowser))
  console.log("pages?", hasMaybeArray("pages", o, isHarPage))
  console.log("entries?", hasArray("entries", o, isHarEntry))
  console.log("comment?", hasMaybeString("comment", o))

  return (
    hasString("version", o)
    && hasField("creator", o, isHarCreator)
    && hasMaybeField("browser", o, isHarBrowser)
    && hasMaybeArray("pages", o, isHarPage)
    && hasArray("entries", o, isHarEntry)
    && hasMaybeString("comment", o)
  )
}

const isHarCreator = (o: any): o is Creator => {
  return (
    hasString("name", o)
    && hasString("version", o)
    && hasMaybeString("comment", o)
  )
}

const isHarBrowser = (o: any): o is Browser => {
  return (
    hasString("name", o)
    && hasString("version", o)
    && hasMaybeString("comment", o)
  )
}

const isHarPage = (o: any): o is Page => {
  return (
    hasString("startedDateTime", o)
    && hasString("id", o)
    && hasString("title", o)
    && hasField("pageTimings", o, isPageTiming)
    && hasMaybeString("comment", o)
  )
}

const isPageTiming = (o: any): o is PageTiming => {
  return (
    hasMaybeNumber("onContentLoad", o)
    && hasMaybeNumber("onLoad", o)
    && hasMaybeString("comment", o)
  )
}

const isHarEntry = (o: any): o is Entry => {
  return (
    hasMaybeString("pageref", o)
    && hasString("startedDateTime", o)
    && hasNumber("time", o)
    && hasField("request", o, isHarRequest)
    && hasField("response", o, isHarResponse)
    && hasField("cache", o, isHarEntryCache)
    && hasField("timings", o, isHarEntryTimings)
    && hasMaybeString("serverIPAddress", o)
    && hasMaybeString("connection", o)
    && hasMaybeString("comment", o)
  )
}

const isHarRequest = (o: any): o is Request => {
  return (
    hasString("method", o)
    && hasString("url", o)
    && hasString("httpVersion", o)
    && hasArray("cookies", o, isHarCookie)
    && hasArray("headers", o, isHarHeader)
    && hasArray("queryString", o, isHarQueryString)
    && hasMaybeField("postData", o, isHarPostData)
    && hasNumber("headersSize", o)
    && hasNumber("bodySize", o)
    && hasMaybeString("comment", o)
  )
}

const isHarResponse = (o: any): o is Response => {
  return (
    hasNumber("status", o)
    && hasString("statusText", o)
    && hasString("httpVersion", o)
    && hasArray("cookies", o, isHarCookie)
    && hasArray("headers", o, isHarHeader)
    && hasField("content", o, isHarContent)
    && hasString("redirectURL", o)
    && hasNumber("headersSize", o)
    && hasNumber("bodySize", o)
    && hasMaybeString("comment", o)
  )
}

const isHarCookie = (o: any): o is Cookie => {
  return (
    hasString("name", o)
    && hasString("value", o)
    && hasMaybeString("path", o)
    && hasMaybeString("domain", o)
    && hasMaybeField("expires", o, (t) => t == null || isAString(t))
    && hasMaybeBoolean("httpOnly", o)
    && hasMaybeBoolean("secure", o)
    && hasMaybeString("comment", o)
  )
}

const isHarHeader = (o: any): o is Header => {
  return (
    hasString("name", o)
    && hasString("value", o)
    && hasMaybeString("comment", o)
  )
}

const isHarQueryString = (o: any): o is QueryString => {
  return (
    hasString("name", o)
    && hasString("value", o)
    && hasMaybeString("comment", o)
  )
}

const isHarPostData = (o: any): o is PostData => {
  return (
    hasString("mimeType", o)
    && hasMaybeString("text", o)
    && hasMaybeArray("param", o, isHarParam)
    && hasMaybeString("comment", o)
  )
}

const isHarParam = (o: any): o is Param => {
  return (
    hasString("name", o)
    && hasMaybeString("value", o)
    && hasMaybeString("fileName", o)
    && hasMaybeString("contentType", o)
    && hasMaybeString("comment", o)
  )
}

const isHarContent = (o: any): o is Content => {
  return (
    hasNumber("size", o)
    && hasMaybeNumber("compression", o)
    && hasString("mimeType", o)
    && hasMaybeString("text", o)
    && hasMaybeString("encoding", o)
    && hasMaybeString("comment", o)
  )
}

const isHarEntryCache = (o: any): o is Cache => {
  return (
    hasMaybeField("beforeRequest", o, isHarCacheDetails)
    && hasMaybeField("afterRequest", o, isHarCacheDetails)
    && hasMaybeString("comment", o)
  )
}

const isHarCacheDetails = (o: any): o is (CacheDetails | null) => {
  return o === null || (
    hasMaybeString("expires", o)
    && hasString("lastAccess", o)
    && hasString("eTag", o)
    && hasNumber("hitCount", o)
    && hasMaybeString("comment", o)
  )
}

const isHarEntryTimings = (o: any): o is Timings => {
  return (
    hasMaybeNumber("blocked", o)
    && hasMaybeNumber("dns", o)
    && hasMaybeNumber("connect", o)
    && hasMaybeNumber("send", o)
    && hasNumber("wait", o)
    && hasNumber("receive", o)
    && hasMaybeNumber("ssl", o)
    && hasMaybeString("comment", o)
  )
}

const hasField = (field: string, o: any, ofType: (b: any) => boolean) => {
  return field in o && ofType(o[field])
}

const hasMaybeField = (field: string, o: any, ofType: (b: any) => boolean): boolean => {
  return (!(field in o) || hasField(field, o, ofType))
}

const hasArray = (field: string, o: any, isChildOfType: (b: any) => boolean): boolean => {
  console.log(field, "has field?", field in o)
  console.log(field, "isArray?", Array.isArray(o[field]))
  console.log(field, "valid children?", o[field].map((item: any) => isChildOfType(item)).reduce((acc: boolean, cur: boolean) => acc && cur, true))

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
const isABoolean = (o: any): o is boolean => typeof o === 'boolean'

const hasString = (field: string, o: any) => hasField(field, o, isAString)
const hasMaybeString = (field: string, o: any) => hasMaybeField(field, o, isAString)
const hasNumber = (field: string, o: any) => hasField(field, o, isANumber)
const hasMaybeNumber = (field: string, o: any) => hasMaybeField(field, o, isANumber)
// const hasBoolean = (field: string, o: any) => hasField(field, o, isABoolean)
const hasMaybeBoolean = (field: string, o: any) => hasMaybeField(field, o, isABoolean)
