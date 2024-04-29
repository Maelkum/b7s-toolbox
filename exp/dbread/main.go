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

const usage = `
dbread - read everything in the database
dbread <peer|function> [id1 id2 id3]`

func main() {

	const (
		dbpath = "db"
	)

	db, err := pebble.Open(dbpath, nil)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	store := store.New(db, JSONSerializer{})

	args := os.Args[1:]
	if len(args) == 0 {
		readAll(store)
		return
	}

	if len(args) < 2 {
		fmt.Println(usage)
		return
	}

	typ := args[0]

	switch typ {
	case "peer":
		readPeers(store, args[1:])
	case "function":
		readFunctions(store, args[1:])
	default:
		fmt.Println(usage)
		return
	}
}

func readPeers(store *store.Store, ids []string) {

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

func readFunctions(store *store.Store, ids []string) {

	for _, id := range ids {
		p, err := store.RetrieveFunction(id)
		if err != nil {
			log.Fatalf("could not retrieve peer: %s", err)
		}

		data, _ := json.Marshal(p)
		fmt.Printf("%s\n", data)
	}
}

func readAll(store *store.Store) {

	peers, err := store.RetrievePeers()
	if err != nil {
		log.Fatalf("could not retrieve peers: %s", err)
	}
	for _, rec := range peers {
		data, _ := json.Marshal(rec)
		fmt.Printf("%s\n", data)
	}

	functions, err := store.RetrieveFunctions()
	if err != nil {
		log.Fatalf("could not retrieve functions: %s", err)
	}
	for _, rec := range functions {
		data, _ := json.Marshal(rec)
		fmt.Printf("%s\n", data)
	}
}
