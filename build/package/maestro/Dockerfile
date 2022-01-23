# Build
FROM golang:1.16-alpine as build

ADD . /go/src/github.com/YuKitsune/maestro
WORKDIR /go/src/github.com/YuKitsune/maestro
RUN go build -o bin/maestro -ldflags "-X 'github.com/yukitsune/maestro.Version=$VERSION'" cmd/maestro/main.go

# Run
FROM alpine:3.15.0

# Versioning information
ARG GIT_COMMIT
ARG GIT_BRANCH=main
ARG GIT_DIRTY='false'
ARG VERSION
LABEL branch=$GIT_BRANCH \
    commit=$GIT_COMMIT \
    dirty=$GIT_DIRTY \
    version=$VERSION

COPY --from=build /go/src/github.com/YuKitsune/maestro/bin/maestro maestro
COPY --from=build /go/src/github.com/YuKitsune/maestro/assets /assets

ENV MAESTRO_API_ASSETS_DIR=/assets

CMD  ["./maestro", "serve"]