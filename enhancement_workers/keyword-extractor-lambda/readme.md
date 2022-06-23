# (Image) Keyword Extractor

This worker is focused on parsing text from a screenshot, then running that data through a keyword extractor. This will produce results like:

```sh
password 0.1234
username 0.0295
```

Where the first value is the word/phrase, and the second value is the importance. Note that low values here represent more relevant keywords.

Under-the-hood, this uses [tesseract](https://github.com/tesseract-ocr/tesseract) and [yake](https://github.com/LIAAD/yake) to perform the work. This imposes a limitation in that only image content can be processed in this way.

Finally, note that this is a lambda-based worker, so is designed to work in an AWS environment.

## Deploying to AWS

See [Amazon's guide](https://docs.aws.amazon.com/lambda/latest/dg/images-create.html) on how to deploy this image.

## Deploying to AShirt

The below is a common sense configuration to deploy this worker within an ashirt environment.

```json
{
    "type": "aws",
    "version": 1,
    "lambdaName": "keyword-extractor",
    "asyncFunction": false
}
```

## Testing with Docker Compose / Standard development environment

The below is a reasonable docker-compose deployment _specifically for local development_:

```yaml
  keyword-extractor:
    build:
      context: enhancement_workers/keyword-extractor-lambda
      dockerfile: Dockerfile
    ports:
      - 3003:8080
    restart: on-failure
    environment:
      ASHIRT_BACKEND_URL: http://backend:3000
      # Note that these below values are pre-set in the standard database seed
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```
