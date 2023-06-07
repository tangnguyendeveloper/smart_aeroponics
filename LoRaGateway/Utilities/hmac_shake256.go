package Utilities

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"

	"golang.org/x/crypto/sha3"
)

func GetPublicKey(id string) (pk []byte, sk []byte) {

	ev := getEV(id)
	if ev == nil {
		log.Printf("ERROR: Unknown device ID %s\n", id)
		return nil, nil
	}

	pk = make([]byte, KEY_SIZE)
	sk = make([]byte, KEY_SIZE)

	rand.Read(pk)

	for idx, b := range pk {
		sk[idx] = ev[b]
	}

	return pk, sk
}

func GetSecretKey(public_key []byte, id string) []byte {
	ev := getEV(id)
	if ev == nil {
		return nil
	}

	sk := make([]byte, KEY_SIZE)

	for idx, b := range public_key {
		sk[idx] = ev[b]
	}

	return sk
}

func HMAC_SHAKE256(message, secret_key []byte) string {
	OutputSize := 64

	hash := sha3.NewShake256()
	// hash.Write(secret_key)
	// hash_key := make([]byte, OutputSize)
	// hash.Read(hash_key)

	// hash.Reset()
	hash.Write(append(message, secret_key...))
	hash_key_msg := make([]byte, OutputSize)
	hash.Read(hash_key_msg)

	hash.Reset()
	hash.Write(append(hash_key_msg, secret_key...))
	hash.Read(hash_key_msg)

	for i := range hash_key_msg[:KEY_SIZE] {
		hash_key_msg[i] ^= hash_key_msg[i+KEY_SIZE]
	}

	hash.Clone()

	return strings.ToUpper(hex.EncodeToString(hash_key_msg[:KEY_SIZE]))

}
