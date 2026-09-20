package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/bytemare/opaque"
)

func main() {
	conf := opaque.DefaultConfiguration()

	privateKey, _ := conf.KeyGen()

	oprfSeed := make([]byte, 64)
	if _, err := rand.Read(oprfSeed); err != nil {
		log.Fatalf("generate OPRF seed: %v", err)
	}

	fmt.Println("# Generated OPAQUE server key material. Put in the .env file. Keep these secret.")
	fmt.Printf("OPAQUE_SERVER_PRIVATE_KEY=%s\n", hex.EncodeToString(privateKey.Encode()))
	fmt.Printf("OPAQUE_OPRF_SEED_SECRET=%s\n", hex.EncodeToString(oprfSeed))
}
