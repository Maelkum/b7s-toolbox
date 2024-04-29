package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/store"
)

func main() {

	const (
		dbpath = "db"
	)

	db, err := pebble.Open(dbpath, nil)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	store := store.New(db, JSONSerializer{})

	ids := os.Args[1:]
	if len(ids) == 0 {
		peers, err := store.RetrievePeers()
		if err != nil {
			log.Fatalf("could not retrieve peers: %s", err)
		}

		for _, peer := range peers {
			data, _ := json.Marshal(peer)
			fmt.Printf("%s\n", data)
		}

	}

	for _, id := range ids {

		peerID, err := peer.Decode(id)
		if err != nil {
			log.Fatalf("could not parse peer ID: %s", err)
		}

		p, err := store.RetrievePeer(peerID)
		if err != nil {
			log.Fatalf("could not retrieve peer: %s", err)
		}

		data, _ := json.Marshal(p)
		fmt.Printf("%s\n", data)
	}
}
