.PHONY: install add-license check-license

ADDLICENSE_CMD=go run github.com/google/addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "patrick-ogrady" -l "mit" -v

install:
	go install ./...

add-license:
	${ADDLICENCE_SCRIPT} .;
	go mod tidy;

check-license:
	${ADDLICENCE_SCRIPT} -check .;
	go mod tidy;
