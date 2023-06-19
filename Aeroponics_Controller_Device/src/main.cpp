#include "utilities.hpp"
#include <SHAKE.h>

SHAKE256 shake256;

uint8_t public_key[KEY_SIZE];
uint8_t secret_key[KEY_SIZE];
// output (64 bytes)
void SHAKE_256(const uint8_t* input, const size_t& input_length, uint8_t* output);
// signature (32 bytes) key (32 bytes)
void HMAC_SHAKE_256(uint8_t* message, size_t length_message, uint8_t* key, char* signature);


DHT dht(DHTPIN, DHTTYPE);
float Temperature, Humidity = 0, HeatIndex = 0;
bool ReadTemperatureAndHumidity();

uint16_t analogBuffer[SCOUNT];
float TDSValue = 1, ECvalue = 0;
const uint16_t getMedianNum(uint16_t* bArray, uint16_t iFilterLen);
void ReadTDSAndECLevel();

// 6 TX, 7 RX on LoRa E32
SoftwareSerial LoRa(6, 7);
DynamicJsonDocument data(256);

const char* deviceID = "2598354207890938361";
uint8_t send_payload[2*HASH_SIZE];
char hex_signature[2*KEY_SIZE+1];
char hex_public_key[2*KEY_SIZE+1];
void SendSensorData();

void ReceiveCommand();
void SendACK(const uint8_t& action, const uint8_t& status);

unsigned long last_send_time = 0, last_receive_time = 0, last_control = 0;

unsigned long pump_on_time = 15UL*1000*60;
unsigned long pump_off_time = 45UL*1000*60;
uint8_t pump_status = 0, light_status = 0;

void setup() {
    Serial.begin(9600);
    Serial.println(F("DHT and TDS test!"));

    pinMode(TdsSensorPin,INPUT);
    dht.begin();
    LoRa.begin(9600);

    pinMode(RELAY_PIN, OUTPUT);
    pinMode(LIGHT_PIN, OUTPUT);

    last_control = millis();
    digitalWrite(RELAY_PIN, HIGH);
    pump_status = 1;
}

void loop() {
   // Wait a few seconds between measurements.

    if (millis() - last_send_time > TIMER_SEND_DATA_PERIOD) {
        SendSensorData();
        last_send_time = millis();
    }

    if (millis() - last_receive_time > TIMER_RECEIVE_COMMAND_PERIOD) {
        ReceiveCommand();
        last_receive_time = millis();
    }

    if (pump_status == 1 && millis() - last_control > pump_on_time) {
        pump_status = 0;
        digitalWrite(RELAY_PIN, LOW);
        last_control = millis();
    }

    if (pump_status == 0 && millis() - last_control > pump_off_time) {
        pump_status = 1;
        digitalWrite(RELAY_PIN, HIGH);
        last_control = millis();
    }

    delay(10);
}

bool ReadTemperatureAndHumidity() {
    Humidity = dht.readHumidity();
    Temperature = dht.readTemperature();

    if (isnan(Humidity) || isnan(Temperature)) {
      Serial.println(F("Failed to read from DHT sensor!"));
      return false;
    }

    HeatIndex = dht.computeHeatIndex(Temperature, Humidity, false);

    return true;
}


const uint16_t getMedianNum(uint16_t* bArray, uint16_t iFilterLen) {

      uint16_t i, j, bTemp;

      for (j = 0; j < iFilterLen - 1; j++) {
        for (i = 0; i < iFilterLen - j - 1; i++) {
          if (*(bArray + i) > *(bArray + i + 1)) {
            bTemp = *(bArray + i);
            *(bArray + i) = *(bArray + i + 1);
            *(bArray + i + 1) = bTemp;
         }
        }
      }
      
    return *(bArray + 16);
}

void ReadTDSAndECLevel() {
    uint8_t analogBufferIndex = 0;

    while (analogBufferIndex < SCOUNT)
    {
        analogBuffer[analogBufferIndex] = (uint16_t)analogRead(TdsSensorPin);    //read the analog value and store into the buffer
        analogBufferIndex++;
        if(analogBufferIndex == SCOUNT) {
          analogBufferIndex = 0U;
          break;
        }
        delay(40U);
    }   
    
    const float averageVoltage = getMedianNum(analogBuffer,SCOUNT) * VREF / 1024.0f; // read the analog value more stable by the median filtering algorithm, and convert to voltage value
    const float compensationCoefficient=1.0f+0.02f*(HeatIndex-25.0f);    //temperature compensation formula: fFinalResult(25^C) = fFinalResult(current)/(1.0+0.02*(fTP-25.0));
    const float compensationVolatge=averageVoltage/compensationCoefficient;  //temperature compensation
    TDSValue=(133.42f*compensationVolatge*compensationVolatge*compensationVolatge - 255.86f*compensationVolatge*compensationVolatge + 857.39f*compensationVolatge)*0.5f; //convert voltage value to tds value

    ECvalue = (TDSValue * 2U) / 1000U;
}


void SHAKE_256(const uint8_t* input, const size_t& input_length, uint8_t* output) {

    size_t posn = 0, len = 0;

    shake256.reset();

    for (posn = 0; posn < input_length; posn += KEY_SIZE) {
         len = input_length - posn;
         if (len > KEY_SIZE)
             len = KEY_SIZE;
         shake256.update(input + posn, len);
    }

    for (posn = 0; posn < HASH_SIZE; posn += KEY_SIZE) {
        len = HASH_SIZE - posn;
        if (len > KEY_SIZE)
            len = KEY_SIZE;
        shake256.extend(output + posn, len);
    }
}

