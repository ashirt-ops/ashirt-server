
.PHONY: build
build:
	docker build -t lambda-test .

.PHONY: run
run:
	docker run --rm \
	-e "ASHIRT_BACKEND_URL=host.docker.internal" \
	-e "ASHIRT_BACKEND_PORT=3000" \
	-e "ASHIRT_ACCESS_KEY=gR6nVtaQmp2SvzIqLUWdedDk" \
	-e "ASHIRT_SECRET_KEY=WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==" \
	--add-host host.docker.internal:host-gateway \
	-p 9000:8080 lambda-test

.PHONY: start-msg
start-msg:
	@printf "\033[0;32mStarting Container...\033[0m\n"

.PHONY: start
start: build start-msg run

.PHONY: test-test
test-test:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "test"}' | json_pp | pygmentize -l json

.PHONY: test-process-img
test-process-img:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "evidence_created", "evidenceUuid": "seed_dursleys", "operationSlug": "HPSS", "contentType": "image"}' | json_pp | pygmentize -l json

.PHONY: test-unsupported
test-unsupported:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "unsupported"}' | json_pp | pygmentize -l json
