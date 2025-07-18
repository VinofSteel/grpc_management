FROM golang:1.24.5-alpine

# Install OS dependencies
RUN apk update && apk add --no-cache \
    git \
    make \
    postgresql-client

# Setting the work directory
WORKDIR /app

# Copying files to container
COPY . .

# Install Go dependencies and tools
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/air-verse/air@latest

# Set this directory as a safe directory for Git
RUN git config --global --add safe.directory /app

# Copy and make executable the start script
WORKDIR /app
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Exposing api port to the world
EXPOSE ${PORT}

# Running our api
CMD ["/start.sh"]