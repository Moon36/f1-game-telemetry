FROM golang:1.24.4 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/udp-server/*.go ./src/udp-server/
COPY src/udp-server/packets/*.go ./src/udp-server/packets/
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/udp-server/

FROM scratch
COPY --from=build /app/main /app/main
USER 1000:1000
EXPOSE ${PORT}
CMD ["/app/main"]