void HMAC_SHAKE_256(uint8_t* message, size_t length_message, uint8_t* key, char* signature) {
    
    // uint8_t hash_key[HASH_SIZE];
    // SHAKE_256(key, 32, hash_key);

    uint8_t* input = new uint8_t[length_message+KEY_SIZE];
    Append(input, length_message+KEY_SIZE, 0, message, length_message);
    Append(input, length_message+KEY_SIZE, length_message, key, KEY_SIZE);

    uint8_t hash_key_msg[HASH_SIZE];
    SHAKE_256(input, length_message+KEY_SIZE, hash_key_msg);
    delete[] input;

    input = new uint8_t[HASH_SIZE+KEY_SIZE];
    Append(input, HASH_SIZE+KEY_SIZE, 0, hash_key_msg, HASH_SIZE);
    Append(input, HASH_SIZE+KEY_SIZE, HASH_SIZE, key, KEY_SIZE);

    SHAKE_256(input, HASH_SIZE+KEY_SIZE, hash_key_msg);
    delete[] input;

    for (uint8_t i = 0; i < KEY_SIZE; i++) hash_key_msg[i] ^= hash_key_msg[i+KEY_SIZE];

    to_hex_string(hash_key_msg, KEY_SIZE, signature);

}


void SendSensorData() {

    bool ok = ReadTemperatureAndHumidity();
    if (ok) {
      ReadTDSAndECLevel();
    } else return;

    data["ID"] = deviceID;
    data["TP"] = Round2Decimal(Temperature); 
    data["HU"] = Round2Decimal(Humidity);
    data["HI"] = Round2Decimal(HeatIndex);
    data["TD"] = Round2Decimal(TDSValue);
    data["EC"] = Round2Decimal(ECvalue);

    size_t n = serializeJson(data, send_payload);
    GetPublicKey(public_key, secret_key);
    HMAC_SHAKE_256(send_payload, n, secret_key, hex_signature);
    data["SG"] = hex_signature;

    to_hex_string(public_key, KEY_SIZE, hex_public_key);
    data["PK"] = hex_public_key;
    
    //n = serializeJson(data, Serial);
    //Serial.print("\n Length: ");
    //Serial.println(n);

    n = measureJson(data);

    byte length[2];
    LengthToBytes(length, n+1);
    LoRa.write(length, 2);
    LoRa.write(SENSOR_DATA);
    n = serializeJson(data, LoRa);
    // Serial.print("Send in: ");
    // Serial.print(n);
    // Serial.print(" bytes \n");

    data.clear();
    //Serial.println("Send Sensor Data OK");

}

void ReceiveCommand() {

    //Serial.print("Receive Command ...");
    String command_msg = LoRa.readStringUntil('}');
    command_msg.trim();
    //Serial.println("Done!");
    if (command_msg[0] == STATUS) {
        if (command_msg[2] == 'p' && pump_status == 1) SendACK(STATUS, PUMP_ON);
        else if (command_msg[2] == 'p') SendACK(STATUS, PUMP_OFF);
        else if (command_msg[2] == 'l' && light_status == 1) SendACK(STATUS, LIGHT_ON);
        else if (command_msg[2] == 'l') SendACK(STATUS, LIGHT_OFF);
        return;
    }

    if (command_msg.length() < 150) return;
    if (command_msg.substring(129, 148) != deviceID) return;
    if (*command_msg.end() != '}') command_msg += '}';
    
    Serial.println(command_msg);

    // hex_string_to_bytes(public_key, KEY_SIZE, command_msg.begin()+65);
    // GetSecretKey(public_key, secret_key);

    // HMAC_SHAKE_256((uint8_t*)command_msg.begin()+65, command_msg.length()-65, secret_key, hex_signature);
    
    // if (memcmp(hex_signature, command_msg.begin() + 1, HASH_SIZE) != 0) {
    //     Serial.println("Verify signature failed!");
    //     return;
    // } 
    
    deserializeJson(data, command_msg.begin()+148);

    switch (command_msg[0]){
        case CONTROL:

            if (data.containsKey("pump")) {
                if (data["pump"].as<uint8_t>() == 0) {
                digitalWrite(RELAY_PIN, LOW);
                SendACK(CONTROL, PUMP_OFF);
                }
                else if (data["pump"].as<uint8_t>() == 1) {
                digitalWrite(RELAY_PIN, HIGH);
                SendACK(CONTROL, PUMP_ON);
                }
                pump_status = data["pump"].as<uint8_t>();
                last_control = millis();
            }

            if (data.containsKey("light")) {
                if (data["light"].as<uint8_t>() == 0) {
                digitalWrite(LIGHT_PIN, LOW);
                SendACK(CONTROL, LIGHT_OFF);
                }
                else if (data["light"].as<uint8_t>() == 1) {
                digitalWrite(LIGHT_PIN, HIGH);
                SendACK(CONTROL, LIGHT_ON);
                }
                light_status = data["light"].as<uint8_t>();
                last_control = millis();
            }
            
            break;
        
        case CONFIGURE:

            if (data.containsKey("ON")) {
                pump_on_time = data["ON"].as<unsigned long>()*60*1000;
            }
            if (data.containsKey("OFF")) {
                pump_off_time = data["OFF"].as<unsigned long>()*60*1000;
            }
            SendACK(CONFIGURE, CONFIGURE);
            
            break;

        default:
            break;
    }

    data.clear();

}


void SendACK(const uint8_t& action,  const uint8_t& status) {
    byte length[2];
    LengthToBytes(length, 22);

    LoRa.write(length, 2);
    LoRa.write(ACK);
    LoRa.write(deviceID);
    LoRa.write(action);
    LoRa.write(status);

}