# syntax=docker/dockerfile:1

FROM golang:1.20

# Set destination for COPY
WORKDIR /app

# Download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/app/.env /app/cmd/app/
# Copy the rest of the application source code to the container
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /stori_challenge cmd/app/main.go
# Set the entrypoint command to run the app when the container starts
CMD ["/stori_challenge","/app/client1.csv"]