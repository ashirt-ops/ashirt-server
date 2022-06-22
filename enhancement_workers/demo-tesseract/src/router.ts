import { Express, Request, Response } from 'express'
import { Logger } from 'pino'

import { handleActionProcess } from './actions'
import { config } from './config'
import { StatusCode } from './helpers/http_status_codes'
import { isSupportedMessage, EvidenceCreatedMessage } from './helpers/request_validation'
import { AShirtService } from './services/ashirt'
import { appWarnLog, getRequestLogger } from './services/logging'

export function addRoutes(app: Express, isDev: boolean) {
  const ashirtService = new AShirtService(
    config().backendUrl,
    config().accessKey,
    config().secretKeyB64
  )

  app.post('/process', async (req: Request, res: Response) => {
    const body = req.body

    if (isSupportedMessage(body)) {
      switch (body.type) {
        case 'test':
          return handleTestRequest(res)
        case 'evidence_created':
          const logger = getRequestLogger()
          logger.info(body, "Process request")
          try {
            const result = await handleProcessRequest(body, res, ashirtService, logger)
            logger.info({ result }, "Request complete")
            return result
          }
          catch (err) {
            logger.info({ err }, "Request failed")
            res.status(StatusCode.INTERNAL_SERVICE_ERROR).send({
              "error": "unable to process message"
            })
          }
        default:
          res.status(StatusCode.NOT_IMPLEMENTED)
      }
    }
    res.status(StatusCode.NOT_IMPLEMENTED)
  })

  if (isDev) {
    appWarnLog("Adding Dev Routes")
    // add Dev routes
    app.get('/', (_: Request, res: Response) => {
      res.send('Server live')
    })
    app.get('/test', async (_: Request, res: Response) => {
      try {
        const content = { message: "No test set. Write one!" }
        res.status(200).send({ content })
      }
      catch (err) {
        res.status(200).send({ message: "Test failed", err })
      }

    })
  }

}

const handleTestRequest = (res: Response) => {
  res.status(StatusCode.OK).send({
    status: "ok",
  })
}

const handleProcessRequest = async (body: EvidenceCreatedMessage, res: Response, svc: AShirtService, logger: Logger) => {
  const actionResponse = await handleActionProcess(body, svc, logger)
  switch (actionResponse.action) {
    case 'deferred':
      res.status(StatusCode.ACCEPTED)
      return
    case 'error':
      actionResponse.content
        ? res.status(StatusCode.INTERNAL_SERVICE_ERROR).send(actionResponse)
        : res.status(StatusCode.INTERNAL_SERVICE_ERROR)
      return
    case 'rejected':
      actionResponse.content
        ? res.status(StatusCode.NOT_ACCEPTABLE).send(actionResponse)
        : res.status(StatusCode.NOT_ACCEPTABLE)
      return
    case 'processed':
      res.status(StatusCode.OK).send(actionResponse)
      return
    default:
      res.status(StatusCode.INTERNAL_SERVICE_ERROR)
  }
}

