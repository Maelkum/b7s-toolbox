package main

import (
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

type DummyExecutor struct{}

func (d DummyExecutor) ExecuteFunction(id string, req execute.Request) (execute.Result, error) {

	res := execute.Result{
		Code:      codes.OK,
		RequestID: id,
		Usage:     execute.Usage{},
		Result: execute.RuntimeOutput{
			Stdout:   "dummy-output",
			Stderr:   "dummy-stderr-output",
			ExitCode: 0,
		},
	}

	return res, nil
}
