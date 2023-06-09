FROM golang:1.20

WORKDIR $GOPATH/src/LoRaGateway

COPY . .
RUN go mod download

RUN go build -o /lora-gateway

CMD [ "/lora-gateway" ]