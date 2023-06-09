FROM golang:1.20

WORKDIR /gateway

COPY go.mod go.sum ./
RUN go mod download

RUN go get github.com/tangnguyendeveloper/smart_aeroponics/LoRaGateway/Utilities

COPY ./*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /lora-gateway

CMD [ "/lora-gateway" ]