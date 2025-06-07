package spammer

import "time"

type testResult struct {
	config    testConfig
	responses []runResponse
}

type runResponse struct {
	success bool
	start   time.Time
	end     time.Time
}
