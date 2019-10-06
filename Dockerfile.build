FROM golang:1.12 as builder
# Force the go compiler to use modules
ENV GO111MODULE=on
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o main .
RUN ls

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main ./
CMD ["/app/main"]