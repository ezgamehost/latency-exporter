package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezgamehost/latency-exporter/internal/app/metrics"
	"github.com/ezgamehost/latency-exporter/internal/app/router"
	configparser "github.com/ezgamehost/latency-exporter/internal/lib/configparser"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var config *configparser.LatencyParserConfig
var routerConfig *router.RouterConfig
var prometheusFormatter *metrics.PrometheusFormatter

func init() {
	var err error
	config, err = configparser.ConfigParser()
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config: %v", err))
	}

	routerConfig = router.NewRouterConfig(config)
	prometheusFormatter = &metrics.PrometheusFormatter{}
}

func main() {
	fmt.Printf("Loaded config with %d destinations\n", len(config.Destinations))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Route for individual resource metrics
	r.Get("/metrics/{resource}", handleResourceMetrics)

	// Route for all metrics
	r.Get("/metrics", handleAllMetrics)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func handleResourceMetrics(w http.ResponseWriter, r *http.Request) {
	resource := chi.URLParam(r, "resource")

	// Find the destination by metrics slug
	dest, err := routerConfig.GetDestinationBySlug(resource)
	if err != nil {
		http.Error(w, fmt.Sprintf("Resource not found: %s", resource), http.StatusNotFound)
		return
	}

	// Run the measurement
	result, err := routerConfig.RunMeasurement(*dest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Measurement failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Format as Prometheus metrics
	metrics := prometheusFormatter.FormatMetrics(*dest, result)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

func handleAllMetrics(w http.ResponseWriter, r *http.Request) {
	destinations := routerConfig.GetAllDestinations()
	results := make(map[string]router.MeasurementResult)

	// Run measurements for all destinations
	for _, dest := range destinations {
		result, err := routerConfig.RunMeasurement(dest)
		if err != nil {
			log.Printf("Failed to measure %s: %v", dest.Name, err)
			// Continue with other measurements even if one fails
			continue
		}
		results[dest.MetricsSlug] = result
	}

	// Format all metrics
	allMetrics := prometheusFormatter.FormatAllMetrics(destinations, results)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(allMetrics))
}
