package main

import "LoRaGateway/Utilities"

func main() {

	const (
		com_port    = "/dev/ttyUSB0"
		baud_rate   = 9600
		mqtt_broker = "mqtt-broker"
		mqtt_port   = 1883

		queue_capacity = 521

		topic_collect_sensor_data = "sensors/aeroponics"
	)

	topic_subscribe := map[string]byte{
		"pump/config": 1, "pump/control": 1, "pump/status": 1,
		"lights/control": 1, "lights/status": 1,
	}

	lora_gateway := Utilities.NewLoRaGateway(com_port, baud_rate, mqtt_broker, mqtt_port)

	lora_gateway.Start(queue_capacity, topic_collect_sensor_data, topic_subscribe)
}
