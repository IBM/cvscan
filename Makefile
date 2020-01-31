all: build

VERSION=$(shell git describe --tags | sed "s/\(-[[:digit:]]\{1,\}\)-[[:alnum:]]\{1,\}$$/\1/" | sed "s/-0$$//")
VER_PKG=$(shell go list -m)/pkg/version
LDFLAGS += -ldflags "-X ${VER_PKG}.version=${VERSION}"

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} ./cmd/...

.PHONY: install
install:
	CGO_ENABLED=0 go install ${LDFLAGS} ./cmd/...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: vet
	go test ./...

PUBLISH_BIN=cvscan

define targeted_build
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build ${LDFLAGS} ./cmd/${PUBLISH_BIN}
	$(eval OUTBIN=$(if $(filter $(1),windows),${PUBLISH_BIN}.exe,${PUBLISH_BIN}))
	tar -czf ${PUBLISH_BIN}-$(1)-$(2).tar.gz ${OUTBIN}
	rm ${OUTBIN}
endef

.PHONY: binaries
binaries:
	$(call targeted_build,darwin,amd64)
	$(call targeted_build,linux,amd64)
	$(call targeted_build,linux,ppc64le)
	$(call targeted_build,linux,s390x)
	$(call targeted_build,windows,386)
