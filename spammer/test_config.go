package spammer

type testConfig struct {
	executions uint
	frequency  uint
}

// TODO: Update these.
var testExecutions = []uint{
	3,
	// 1_000,
	// 10_000,
	// 100_000,
}

var testFrequencies = []uint{
	10,
	// 100,
	// 1000,
}

var testFunction = struct {
	cid    string
	method string
}{
	cid:    "bafybeie3nlygbnuxhvqv3gvwa2hmd4tcfzk5jtvscwl6qs3ljn5tknlt4q",
	method: "echo.wasm",
}

func getTestProfiles() []testConfig {

	var cfgs []testConfig

	for _, ex := range testExecutions {
		for _, f := range testFrequencies {

			cfgs = append(cfgs,
				testConfig{
					executions: ex,
					frequency:  f,
				})
		}
	}

	return cfgs
}
