package spammer

import (
	"log/slog"
	"os"
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

func processResults(stats *sync.Map, opts ...tableOption) {

	cfg := outputConfig{
		detailed:   false,
		timeFormat: "15:04:05.000",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	total := 0
	ok := 0

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if cfg.detailed {
		// Table for detailed output.
		t.AppendHeader(table.Row{"Request", "Start", "End", "Duration", "Success"})
		t.AppendSeparator()
	}

	// Traverse the stats store.
	stats.Range(func(k, v any) bool {
		key := k.(string)
		val := v.(runResponse)

		total++
		if val.success {
			ok++
		}

		slog.Debug("processing result",
			"key", key,
		)

		if cfg.detailed {
			t.AppendRow(table.Row{
				key, val.start.Format(cfg.timeFormat), val.end.Format(cfg.timeFormat), val.end.Sub(val.start), val.success,
			})
		}

		return true
	})

	t.AppendSeparator()
	t.AppendSeparator()

	t.Render()

	// Rest table for the global summary
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		table.Row{
			"Metric", "Count",
		},
	)

	t.AppendSeparator()
	t.AppendRow(table.Row{"total", total})
	t.AppendRow(table.Row{"success", ok})

	t.Render()
}
