# Use base go image just to build the executable file which is then passed to alpine for making it lightweight (i.e. no need go dependencies can just run the app from alpine)
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

# Set env var
# NOT using any C library just Go standard lib
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# Just to make sure that it's executable
RUN chmod +x /app/brokerApp

# Build a NEW tiny docker image (different from above)
FROM alpine:latest

RUN mkdir /app

# Copy the go build file from builder to the NEW alpine image
COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]