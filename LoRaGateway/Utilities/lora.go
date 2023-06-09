package Utilities

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tarm/serial"
)

type LoRaGateway struct {
	serial_port string
	baud_rate   uint
	mqtt_broker string
	mqtt_port   uint

	LoRa *serial.Port

	queue_lora_receive chan []byte
	queue_mqtt_receive chan MQTT.Message

	mqtt_client             MQTT.Client
	mqtt_messagePubHandler  MQTT.MessageHandler
	mqtt_connectHandler     MQTT.OnConnectHandler
	mqtt_connectLostHandler MQTT.ConnectionLostHandler
}

func NewLoRaGateway(serial_port string, baud_rate uint, mqtt_broker string, mqtt_port uint) *LoRaGateway {

	lora_gateway := &LoRaGateway{
		serial_port: serial_port, baud_rate: baud_rate,
		mqtt_broker: mqtt_broker, mqtt_port: mqtt_port,
	}

	lora_gateway.mqtt_connectLostHandler = func(client MQTT.Client, err error) {
		log.Printf("WARNING: Connect lost to MQTT server: %v\n", err)
	}

	return lora_gateway
}

func (lg *LoRaGateway) Start(queue_capacity uint, topic_collect_sensor_data string, topic_subscribe map[string]byte) {
	config := &serial.Config{Name: lg.serial_port, Baud: int(lg.baud_rate)}
	LoRa, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	lg.LoRa = LoRa

	lg.queue_lora_receive = make(chan []byte, queue_capacity)
	lg.queue_mqtt_receive = make(chan MQTT.Message, queue_capacity)

	lg.mqtt_messagePubHandler = func(client MQTT.Client, msg MQTT.Message) {
		lg.queue_mqtt_receive <- msg
	}

	lg.mqtt_connectHandler = func(client MQTT.Client) {
		log.Printf("INFO: Connected to MQTT Broker tcp://%s:%d\n", lg.mqtt_broker, lg.mqtt_port)
		if token := client.SubscribeMultiple(topic_subscribe, nil); token.Wait() && token.Error() != nil {
			log.Printf("ERROR: MQTT Subscribe %s\n", token.Error())
		}
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", lg.mqtt_broker, lg.mqtt_port))
	opts.SetClientID("LoRa Gateway")
	opts.SetDefaultPublishHandler(lg.mqtt_messagePubHandler)
	opts.OnConnect = lg.mqtt_connectHandler
	opts.OnConnectionLost = lg.mqtt_connectLostHandler
	opts.AutoReconnect = true

	lg.mqtt_client = MQTT.NewClient(opts)
	if token := lg.mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connect %s\n", token.Error())
	}

	log.Printf("Start LoRa receive... %s\n", config.Name)

	go lg.forwardToMQTTBroker(topic_collect_sensor_data)
	go lg.mqttMessageHandle()

	defer close(lg.queue_lora_receive)
	for {
		lb := ReceiveBytes(LoRa, 2)
		if lb == nil {
			continue
		}
		length := binary.BigEndian.Uint16(lb)
		if length > 2048 {
			continue
		}
		payload := ReceiveBytes(LoRa, int(length))
		if payload == nil {
			continue
		}

		lg.queue_lora_receive <- payload
	}

}

func (lg *LoRaGateway) forwardToMQTTBroker(topic_collect_sensor_data string) {
	for {
		payload, ok := <-lg.queue_lora_receive
		if !ok {
			time.Sleep(time.Millisecond)
			continue
		}

		switch payload[0] {
		case SENSOR_DATA:
			mqtt_paylod := NewMQTTMessageJson(payload[1:])
			if mqtt_paylod == nil {
				continue
			}

			if token := lg.mqtt_client.Publish(topic_collect_sensor_data, 0, false, mqtt_paylod); token.Wait() && token.Error() != nil {
				log.Printf("ERROR: MQTT Publish %s\n", token.Error())
			}
			log.Println("SENSOR_DATA: ", string(mqtt_paylod))

		case ACK:
			go lg.handleLoRaACK(payload[1:])

		default:
			continue
		}

	}
}

func (lg *LoRaGateway) handleLoRaACK(ack []byte) {
	id := string(ack[:19])
	action := ""
	switch ack[19] {
	case CONTROL:
		action = "CONTROL-ACK"
	case STATUS:
		action = "CHECK-STATUS-ACK"
	case CONFIGURE:
		action = "CONFIGURE-ACK"
	}

	status := ""
	switch ack[20] {
	case PUMP_ON:
		action = "PUMP_ON"
	case PUMP_OFF:
		action = "PUMP_OFF"
	case LIGHT_ON:
		action = "LIGHT_ON"
	case LIGHT_OFF:
		action = "LIGHT_OFF"
	case CONFIGURE:
		action = "CONFIGURE_OK"
	}

	js_ack := fmt.Sprintf(`{"id": "%s", "action": "%s", "status": "%s"}`, id, action, status)
	if token := lg.mqtt_client.Publish("lora/ack", 0, false, []byte(js_ack)); token.Wait() && token.Error() != nil {
		log.Printf("ERROE: MQTT Publish %s\n", token.Error())
	}
	log.Println(js_ack)

}
