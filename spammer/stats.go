package spammer

import (
	"encoding/json"
	"fmt"
)

func getMessagePayload(args ...string) []byte {

	rec := Execute{
		Request: Request{
			FunctionID: testFunction.cid,
			Method:     testFunction.method,
			Arguments:  append([]string{}, args...),
			Config: RequestConfig{
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
