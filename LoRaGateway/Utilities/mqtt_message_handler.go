package Utilities

import (
	"encoding/hex"
	"log"
	"strings"
	"time"
)

// TYPE || Sig || PK || MQTT payload == ID || msg

func (lg *LoRaGateway) mqttMessageHandle() {
	for {
		msg, ok := <-lg.queue_mqtt_receive
		if !ok {
			time.Sleep(time.Millisecond)
			continue
		}

		n := len(msg.Payload())

		if msg.Topic() == "pump/status" || msg.Topic() == "lights/status" {
			n = 22
		}

		if n < 20 {
			continue
		}

		log.Println(msg.Topic(), string(msg.Payload()))

		switch msg.Topic() {
		case "pump/config":
			pk, sk := GetPublicKey(string(msg.Payload()[:19]))
			message := append([]byte(strings.ToUpper(hex.EncodeToString(pk))), msg.Payload()...)
			signature := HMAC_SHAKE256(message, sk)
			message = append([]byte(signature), message...)

			_, err := lg.LoRa.Write(append([]byte{CONFIGURE}, message...))
			if err != nil {
				log.Printf("ERROR: Send CONFIGURE, %s\n", err)
			}

		case "pump/control", "lights/control":
			pk, sk := GetPublicKey(string(msg.Payload()[:19]))
			message := append([]byte(strings.ToUpper(hex.EncodeToString(pk))), msg.Payload()...)
			signature := HMAC_SHAKE256(message, sk)
			message = append([]byte(signature), message...)

			_, err := lg.LoRa.Write(append([]byte{CONTROL}, message...))
			if err != nil {
				log.Printf("ERROR: Send CONTROL, %s\n", err)
			}

		case "pump/status", "lights/status":
			// pk, sk := GetPublicKey(string(msg.Payload()[:19]))
			// message := append([]byte(strings.ToUpper(hex.EncodeToString(pk))), msg.Payload()...)
			// signature := HMAC_SHAKE256(message, sk)
			// message = append([]byte(signature), message...)

			_, err := lg.LoRa.Write(append([]byte{STATUS}, msg.Payload()...))
			if err != nil {
				log.Printf("ERROR: Send STATUS, %s\n", err)
			}

		default:
			continue
		}
	}
}
