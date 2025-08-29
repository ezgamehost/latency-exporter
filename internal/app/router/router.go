package router

import (
	"fmt"

	configparser "github.com/ezgamehost/latency-exporter/internal/lib/configparser"
)

// MeasurementResult represents a generic measurement result
type MeasurementResult struct {
	Success   bool
	Latency   float64 // in milliseconds
	ErrorMsg  string
	ExtraData map[string]interface{} // for method-specific data
}

// Runner defines the interface for different measurement methods
type Runner interface {
	Run(endpoint string) MeasurementResult
	GetMethodName() string
}

// RouterConfig manages routing to different runners based on configuration
type RouterConfig struct {
	runners map[string]Runner
	config  *configparser.LatencyParserConfig
}

// NewRouterConfig creates a new router configuration
func NewRouterConfig(config *configparser.LatencyParserConfig) *RouterConfig {
	rc := &RouterConfig{
		runners: make(map[string]Runner),
		config:  config,
	}

	// Register available runners
	rc.registerRunners()
	return rc
}

// registerRunners registers all available measurement runners
func (rc *RouterConfig) registerRunners() {
	rc.runners["http"] = &HTTPRunner{}
	rc.runners["https"] = &HTTPRunner{}
	rc.runners["icmp"] = &ICMPRunner{}
	rc.runners["ping"] = &ICMPRunner{} // alias for icmp
}

// GetRunner returns the appropriate runner for a given method
func (rc *RouterConfig) GetRunner(method string) (Runner, error) {
	runner, exists := rc.runners[method]
	if !exists {
		return nil, fmt.Errorf("unsupported measurement method: %s", method)
	}
	return runner, nil
}

// GetAllDestinations returns all configured destinations
func (rc *RouterConfig) GetAllDestinations() []configparser.DestinationConfig {
	return rc.config.Destinations
}

// GetDestinationBySlug finds a destination by its metrics slug
func (rc *RouterConfig) GetDestinationBySlug(slug string) (*configparser.DestinationConfig, error) {
	for _, dest := range rc.config.Destinations {
		if dest.MetricsSlug == slug {
			return &dest, nil
		}
	}
	return nil, fmt.Errorf("destination with slug '%s' not found", slug)
}

// RunMeasurement executes a measurement for a specific destination
func (rc *RouterConfig) RunMeasurement(dest configparser.DestinationConfig) (MeasurementResult, error) {
	runner, err := rc.GetRunner(dest.Method)
	if err != nil {
		return MeasurementResult{}, err
	}

	result := runner.Run(dest.Endpoint)
	return result, nil
}
