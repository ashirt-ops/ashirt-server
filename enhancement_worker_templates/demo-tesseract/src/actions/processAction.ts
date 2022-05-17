import { default as axios, AxiosError } from 'axios'
import tesseract from 'node-tesseract-ocr'
import { Logger } from 'pino'

import { AShirtService } from "src/services/ashirt"
import { ProcessRequest } from "src/helpers/request_validation"

export type ProcessResultDTO =
  | ProcessResultNormal
  | ProcessResultComplete
  | ProcessResultDeferred

type ProcessResultNormal = {
  action: "rejected" | "error"
  content?: string
}

type ProcessResultComplete = {
  action: "processed"
  content: string
}

type ProcessResultDeferred = {
  action: "deferred"
}

export const handleActionProcess = async (
  body: ProcessRequest,
  svc: AShirtService,
  reqLog: Logger,
): Promise<ProcessResultDTO> => {

  if (body.contentType !== 'image') {
    return {
      action: 'rejected',
    }
  }

  try {
    const resp = await svc.getEvidenceContent(body.operationSlug, body.evidenceUuid)
    const content = await tesseract.recognize(resp.data)

    return {
      action: 'processed',
      content
    }
  }
  catch (err: unknown) {
    reqLog.error({ err }, "Unable to process evidence")
    const content = axios.isAxiosError(err)
      ? (err as AxiosError).message
      : null

    const rtn: ProcessResultNormal = {
      action: 'error',
      ...(content ? { content } : {}),
    }
    return rtn
  }
}
