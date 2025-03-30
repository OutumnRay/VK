FROM golang:1.24

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go mod tidy && go build -o mattermost-bot ./cmd/main.go

CMD ["./mattermost-bot"]
