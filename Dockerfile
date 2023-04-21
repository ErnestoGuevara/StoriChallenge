# syntax=docker/dockerfile:1

FROM golang:1.20

# Set destination for COPY
WORKDIR /app
COPY test ./test
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code. Note the slash at the end, as explained in
COPY *.go ./
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /stori_challenge
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8080
#Run
CMD ["/stori_challenge"]