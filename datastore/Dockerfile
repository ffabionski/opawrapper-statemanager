# Build the application from source
FROM golang:latest AS build

WORKDIR /app

COPY go.* ./
COPY *.go ./
COPY *.json /*.json

RUN CGO_ENABLED=0 GOOS=linux go build -o /datastore

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS relase

WORKDIR /

COPY --from=build /*.json /*.json
COPY --from=build /datastore /datastore

EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/datastore"]