package Utilities

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
)

func ReceiveBytes(rd io.Reader, size int) []byte {
	if size < 0 {
		return nil
	}

	reader := bufio.NewReader(rd)
	result := make([]byte, size)

	for i := 0; i < size; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return nil
		}
		result[i] = b
	}
	return result
}

type LoRaMessageJson struct {
	ID string  `json:"ID"`
	TP float32 `json:"TP"`
	HU float32 `json:"HU"`
	HI float32 `json:"HI"`
	TD float32 `json:"TD"`
	EC float32 `json:"EC"`
	SG string  `json:"SG"`
	PK string  `json:"PK"`
}

type MQTTMessageJson struct {
	Temperature float32 `json:"Temperature"`
	Humidity    float32 `json:"Humidity"`
	Heat_Index  float32 `json:"HeatIndex"`
	TDS_Value   float32 `json:"TDSValue"`
	EC_Value    float32 `json:"ECValue"`
}

type LoRaPayload struct {
	ID string  `json:"ID"`
	TP float32 `json:"TP"`
	HU float32 `json:"HU"`
	HI float32 `json:"HI"`
	TD float32 `json:"TD"`
	EC float32 `json:"EC"`
}

func VerifySignature(lora_message *LoRaMessageJson, secret_key []byte) bool {
	payload := LoRaPayload{
		ID: lora_message.ID, TP: lora_message.TP, HU: lora_message.HU,
		HI: lora_message.HI, TD: lora_message.TD, EC: lora_message.EC,
	}

	js_payload, err := json.Marshal(payload)
	if err != nil {
		log.Println("ERROR:", err)
		return false
	}

	verify := HMAC_SHAKE256(js_payload, secret_key)

	return verify == lora_message.SG
}

func NewMQTTMessageJson(payload []byte) []byte {
	var lrjs LoRaMessageJson
	err := json.Unmarshal(payload, &lrjs)
	if err != nil {
		log.Println("ERROR:", err)
		return nil
	}

	public_key, err := hex.DecodeString(lrjs.PK)
	if err != nil {
		log.Printf("ERROR: public_key must be represented by the hexadecimal string %s\n", err)
		return nil
	}

	secret_key := GetSecretKey(public_key, lrjs.ID)
	if secret_key == nil {
		log.Printf("WARNING: Unknown device ID %s\n", lrjs.ID)
		return nil
	}

	ok := VerifySignature(&lrjs, secret_key)
	if !ok {
		log.Printf("WARNING: No-Trust device (%s) because Verify Signature is failed.\n", lrjs.ID)
		return nil
	} else {
		log.Printf("INFO: Trust device (%s) because Verify Signature is passed.\n", lrjs.ID)
	}

	mqtt_js := MQTTMessageJson{
		Temperature: lrjs.TP, Humidity: lrjs.HU,
		Heat_Index: lrjs.HI, TDS_Value: lrjs.TD,
		EC_Value: lrjs.TD / 707.0,
	}

	result, err := json.Marshal(mqtt_js)
	if err != nil {
		log.Println("ERROR:", err)
		return nil
	}
	return result
}
