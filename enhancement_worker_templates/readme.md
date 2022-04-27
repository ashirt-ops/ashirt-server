# Enhancement Worker Templates

This directory acts as a repository for service-based workers, as well as a test worker for local testing. You can use this as a jumping off place for your own worker, or as inspiration for your own worker.

## Config for adding to AShirt Backend

I prefer the name `Hexer` for this project, so if you don't like that, you'll need to rename the below.

### Docker Compose change

Note that the Access key and secret key will need to be changed. Currently, the process is:
1. Provide the hexer config below
2. Start the backend, and create a new API key
3. Update the docker-compose area below with the new key
4. stop docker-compose (don't destroy the containers)
5. restart docker-compose

```yml
  hexer:
    build:
      context: enhancement_worker_templates/test_worker
      dockerfile: Dockerfile
    ports:
      - 3001:3001
    restart: on-failure
    environment:
      PORT: 3001
      ENABLE_DEV: true
      ASHIRT_BACKEND_URL: http://backend:3000
      ASHIRT_ACCESS_KEY: 2tMRBB1-XsNzxrbpYca-BBPq
      ASHIRT_SECRET_KEY: mP/2rfq7czA+Kxhj5n46HJhRsCMKW0GuYlZ3Wq/w4jcRZDkN1S7Xi6Lz9+RM5ydAJCg7FNJ77kWOuWOKvXHmyw==
```

### Loading into AShirt

Service Name: Hexer
Service Config:

```json
{
    "type": "web", 
    "version": 1,
    "url": "http://hexer:3001/process"
}
```
