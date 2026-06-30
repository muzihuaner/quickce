package probe

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ProbeTraceroute(ctx context.Context, target string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	start := time.Now()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "tracert", "-d", "-h", "5", target)
	} else {
		cmd = exec.CommandContext(ctx, "traceroute", "-n", "-m", "5", target)
	}

	output, err := cmd.CombinedOutput()
	totalLatency := time.Since(start).Milliseconds()

	hops := parseTraceroute(string(output), runtime.GOOS == "windows")

	if err != nil && len(hops) == 0 {
		return &Result{
			Error:     err.Error() + ": " + strings.TrimSpace(string(output)),
			LatencyMs: totalLatency,
			Detail: TracerouteDetail{
				Hops: hops,
			},
		}
	}

	return &Result{
		Success:   true,
		LatencyMs: totalLatency,
		Detail: TracerouteDetail{
			Hops: hops,
		},
	}
}

func parseTraceroute(output string, isWindows bool) []Hop {
	lines := strings.Split(output, "\n")
	var hops []Hop

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if isWindows {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			ttl, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}
			if ttl < 1 || ttl > 5 {
				continue
			}
			if len(fields) < 3 {
				continue
			}
			if fields[1] == "*" || fields[2] == "*" {
				hops = append(hops, Hop{TTL: ttl, Address: "*", RTT: 0})
				continue
			}
			addr := fields[len(fields)-1]
			if addr == "" || addr == "Request" || addr == "ms" {
				continue
			}
			rttStr := fields[1]
			rttVal := 0.0
			if strings.HasSuffix(rttStr, "ms") {
				rttVal, _ = strconv.ParseFloat(rttStr[:len(rttStr)-2], 64)
			}
			hops = append(hops, Hop{
				TTL:     ttl,
				Address: addr,
				RTT:     time.Duration(rttVal) * time.Millisecond,
			})
		} else {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			ttl, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}
			if ttl < 1 || ttl > 5 {
				continue
			}
			if fields[1] == "*" {
				hops = append(hops, Hop{TTL: ttl, Address: "*", RTT: 0})
				continue
			}
			addr := fields[1]
			rttVal := 0.0
			for _, f := range fields[2:] {
				if strings.HasSuffix(f, "ms") {
					rttVal, _ = strconv.ParseFloat(f[:len(f)-2], 64)
					break
				}
			}
			hops = append(hops, Hop{
				TTL:     ttl,
				Address: addr,
				RTT:     time.Duration(rttVal) * time.Millisecond,
			})
		}
	}

	return hops
}
