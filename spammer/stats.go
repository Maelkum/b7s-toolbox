package spammer

import (
	"encoding/json"
	"fmt"

	"github.com/Maelkum/b7s/models/execute"
	"github.com/Maelkum/b7s/models/request"
)

func getMessagePayload(args ...string) []byte {

	params := make([]execute.Parameter, 0, len(args))
	for _, arg := range args {
		params = append(params, execute.Parameter{Value: arg})
	}

	rec := request.Execute{
		Request: execute.Request{
			FunctionID: testFunction.cid,
			Method:     testFunction.method,
			Parameters: params,
			Config: execute.Config{
				NodeCount: 1,
			},
		},
	}

	data, err := json.Marshal(rec)
	if err != nil {
		panic("could not marshal message")
	}

	data = append(data, '\n')

	return data
}

func executionMapKey(i uint, total uint, frequency uint) string {
	return fmt.Sprintf("%v/%v-f%v", i+1, total, frequency)
}
