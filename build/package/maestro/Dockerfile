FROM golang:1.16-alpine

# Versioning information
ARG GIT_COMMIT
ARG GIT_BRANCH=main
ARG GIT_DIRTY='false'
ARG VERSION
LABEL branch=$GIT_BRANCH \
    commit=$GIT_COMMIT \
    dirty=$GIT_DIRTY \
    version=$VERSION

ADD . /go/src/github.com/YuKitsune/maestro
WORKDIR /go/src/github.com/YuKitsune/maestro
RUN go build -o bin/maestro -ldflags "-X 'github.com/yukitsune/maestro.Version=$VERSION'" cmd/maestro/main.go
CMD ["./bin/maestro", "serve"]