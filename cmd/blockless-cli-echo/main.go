package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {

	o := struct {
		Args []string `json:"args"`
		Env  []string `json:"env"`
	}{
		Args: os.Args,
		Env:  os.Environ(),
	}

	payload, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		log.Fatalf("could not JSON encode output: %s", err)
	}

	fmt.Printf("%s\n", payload)
}
