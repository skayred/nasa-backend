FROM registry.semaphoreci.com/golang:1.18 as builder

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

RUN mkdir -p /src/projects
WORKDIR /src/projects
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o app

FROM registry.semaphoreci.com/golang:1.18
COPY --from=builder /src/projects/app /bin/app

ENTRYPOINT ["/bin/app"]
