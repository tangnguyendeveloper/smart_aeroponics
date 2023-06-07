#include "utilities.hpp"
#include <EEPROM.h>

float Round2Decimal(float value) {
   return (int)(value * 100 + 0.5) / 100.0;
}

void LengthToBytes(byte* b, uint16_t length) {
   *b = length >> 8;
   *(b + 1) = (uint8_t) (length & 0xFF);
}



void to_hex_string(uint8_t* uintArray, size_t length, char* hexArray) {
  for (size_t i = 0; i < length; i++) {
    sprintf(hexArray + (i * 2), "%02X", uintArray[i]);
  }
  hexArray[length * 2] = '\0';
}

size_t hex_string_to_bytes(uint8_t *dest, size_t count, const char *src) {
    size_t i = 0;
    int value = 0;
    for (i = 0; i < count && sscanf(src + i * 2, "%2x", &value) == 1; i++) {
        dest[i] = value;
    }
    return i;
}

void GetPublicKey(uint8_t* public_key, uint8_t* secret_key) {
    unsigned long seed = 0;
    uint8_t index = 0;
    seed += (analogRead(5) | ((unsigned long)analogRead(5) << 16));

    randomSeed(seed);

    while (index < KEY_SIZE) {
        *(public_key + index) = (uint8_t)random(256);
        *(secret_key + index) = EEPROM.read(*(public_key + index));
        index++;
    }

}

void GetSecretKey(uint8_t* public_key, uint8_t* secret_key){
    uint8_t index = 0;
    while (index < KEY_SIZE) {
        *(secret_key + index) = EEPROM.read(*(public_key + index));
        index++;
    }
}

void Append(uint8_t* dst, size_t dst_length, size_t index, uint8_t* src, size_t src_length){
    size_t n = 0;
    if (index + src_length > dst_length) n = dst_length - index;
    else n = src_length;

    for (size_t i = 0; i < n; i++) *(dst + index + i) = *(src + i);
}

