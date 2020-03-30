PACKAGE_NAME = osbuild-installer-backend

.PHONY: build
build:
	go build -o osbuild-installer ./cmd/osbuild-installer/
	go build -o osbuild-installer-openapi ./cmd/osbuild-installer-openapi/

# pip3 install openapi-spec-validator
.PHONY: check-api-spec
check-api-spec:
	 openapi-spec-validator osbuild_installer/openapi/api.spec.yaml

# go get github.com/deepmap/oapi-codegen
.PHONY: generate-api
generate-api:
	oapi-codegen osbuild_installer/openapi/api.spec.yaml > internal/api/api.go

# https://github.com/OpenAPITools/openapi-generator
# generates stub, copy out the useful stuff
.PHONY: generate-api-openapitools
generate-api-openapitools:
	podman pull openapitools/openapi-generator-cli
	mkdir -p tmp/openapitools
	mkdir -p internal/api-openapitools
	podman run --rm -v $(shell pwd):/local openapitools/openapi-generator-cli generate \
		-i /local/osbuild_installer/openapi/api.spec.yaml \
		-g go-server -o /local/tmp/openapitools
	cp tmp/openapitools/go/routers.go internal/api-openapitools/
	cp tmp/openapitools/go/model_version.go internal/api-openapitools/
	cp tmp/openapitools/go/logger.go internal/api-openapitools/
	cp tmp/openapitools/go/api.go internal/api-openapitools/
	cp tmp/openapitools/go/api_default.go internal/api-openapitools/
	cp tmp/openapitools/go/api_default_service.go internal/api-openapitools/
	[ -z "${SUDO_USER}" ] || chown -R ${SUDO_USER}:${SUDO_USER} internal/api-openapitools
	rm -rf tmp

