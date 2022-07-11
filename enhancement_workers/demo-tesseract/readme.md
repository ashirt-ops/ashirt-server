# Demo Worker -- Tesseract

This is a demo worker demonstrating how the enhancement workers operate. This particular worker executes tesseract on an incoming piece of image evidence (all others are ignored). While this might be useful for your purposes, you will likely want to spend time tuning your worker to get better tesseract results.

This worker iis adapted from the [Typescript Web Worker](/enhancement_worker_templates/web/typescript_express/readme.md). Details on how to modify the application can be found there.

## Typical Configuration

```ts
{
    "type": "web", 
    "version": 1,
    "url": "http://ts-tess/process"
}
```

## Testing with Docker Compose / Standard development environment

The following docker-compose configuration section should work:

```yml
  ts-tess:
    build:
      context: enhancement_workers/demo-tesseract
      dockerfile: Dockerfile.dev
    ports:
      - 3003:3003
    restart: on-failure
    volumes:
      - ./enhancement_workers/demo-tesseract/src:/app/src
    environment:
      ENABLE_DEV: true
      ASHIRT_BACKEND_URL: http://backend:3000
      PORT: 3003
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==

```
