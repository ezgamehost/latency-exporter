package metrics

import (
	"fmt"
	"strings"

	"github.com/ezgamehost/latency-exporter/internal/app/router"
	configparser "github.com/ezgamehost/latency-exporter/internal/lib/configparser"
)

// PrometheusFormatter formats measurement results as Prometheus metrics
type PrometheusFormatter struct{}

// FormatMetrics converts measurement results to Prometheus format
func (pf *PrometheusFormatter) FormatMetrics(dest configparser.DestinationConfig, result router.MeasurementResult) string {
	var output strings.Builder

	// Create consistent labels for all metrics
	labels := fmt.Sprintf(`target="%s",method="%s",endpoint="%s",name="%s"`,
		dest.MetricsSlug, dest.Method, dest.Endpoint, dest.Name)

	// Latency metric
	output.WriteString("latency_measurement_seconds{")
	output.WriteString(labels)
	output.WriteString(fmt.Sprintf("} %.6f\n", result.Latency/1000.0)) // Convert ms to seconds

	// Success metric
	successValue := 0
	if result.Success {
		successValue = 1
	}
	output.WriteString("latency_measurement_success{")
	output.WriteString(labels)
	output.WriteString(fmt.Sprintf("} %d\n", successValue))

	// Method-specific metrics
	switch dest.Method {
	case "http", "https":
		if statusCode, ok := result.ExtraData["status_code"].(int); ok {
			output.WriteString("latency_measurement_http_status_code{")
			output.WriteString(labels)
			output.WriteString(fmt.Sprintf("} %d\n", statusCode))
		}
	case "icmp", "ping":
		// ICMP-specific metrics could be added here if needed
	}

	return output.String()
}

// FormatAllMetrics formats metrics for all destinations
func (pf *PrometheusFormatter) FormatAllMetrics(destinations []configparser.DestinationConfig, results map[string]router.MeasurementResult) string {
	var output strings.Builder

	// Write metric definitions once at the top
	output.WriteString("# HELP latency_measurement_seconds Latency measurement in seconds\n")
	output.WriteString("# TYPE latency_measurement_seconds gauge\n")
	output.WriteString("# HELP latency_measurement_success Success indicator (1=success, 0=failure)\n")
	output.WriteString("# TYPE latency_measurement_success gauge\n")
	output.WriteString("# HELP latency_measurement_http_status_code HTTP status code for HTTP measurements\n")
	output.WriteString("# TYPE latency_measurement_http_status_code gauge\n")

	for _, dest := range destinations {
		if result, exists := results[dest.MetricsSlug]; exists {
			output.WriteString(pf.FormatMetrics(dest, result))
		}
	}

	return strings.TrimSpace(output.String())
}

// sanitizeMetricName ensures metric names are valid for Prometheus
func sanitizeMetricName(name string) string {
	// Replace invalid characters with underscores
	result := strings.ReplaceAll(name, ".", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ToLower(result)

	// Ensure it starts with a letter or underscore
	if len(result) > 0 && (result[0] >= '0' && result[0] <= '9') {
		result = "_" + result
	}

	return result
}
