package probe

import (
	"context"
	"net"
	"strconv"
	"time"
)

func ProbeTCP(ctx context.Context, target string, port int, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	addr := net.JoinHostPort(target, strconv.Itoa(port))

	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, time.Duration(timeoutSec)*time.Second)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &Result{
			Success:   false,
			LatencyMs: latency,
			Error:     err.Error(),
			Detail: TCPDetail{
				PortOpen: false,
				RTT:      latency,
			},
		}
	}
	conn.Close()

	return &Result{
		Success:   true,
		LatencyMs: latency,
		Detail: TCPDetail{
			PortOpen: true,
			RTT:      latency,
		},
	}
}
