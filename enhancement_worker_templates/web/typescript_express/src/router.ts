import { Express, Request, Response } from 'express'
import { handleEvidenceCreatedAction } from './actions'
import { config } from './config'
import { StatusCode } from './helpers/http_status_codes'
import { isSupportedMessage, EvidenceCreatedMessage } from './helpers/request_validation'
import { AShirtService } from './services/ashirt'


/**
 * addRoutes creates the routes necessary to support integration with the AShirt backend.
 * This also provides the area where more routes can be added to support other features.
 * 
 * @param app The express app hosting these routes
 * @param isDev If true, enables additional routes, for testing purposes
 */
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
          return await handleEvidenceCreatedRequest(body, res, ashirtService)
        default:
          res.status(StatusCode.NOT_IMPLEMENTED)
      }
    }
    res.status(StatusCode.NOT_IMPLEMENTED)
  })

  if (isDev) {
    console.log("==============> Adding Dev Routes <=================")
    // add Dev routes
    app.get('/', (_: Request, res: Response) => {
      res.send('Server live')
    })
    app.get('/test', async (_: Request, res: Response) => {
      try {
        res.status(200).send({ working: "yes" })
      }
      catch (err) {
        res.status(200).send({ message: "Test failed", err })
      }
    })
  }

}

/**
 * handleTestRequest processes test-type requests originating from the ashirt backend.
 * This is a canned response of "OK", but might need to be more detailed for your needs.
 * @param res 
 */
const handleTestRequest = (res: Response) => {
  res.status(StatusCode.OK).send({
    status: "ok",
  })
}

/**
 * handleProcessRequest processes process-type requests originating from the ashirt backend.
 * This defers the logic to actions/handleActionProcess and instead focuses on handling the
 * response structure, given the response from handleActionProcess
 * 
 * @param body The body of the process request
 * @param res The response
 * @param svc The AShirt service, which can be used to gather information on the evidence
 */
const handleEvidenceCreatedRequest = async (body: EvidenceCreatedMessage, res: Response, svc: AShirtService) => {
  const actionResponse = await handleEvidenceCreatedAction(body, svc)
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

