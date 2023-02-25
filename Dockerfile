
# Build backend
FROM golang:1.20-alpine as build-backend

ADD . /src/maestro
WORKDIR /src/maestro
RUN go build -o bin/maestro -ldflags "-X 'github.com/yukitsune/maestro.Version=$VERSION'" cmd/maestro/main.go

# Build frontend
FROM node:16-bullseye-slim as build-frontend

RUN mkdir -p /opt/maestro/frontend
WORKDIR /opt/maestro/frontend

COPY web/frontend-remix .
RUN npm install
RUN npm run build

# Run
FROM node:19-alpine3.16

# Copy backend
RUN mkdir /opt/maestro
COPY --from=build-backend /src/maestro/bin/maestro /opt/maestro/maestro
COPY --from=build-backend /src/maestro/assets /opt/maestro/assets

# Copy frontend
RUN mkdir /opt/maestro/frontend
COPY --from=build-frontend /opt/maestro/frontend /opt/maestro/frontend

# Install concurrently
RUN npm install -g concurrently

# Versioning information
ARG GIT_COMMIT
ARG GIT_BRANCH=main
ARG GIT_DIRTY='false'
ARG VERSION
LABEL branch=$GIT_BRANCH \
    commit=$GIT_COMMIT \
    dirty=$GIT_DIRTY \
    version=$VERSION

# Backend env vars
ENV MAESTRO_API_ASSETS_DIR=/opt/maestro/assets

# Frontend env vars
ENV NODE_ENV=production
ENV API_URL="http://localhost:8182"

WORKDIR /opt/maestro

CMD ["concurrently", "'./maestro serve'", "'npm start --prefix ./frontend'"]