package Utilities

import (
	"log"
	"strconv"
)

const (
	KEY_SIZE = 32

	// sensor data message type
	SENSOR_DATA = 40
	// configure message type
	CONFIGURE = 41
	// ACK message type
	ACK = 42

	CONTROL = 43
	STATUS  = 44

	PUMP_ON   = 45
	PUMP_OFF  = 46
	LIGHT_ON  = 47
	LIGHT_OFF = 48
)

func getEV(id string) []byte {

	n_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Printf("ERROR: Parse Uint64 from %s, %s\n ", id, err)
		return nil
	}

	switch n_id {
	case 2598354207890938361:
		return []byte{0x57, 0x69, 0x48, 0x72, 0x48, 0x5D, 0x31, 0x2D, 0x44, 0x4C, 0x25, 0x41, 0x77, 0x59, 0x41, 0x3C, 0x71, 0x49, 0x37, 0x6B, 0x30, 0x64, 0x77, 0x70, 0x33, 0x2C, 0x75, 0x2A, 0x3A, 0x5E, 0x35, 0x5D, 0x2D, 0x31, 0x55, 0x40, 0x4B, 0x3B, 0x3A, 0x51, 0x31, 0x44, 0x64, 0x4A, 0x62, 0x42, 0x2D, 0x53, 0x34, 0x56, 0x60, 0x4C, 0x62, 0x65, 0x47, 0x7C, 0x25, 0x2B, 0x78, 0x6C, 0x44, 0x6A, 0x2A, 0x46, 0x41, 0x2D, 0x21, 0x25, 0x5D, 0x69, 0x5C, 0x2B, 0x49, 0x5E, 0x74, 0x45, 0x69, 0x56, 0x48, 0x4D, 0x2F, 0x6F, 0x3B, 0x40, 0x53, 0x33, 0x41, 0x77, 0x78, 0x3E, 0x7B, 0x7A, 0x34, 0x20, 0x6D, 0x7C, 0x48, 0x54, 0x6C, 0x50, 0x4D, 0x79, 0x45, 0x69, 0x4A, 0x5C, 0x2E, 0x4B, 0x7D, 0x20, 0x47, 0x69, 0x44, 0x49, 0x78, 0x55, 0x2A, 0x3D, 0x2B, 0x63, 0x50, 0x5C, 0x4C, 0x3A, 0x50, 0x38, 0x58, 0x78, 0x65, 0x59, 0x45, 0x26, 0x41, 0x49, 0x5E, 0x78, 0x26, 0x3F, 0x53, 0x25, 0x5F, 0x54, 0x6E, 0x5D, 0x55, 0x25, 0x34, 0x55, 0x75, 0x65, 0x48, 0x7B, 0x3E, 0x35, 0x2B, 0x50, 0x78, 0x7C, 0x54, 0x3D, 0x69, 0x40, 0x39, 0x23, 0x5C, 0x40, 0x2B, 0x78, 0x31, 0x3A, 0x6D, 0x47, 0x55, 0x3D, 0x76, 0x2A, 0x38, 0x78, 0x74, 0x3D, 0x75, 0x43, 0x46, 0x5C, 0x54, 0x2C, 0x27, 0x40, 0x29, 0x21, 0x63, 0x61, 0x2E, 0x3C, 0x39, 0x40, 0x47, 0x54, 0x72, 0x26, 0x58, 0x3C, 0x4F, 0x50, 0x2B, 0x60, 0x4E, 0x54, 0x5D, 0x5C, 0x2F, 0x6A, 0x60, 0x43, 0x33, 0x37, 0x45, 0x49, 0x5C, 0x5B, 0x67, 0x4D, 0x58, 0x76, 0x58, 0x7D, 0x4C, 0x20, 0x3C, 0x4B, 0x51, 0x7D, 0x3A, 0x55, 0x37, 0x2C, 0x45, 0x75, 0x6A, 0x4B, 0x57, 0x66, 0x2E, 0x3E, 0x64, 0x37, 0x73, 0x47, 0x21, 0x64, 0x40, 0x70, 0x6B, 0x47, 0x35, 0x70}

	default:
		return nil
	}

}
