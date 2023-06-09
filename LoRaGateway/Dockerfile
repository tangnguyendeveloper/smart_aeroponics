FROM golang:1.20

WORKDIR /gateway

COPY go.mod go.sum ./
RUN go mod download

COPY * ./
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /lora-gateway

CMD [ "/lora-gateway" ]