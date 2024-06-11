# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /
COPY go.mod go.sum ./
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /clerk

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
COPY --from=build-stage /clerk /clerk
COPY --from=build-stage config.yaml ./

ENTRYPOINT ["/clerk", "-l", "debug"]
