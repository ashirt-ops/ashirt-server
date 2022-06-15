import { default as express, Express } from 'express'
import { config } from './config'

import { addRoutes } from './router'

async function bootstrap() {
  const app: Express = express()

  // expect json responses
  app.use(express.json())

  // Set routes to handle ashirt requests
  addRoutes(app, config().enableDev)
  
  // start the server
  const port = config().port
  app.listen(port, () => {
    console.log(`server started on port ${port}`)
  })
}
bootstrap()
