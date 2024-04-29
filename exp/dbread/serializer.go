package main

import "encoding/json"

type JSONSerializer struct{}

func (s JSONSerializer) Marshal(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func (s JSONSerializer) Unmarshal(data []byte, target any) error {
	return json.Unmarshal(data, target)
}
