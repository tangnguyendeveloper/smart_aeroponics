#ifndef UTILITIES_HPP
#define UTILITIES_HPP

#include <Arduino.h>
#include "DHT.h"
#include <SoftwareSerial.h>
#include <ArduinoJson.h>

#define DHTPIN 2
#define DHTTYPE DHT22

#define TdsSensorPin A2
#define VREF 5.0f
#define SCOUNT  33U


#define RELAY_PIN 3U
#define LIGHT_PIN 4U

#define HASH_SIZE 64U
#define KEY_SIZE 32U

#define TIMER_SEND_DATA_PERIOD 120000U // 2 minutes
#define TIMER_RECEIVE_COMMAND_PERIOD 100U // 3 seconds

// sensor data message type
#define SENSOR_DATA 40U
// configure message type
#define CONFIGURE 41U
// ACK message type
#define ACK 42U
#define CONTROL 43U
#define STATUS 44U
#define PUMP_ON 45U
#define PUMP_OFF 46U
#define LIGHT_ON 47U
#define LIGHT_OFF 48U

// public_key and secret_key (32 bytes) 
void GetPublicKey(uint8_t* public_key, uint8_t* secret_key);
void GetSecretKey(uint8_t* public_key, uint8_t* secret_key);


void to_hex_string(uint8_t* uintArray, size_t length, char* hexArray);
size_t hex_string_to_bytes(uint8_t *dest, size_t count, const char *src);

const float Round2Decimal(float x);
void LengthToBytes(byte* b, uint16_t length);
// Append src to dst, start at index
void Append(uint8_t* dst, size_t dst_length, size_t index, uint8_t* src, size_t src_length);


#endif