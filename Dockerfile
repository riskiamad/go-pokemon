# Start from golang base image
FROM golang:1.17-alpine as builder

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

# Set the current working directory inside the container 
RUN mkdir /app

COPY . /app

WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o myPokemonService .

# Start a new stage from scratch
FROM alpine:latest as deploy

RUN mkdir /app

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/myPokemonService /app/myPokemonService
COPY --from=builder /app/.env /app/.env       

#Command to run the executable
CMD ["./app/myPokemonService"]