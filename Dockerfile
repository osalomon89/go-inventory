FROM golang:1.19-alpine

# Set destination for COPY
WORKDIR /app

ENV GO_ENVIRONMENT=production
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_NAME=inventory
ENV DB_USER=root
ENV DB_PASS=secret

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code.
COPY . .

# Build
RUN go build -o /api

# This is for documentation purposes only.
# To actually open the port, runtime parameters
EXPOSE 8080

# Run
CMD [ "/api" ]