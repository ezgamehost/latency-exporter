package router

import (
	"github.com/ezgamehost/latency-exporter/internal/lib/fping"
)

// ICMPRunner implements the Runner interface for ICMP ping measurements
type ICMPRunner struct{}

// Run executes an ICMP ping measurement
func (i *ICMPRunner) Run(endpoint string) MeasurementResult {
	pingResult, err := fping.ExecPing(endpoint)

	if err != nil {
		return MeasurementResult{
			Success:   false,
			Latency:   0,
			ErrorMsg:  err.Error(),
			ExtraData: make(map[string]interface{}),
		}
	}

	extraData := map[string]interface{}{
		"host": pingResult.Host,
	}

	return MeasurementResult{
		Success:   pingResult.Alive,
		Latency:   pingResult.Latency,
		ErrorMsg:  "",
		ExtraData: extraData,
	}
}

// GetMethodName returns the method name for this runner
func (i *ICMPRunner) GetMethodName() string {
	return "icmp"
}
