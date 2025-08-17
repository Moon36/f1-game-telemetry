FROM golang:1.24.4 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch
COPY --from=build /app/main /app/main
USER 1000:1000
EXPOSE ${PORT}
CMD ["/app/main"]