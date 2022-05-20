# Enhancement Worker Templates (& Demo!)

This directory acts as a repository for service-based workers, as well as a demo worker for local testing. You can use this as a jumping off place for your own worker, or as inspiration for your own worker.

## Config for adding to AShirt Backend

The microservice-style apps can be added with the below change. The name given here is arbitrary.

### Docker Compose change

There are really two ways to add support for services to your docker-compose file. The first is the easy way, where we use the pre-configured headless user from the seed data to do the work. In practice, you'd likely want to have a unique headless user per app, but in low security environments -- i.e. development testing, this solution works well. To use this, simply use the yaml below, and add it to your docker-compose file.

```yml
  demo:
    build:
      context: enhancement_worker_templates/demo-tesseract
      dockerfile: Dockerfile.dev
    ports:
      - 3001:3001
    restart: on-failure
    environment:
      PORT: 3001
      ENABLE_DEV: true
      ASHIRT_BACKEND_URL: http://backend:3000
      ASHIRT_ACCESS_KEY: gR6nVtaQmp2SvzIqLUWdedDk
      ASHIRT_SECRET_KEY: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
```

If, however, your environment needs to be more secure, you'll need to pre-generate your key and secret, and update the above config, _then_ add the file to your docker_compose file, and restart the docker-compose server without destroying the database.

### Loading into AShirt

Service Name: Demo
Service Config:

```json
{
    "type": "web", 
    "version": 1,
    "url": "http://demo:3001/process"
}
```

### Permanently loading into AShirt seeding

If you want to keep this service around for awhile, you can define the docker-compose value above, and in addition, you can add the configuration and service name to the `hp_seed_data.go` (search for `newHPServiceWorker` to find examples on how to create your own)
