package router

import (
	"github.com/ezgamehost/latency-exporter/internal/app/exporter"
)

// HTTPRunner implements the Runner interface for HTTP/HTTPS measurements
type HTTPRunner struct{}

// Run executes an HTTP measurement
func (h *HTTPRunner) Run(endpoint string) MeasurementResult {
	measurement := exporter.HttpMeasure(endpoint)

	extraData := map[string]interface{}{
		"status_code": measurement.StatusCode,
	}

	return MeasurementResult{
		Success:   measurement.Success,
		Latency:   float64(measurement.SpeedMS),
		ErrorMsg:  "",
		ExtraData: extraData,
	}
}

// GetMethodName returns the method name for this runner
func (h *HTTPRunner) GetMethodName() string {
	return "http"
}
