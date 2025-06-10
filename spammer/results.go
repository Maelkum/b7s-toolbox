package spammer

import (
	"io"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type testResult struct {
	config    testConfig
	responses []runResponse
}

type runResponse struct {
	success bool
	start   time.Time
	end     time.Time
}

type outputConfig struct {
	detailed   bool
	timeFormat string
}

type tableOption func(*outputConfig)

func detailedTable(cfg *outputConfig) {
	cfg.detailed = true
}

func processResults(stats *sync.Map, out io.WriteCloser, opts ...tableOption) {
	defer out.Close()

	cfg := outputConfig{
		detailed:   false,
		timeFormat: "15:04:05.000",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	var (
		total                            = 0
		ok                               = 0
		totalTime                        time.Duration
		totalTimeForSuccessfulExecutions time.Duration
	)

	t := table.NewWriter()
	t.SetOutputMirror(out)

	if cfg.detailed {
		// Table for detailed output.
		t.AppendHeader(table.Row{"Request", "Start", "End", "Duration", "Success"})
		t.AppendSeparator()
	}

	// Traverse the stats store.
	stats.Range(func(k, v any) bool {
		key := k.(string)
		val := v.(runResponse)

		duration := val.end.Sub(val.start)

		total++
		totalTime += duration
		if val.success {
			ok++
			totalTimeForSuccessfulExecutions += duration
		}

		log().Debug("processing stat",
			"key", key,
		)

		if cfg.detailed {
			t.AppendRow(table.Row{
				key, val.start.Format(cfg.timeFormat), val.end.Format(cfg.timeFormat), duration, val.success,
			})
		}

		return true
	})

	t.AppendSeparator()
	t.AppendSeparator()

	t.Render()

	// Rest table for the global summary
	t = table.NewWriter()
	t.SetOutputMirror(out)
	t.AppendHeader(
		table.Row{
			"Metric", "Count",
		},
	)

	var (
		tpe  string = "N/A" // time per execution
		tpse string = "N/A" // time per successful excution
	)
	if total > 0 {
		tpe = time.Duration(int64(totalTime) / int64(total)).String()
	}

	if ok > 0 {
		tpse = time.Duration(int64(totalTimeForSuccessfulExecutions) / int64(ok)).String()
	}

	t.AppendSeparator()
	t.AppendRow(table.Row{
		"total", total,
	})
	t.AppendRow(table.Row{
		"success", ok,
	})
	t.AppendRow(table.Row{
		"time", totalTime.String(),
	})
	t.AppendRow(table.Row{
		"time per execution", tpe,
	})
	t.AppendRow(table.Row{
		"time for successful executions", totalTimeForSuccessfulExecutions.String(),
	})
	t.AppendRow(table.Row{
		"time per successful executions", tpse,
	})

	t.Render()
}
