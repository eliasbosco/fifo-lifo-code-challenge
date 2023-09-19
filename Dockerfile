FROM golang:1.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build

# Run the tests in the container
FROM builder AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/unicorn .
COPY --from=builder /app/*txt .

EXPOSE 8888

ENTRYPOINT [ "./unicorn"]