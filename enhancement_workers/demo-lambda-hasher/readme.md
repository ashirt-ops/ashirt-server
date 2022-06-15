# Demo Worker -- Hasher

This is a demo worker demonstrating how the enhancement workers operate. This particular worker simply provides a sha256 hash of the content. This is likely not useful for anything other than demonstration purposes.

This worker is adapted from the [Javascript lambda worker](/enhancement_worker_templates/lambda/js-container/readme.md). Details on how to modify the application can be found there.

## Typical Configuration

```ts
{
    "type": "aws",
    "version": 1,
    "lambdaName": "hasher",
    "asyncFunction": false
}
```

## Testing with Docker Compose / Standard development environment

The following docker-compose configuration section should work:

```yml
  hasher:
    build:
      context: enhancement_workers/demo-lambda-hasher
      dockerfile: Dockerfile
    ports:
      - 3002:8080
    restart: on-failure
    environment:
      ASHIRT_BACKEND_URL: backend
      ASHIRT_BACKEND_PORT: 3000
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```
