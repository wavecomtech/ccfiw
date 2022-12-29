
FROM golang:1.18-alpine AS development

ENV PROJECT_NAME=ccfiw
ENV PROJECT_PATH=/ccfiw
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$PROJECT_PATH/build:$GOBIN
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache make git bash alpine-sdk

RUN mkdir -p $PROJECT_PATH
WORKDIR $PROJECT_PATH
COPY go.mod go.sum $PROJECT_PATH/
RUN go mod download
COPY . $PROJECT_PATH/

RUN make build

FROM alpine:3.15.0 AS production

ENV PROJECT_NAME=ccfiw
RUN apk --no-cache add ca-certificates
COPY --from=development /$PROJECT_NAME/build/$PROJECT_NAME /usr/bin/$PROJECT_NAME
USER nobody:nogroup
ENTRYPOINT /usr/bin/$PROJECT_NAME
