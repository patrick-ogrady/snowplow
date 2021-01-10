.PHONY: build install add-license check-license compile

ADDLICENSE_CMD=go run github.com/google/addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "patrick-ogrady" -l "mit" -v

build:
	go build ./...

install:
	go install ./...

add-license:
	${ADDLICENCE_SCRIPT} .;
	go mod tidy;

check-license:
	${ADDLICENCE_SCRIPT} -check .;
	go mod tidy;

compile:
	./scripts/compile.sh $(version)
