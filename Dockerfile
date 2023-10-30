# Compilation: docker build -t gobeats .
# Run: docker run -it -v ${PWD}:/app gobeats

FROM golang:1.21.3

RUN apt-get update && apt-get install -y libasound2-dev

WORKDIR /app

VOLUME /app
