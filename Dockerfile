FROM golang:1.25.5 AS build

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /nakup

FROM scratch

WORKDIR /app

COPY --from=build /nakup /app/nakup

EXPOSE 8081

ENTRYPOINT [ "/app/nakup" ]