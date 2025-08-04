package spammer

type testConfig struct {
	executions uint
	frequency  uint
}

var testProfiles = []testConfig{
	// executions, frequency (requests per second)
	// {100, 10},
	{1000, 10},
	{10_000, 50},
	{10_000, 100},
	{10_000, 140},
	{10_000, 200},
	// {10_000, 300},
	// {10_000, 500},
	// {10_000, 1000},
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
