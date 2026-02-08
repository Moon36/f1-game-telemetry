FROM golang:1.24.4 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/database/ ./src/database/
COPY src/internal/ ./src/internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/database/

FROM scratch
COPY --from=build /app/main /app/main
USER 1000:1000
CMD ["/app/main"]