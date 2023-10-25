package main

import (
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("specify key file")
	}

	path := os.Args[1]
	key, err := readPrivateKey(path)
	if err != nil {
		log.Fatalf("could not read key file (%s): %s", path, err)
	}

	id, err := peer.IDFromPublicKey(key.GetPublic())
	if err != nil {
		log.Fatalf("could not determine identity: %s", err)
	}

	fmt.Printf("%s\n", id)

}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}
