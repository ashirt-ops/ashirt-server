import { AShirtService } from "src/services/ashirt"
import { EvidenceCreatedMessage } from "src/helpers/request_validation"
import { default as axios, AxiosError } from 'axios'

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

/**
 * handleEvidenceCreatedAction is the location where actual request processing takes place. The
 * implementation here is very basic and focuses on the form, rather than the function.
 * 
 * @param body The message received from the AShirt backend, in the ProcessRequest type
 * @param svc The AShirt service, which can be used to get the evidence content, if needed
 * @returns The result of the processing, which will be further interpreted back in router.ts
 */
export const handleEvidenceCreatedAction = async (
  body: EvidenceCreatedMessage,
  svc: AShirtService
): Promise<ProcessResultDTO> => {

  // Restrict the content types being processed -- here we require that all processable content
  // be images
  if (body.contentType !== 'image') {
    return {
      action: 'rejected',
    }
  }

  try {
    // Your logic goes here

    return {
      action: 'processed',
      content: 'done!', // replace me with your content!
    }
  }
  catch (err: unknown) {
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
