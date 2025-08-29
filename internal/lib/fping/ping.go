package fping

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PingResult struct {
	Host    string
	Alive   bool
	Latency float64 // ms
}

func ExecPing(address string) (*PingResult, error) {
	// Run fping
	cmd := exec.Command("fping", "-c1", address)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		// fping returns nonzero if host unreachable, but we still parse output
		return nil, fmt.Errorf("command error: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected output: %s", string(out))
	}

	line1, line2 := lines[0], lines[1]

	// Regex patterns
	line1Re := regexp.MustCompile(`^([\d\.]+)\s+:.*?,\s+([\d\.]+)\s+ms.*\(([\d\.]+)\s+avg,\s+(\d+)% loss\)`)
	line2Re := regexp.MustCompile(`^([\d\.]+)\s+:.*avg/max\s+=\s+[\d\.]+/([\d\.]+)/[\d\.]+`)

	res := &PingResult{}

	if m := line1Re.FindStringSubmatch(line1); m != nil {
		res.Host = m[1]
		lat, _ := strconv.ParseFloat(m[2], 64)
		avg, _ := strconv.ParseFloat(m[3], 64)
		loss, _ := strconv.Atoi(m[4])

		res.Latency = avg
		res.Alive = (loss == 0 && lat > 0)
		return res, nil
	}

	if m := line2Re.FindStringSubmatch(line2); m != nil {
		res.Host = m[1]
		avg, _ := strconv.ParseFloat(m[2], 64)
		res.Latency = avg
		res.Alive = true
		return res, nil
	}

	return nil, fmt.Errorf("could not parse output: %s", string(out))
}
