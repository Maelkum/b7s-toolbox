package spammer

type testConfig struct {
	executions uint
	frequency  uint
}

var testProfiles = []testConfig{
	// executions, frequency (requests per second)
	{1000, 10},
	{1000, 100},
	{1000, 200},
	// {2_000, 1000},
}

var testFunction = struct {
	cid    string
	method string
}{
	cid:    "bafybeie3nlygbnuxhvqv3gvwa2hmd4tcfzk5jtvscwl6qs3ljn5tknlt4q",
	method: "echo.wasm",
}

func getTestProfiles() []testConfig {

	return testProfiles
}
