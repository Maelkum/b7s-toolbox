package spammer

import (
	"encoding/json"
	"errors"

	"github.com/blessnetwork/b7s/models/bls"
	"github.com/hashicorp/go-multierror"
)

// TODO: Delete all of this shit - it's just a workaround because of the `arguments` change in b7s request format.

var _ (json.Marshaler) = (*Execute)(nil)

// Execute describes the `MessageExecute` request payload.
type Execute struct {
	bls.BaseMessage

	Request // execute request is embedded.

	Topic string `json:"topic,omitempty"`
}

func (Execute) Type() string { return bls.MessageExecute }

func (e Execute) MarshalJSON() ([]byte, error) {
	type Alias Execute
	rec := struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(e),
		Type:  e.Type(),
	}
	return json.Marshal(rec)
}

// Request describes an execution request.
type Request struct {
	FunctionID string        `json:"function_id"`
	Method     string        `json:"method"`
	Arguments  []string      `json:"arguments,omitempty"`
	Config     RequestConfig `json:"config"`

	// Optional signature of the request.
	Signature string `json:"signature,omitempty"`
}

func (r Request) Valid() error {

	var err *multierror.Error

	if r.FunctionID == "" {
		err = multierror.Append(err, errors.New("function ID is required"))
	}

	if r.Method == "" {
		err = multierror.Append(err, errors.New("method is required"))
	}

	return err.ErrorOrNil()
}

// Parameter represents an execution parameter, modeled as a key-value pair.
type Parameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Config represents the configurable options for an execution request.
type RequestConfig struct {
	Environment []EnvVar `json:"env_vars,omitempty"`
	Stdin       *string  `json:"stdin,omitempty"`

	// NodeCount specifies how many nodes should execute this request.
	NodeCount int `json:"number_of_nodes,omitempty"`

	// When should the execution timeout
	Timeout int `json:"timeout,omitempty"`

	// Threshold (percentage) defines how many nodes should respond with a result to consider this execution successful.
	Threshold float64 `json:"threshold,omitempty"`
}

// EnvVar represents the name and value of the environment variables set for the execution.
type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ResultAggregation struct {
	Enable     bool        `json:"enable,omitempty"`
	Type       string      `json:"type,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty"`
}
