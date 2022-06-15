imageNameDev = samplePython-dev
imageName = samplePython
hostPort = 9000

.PHONY: build_dev
build_dev:
	docker build -t $(imageNameDev) -f Dockerfile.dev .

.PHONY: build_release
build_release:
	docker build -t $(imageName) -f Dockerfile.prod .

run_dev:
	docker-compose up --build

run_prod:
	docker-compose -f docker-compose.prod.yml up --build
