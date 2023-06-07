// package main

// import (
// 	"LoRaGateway/Utilities"
// 	"crypto/hmac"
// 	"encoding/hex"
// 	"fmt"
// 	"log"
// 	"strings"

// 	"golang.org/x/crypto/sha3"
// )

// func main() {
// 	// Key and message
// 	key := GetSecretKey("a7f1d92a82c8d8fe434d98558ce2b347171198542f112d0558f56bd688079992")

// 	message := []byte(`{"TP":31.1,"HU":81.2,"HI":41.65,"TD":0,"EC":0}`)
// 	// Create a new HMAC hasher with SHA3-256
// 	hash := hmac.New(sha3.New256, key)

// 	// Write the message to the hasher
// 	hash.Write(message)

// 	// Get the HMAC hash
// 	hmacHash := hash.Sum(nil)

// 	// Convert the hash to a hex string
// 	hmacString := hex.EncodeToString(hmacHash)

// 	fmt.Println("key: ", strings.ToUpper(hex.EncodeToString(key)))
// 	fmt.Println("HMAC (SHA3-256):", hmacString)
// 	fmt.Println("Block size:", hash.BlockSize())
// }

// func GetSecretKey(hex_string string) []byte {
// 	key, err := hex.DecodeString(hex_string)
// 	if err != nil {
// 		log.Printf("ERROR: Get secret key from public key (%x), %s\n", hex_string, err)
// 	}

// 	for i := 0; i < Utilities.KEY_SIZE; i++ {
// 		key[i] = Utilities.EntropyVector[key[i]]
// 	}

// 	return key
// }

package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

func main() {
	message := []byte("abcdefghiklmnojhgbnkiuybnjkiuygjkiuygbjkiuygbnkiuytfhjkiuytgjkiuytghnjkiuyhnkiuyhnkiuyhnjmkiuyhgnkiuygbjuytgbhytg")
	key := []byte("123456789012345678901234567890ab")
	OutputSize := 64

	hash := sha3.NewShake256()
	hash.Write(key)
	hash_key := make([]byte, OutputSize)
	hash.Read(hash_key)

	hash.Reset()
	hash.Write([]byte(string(message) + string(key)))
	hash_key_msg := make([]byte, OutputSize)
	hash.Read(hash_key_msg)

	hash.Reset()
	hash.Write([]byte(string(hash_key_msg) + string(hash_key)))
	hash.Read(hash_key_msg)

	for i := range hash_key_msg[:32] {
		hash_key_msg[i] += hash_key_msg[i+32]
	}

	fmt.Println("HMAC: ", hex.EncodeToString(hash_key_msg[:32]))

}
