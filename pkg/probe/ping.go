package probe

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ProbePing(ctx context.Context, target string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	start := time.Now()

	args := []string{"-n", "1", "-w", strconv.Itoa(timeoutSec * 1000), target}
	if runtime.GOOS != "windows" {
		args = []string{"-c", "1", "-W", strconv.Itoa(timeoutSec), target}
	}

	cmd := exec.CommandContext(ctx, "ping", args...)
	output, err := cmd.CombinedOutput()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &Result{
			Error:     err.Error() + ": " + strings.TrimSpace(string(output)),
			LatencyMs: latency,
		}
	}

	text := string(output)
	var rtt time.Duration

	if runtime.GOOS == "windows" {
		if idx := strings.Index(text, "平均 = "); idx > 0 {
			part := text[idx+len("平均 = "):]
			if end := strings.IndexAny(part, "ms\r\n"); end > 0 {
				if val, err := strconv.ParseFloat(part[:end], 64); err == nil {
					rtt = time.Duration(val) * time.Millisecond
				}
			}
		} else if idx := strings.Index(text, "Average = "); idx > 0 {
			part := text[idx+len("Average = "):]
			if end := strings.IndexAny(part, "ms\r\n"); end > 0 {
				if val, err := strconv.ParseFloat(part[:end], 64); err == nil {
					rtt = time.Duration(val) * time.Millisecond
				}
			}
		}
	} else {
		if idx := strings.LastIndex(text, "time="); idx > 0 {
			part := text[idx+5:]
			if end := strings.IndexAny(part, " ms\n"); end > 0 {
				if val, err := strconv.ParseFloat(part[:end], 64); err == nil {
					rtt = time.Duration(val) * time.Millisecond
				}
			}
		}
	}

	if !strings.Contains(text, "TTL=") && !strings.Contains(text, "ttl=") && !strings.Contains(text, "time=") {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Destination host unreachable") ||
				strings.Contains(line, "Request timed out") ||
				strings.Contains(line, "100% loss") {
				return &Result{
					Error:     "unreachable: " + strings.TrimSpace(line),
					LatencyMs: latency,
				}
			}
		}
	}

	if rtt == 0 {
		rtt = time.Duration(latency) * time.Millisecond
	}

	return &Result{
		Success:   true,
		LatencyMs: int64(rtt.Milliseconds()),
		Detail: PingDetail{
			RTT: int64(rtt.Milliseconds()),
		},
	}
}
