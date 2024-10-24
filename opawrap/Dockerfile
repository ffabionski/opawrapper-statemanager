# Build the application from source
FROM golang:latest AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /opawrap

# Run the tests in the container
FROM build AS test
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build /opawrap /opawrap

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/opawrap"]