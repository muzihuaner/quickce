package probe

import (
	"context"
	"net"
	"time"
)

func ProbeDNS(ctx context.Context, target string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	r := &net.Resolver{}

	start := time.Now()
	ips, err := r.LookupHost(ctx, target)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &Result{
			LatencyMs: latency,
			Error:     err.Error(),
			Detail: DNSDetail{
				Ips: []string{},
			},
		}
	}

	if len(ips) == 0 {
		return &Result{
			LatencyMs: latency,
			Error:     "no records found",
			Detail: DNSDetail{
				Ips: []string{},
			},
		}
	}

	return &Result{
		Success:   true,
		LatencyMs: latency,
		Detail: DNSDetail{
			Ips: ips,
		},
	}
}
