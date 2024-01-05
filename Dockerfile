FROM golang:latest
WORKDIR /builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main main.go


FROM alpine:latest
WORKDIR /app
COPY --from=0 /builder/main .
COPY config.yaml .
CMD ["./main"]