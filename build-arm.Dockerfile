FROM golang
WORKDIR /app
COPY *.go ./
COPY go.mod go.sum ./
RUN go mod download
RUN go build