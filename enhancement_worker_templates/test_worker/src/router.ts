import { Express, Request, Response } from 'express'
import { handleActionProcess } from './actions'
import { config } from './config'
import { StatusCode } from './helpers/http_status_codes'
import { isSupportedRequest, ProcessRequest } from './helpers/request_validation'
import { AShirtService } from './ashirt'

export function addRoutes(app: Express, isDev: boolean) {
  const ashirtService = new AShirtService(
    config().backendUrl,
    config().accessKey,
    config().secretKeyB64
  )

  app.post('/process', async (req: Request, res: Response) => {
    const body = req.body

    if (isSupportedRequest(body)) {
      switch (body.type) {
        case 'test':
          return handleTestRequest(res)
        case 'process':
          return await handleProcessRequest(body, res, ashirtService)
        default:
          res.sendStatus(StatusCode.NOT_IMPLEMENTED)
      }
    }
    res.sendStatus(StatusCode.NOT_IMPLEMENTED)
  })

  if (isDev) {
    console.log("==============> Adding Dev Routes <=================")
    // add Dev routes
    app.get('/', (_: Request, res: Response) => {
      res.send('Server live')
    })
    app.get('/test', async (_: Request, res: Response) => {
      try {
        // const data = await ashirtService.getEvidenceContent("HPSS", "seed_dursleys")
        // res.status(200).send(data)
        await handleProcessRequest({
          contentType: 'image',
          evidenceUuid: 'seed_dursleys',
          operationSlug: 'HPSS',
          type: 'process'
        }, res, ashirtService)

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

const handleProcessRequest = async (body: ProcessRequest, res: Response, svc: AShirtService) => {
  const actionResponse = await handleActionProcess(body, svc)
  switch (actionResponse.action) {
    case 'deferred':
      res.sendStatus(StatusCode.ACCEPTED)
      return
    case 'error':
      actionResponse.content
        ? res.status(StatusCode.INTERNAL_SERVICE_ERROR).send(actionResponse)
        : res.sendStatus(StatusCode.INTERNAL_SERVICE_ERROR)
      return
    case 'rejected':
      actionResponse.content
        ? res.status(StatusCode.NOT_ACCEPTABLE).send(actionResponse)
        : res.sendStatus(StatusCode.NOT_ACCEPTABLE)
      return
    case 'processed':
      res.status(StatusCode.OK).send(actionResponse)
      return
    default:
      res.sendStatus(StatusCode.INTERNAL_SERVICE_ERROR)
  }
}

