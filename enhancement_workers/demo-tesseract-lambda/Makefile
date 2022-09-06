
.PHONY: build
build:
	docker build -t lambda-test .

.PHONY: run
run:
	docker run --rm --env-file ./.env -p 9000:8080 lambda-test


.PHONY: start-msg
start-msg:
	@printf "\033[0;32mStarting Container...\033[0m\n"

.PHONY: start
start: build start-msg run

.PHONY: test-test
test-test:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "test"}' | jq | pygmentize -l json

.PHONY: test-process-img
test-process-img:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "evidence_created", "evidenceUuid": "seed_dursleys", "operationSlug": "HPSS", "contentType": "image"}' | jq | pygmentize -l json

.PHONY: test-unsupported
test-unsupported:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"type": "unsupported"}' | jq | pygmentize -l json
