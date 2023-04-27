package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	sleepEnvVar = "B7S_SLEEP"
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

	sec, _ := strconv.ParseInt(os.Getenv(sleepEnvVar), 10, 32)

	if sec > 0 {
		time.Sleep(time.Duration(sec) * time.Second)
	}
}
