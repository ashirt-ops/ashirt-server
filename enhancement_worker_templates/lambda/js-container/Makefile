
.PHONY: build
build:
	docker build -t lambda-test .

.PHONY: run
run:
	docker run --rm -e "ASHIRT_ENV_VAR=some_value" -p 9000:8080 lambda-test

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
