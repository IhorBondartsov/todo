FROM golang:1.17-alpine  AS builder

WORKDIR /src
RUN go env
COPY go.mod ./
COPY go.sum ./
RUN  go mod download
COPY . ./
RUN ls
ARG  GOFLAGS
ARG  VERSION
RUN $GOFLAGS go build -ldflags "-X main.Version=$(VERSION) -s -w" -o /src/bin/todo-app cmd/main.go


# Build runtime image.
FROM       alpine
COPY       --from=builder /src/bin /src/app
ENTRYPOINT ["/src/app/todo-app"]