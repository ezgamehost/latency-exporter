package exporter

import (
	"net/http"
	"time"
)

type HttpMeasurement struct {
	Success    bool
	SpeedMS    int64
	StatusCode int
}

func HttpMeasure(endpoint string) HttpMeasurement {
	start := time.Now()
	resp, err := http.Get(endpoint)
	end := time.Now()
	if err != nil {
		return HttpMeasurement{Success: false, SpeedMS: 0, StatusCode: 500}
	}
	defer resp.Body.Close()

	speed := end.Sub(start).Milliseconds()
	statuscode := resp.StatusCode

	return HttpMeasurement{SpeedMS: speed, StatusCode: statuscode, Success: true}
}
