# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache update
RUN apk --no-cache upgrade
RUN apk --no-cache add build-base git mercurial gcc
ADD . /src
RUN cd /src && go build -o main ./cmd/bus

FROM alpine
WORKDIR /app
COPY --from=build-env /src/main /app/
ENTRYPOINT ./main
