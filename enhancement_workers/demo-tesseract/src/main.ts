import express, { Express } from 'express'
import { config } from './config'

import { addRoutes } from './router'
import { appInfoLog } from './services/logging'

async function bootstrap() {
  const app: Express = express()
  app.use(express.json())
  
  addRoutes(app, config().enableDev)
  
  const port = config().port
  app.listen(port, () => {
    appInfoLog(`server started on port ${port}`)
  })
}
bootstrap()
