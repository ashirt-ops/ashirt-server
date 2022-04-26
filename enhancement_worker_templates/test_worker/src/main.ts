import express, { Express } from 'express'
import { config } from './config'

import { addRoutes } from './router'

async function bootstrap() {
  const app: Express = express()
  app.use(express.json())

  addRoutes(app, config().enableDev)
  
  const port = config().port
  app.listen(port, () => {
    console.log(`server started on port ${port}`)
  })
}
bootstrap()
