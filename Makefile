PKG := "github.com/wexinc/ps-tag-onboarding-go"
PKG_LIST := $(shell go list ${PKG}/...)
INTEGRATION_PKG_LIST := "./test/integration/..."
UNIT_PKG_LIST := $(shell go list ${PKG}/... | grep -v "${INTEGRATION_PKG_LIST}")

fmt:
	gofmt -w -s ${PKG_LIST}

init-checker:
	@go install honnef.co/go/tools/cmd/staticcheck@latest

unit-test:
	@go test -v ${UNIT_PKG_LIST}

integration-test:
	@go test -v ${INTEGRATION_PKG_LIST}

race:
	@go test -race ${UNIT_PKG_LIST}

benchmark:
	@go test -run="-" -bench=".*" ${UNIT_PKG_LIST}

lint:
	@staticcheck ${PKG_LIST}

vendor:
	@go mod vendor

check: init-checker lint
test: unit-test benchmark race