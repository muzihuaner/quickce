package agent

import (
	"context"
	"time"

	"quickce/pkg/probe"
)

type probeResult struct {
	Success   bool        `json:"success"`
	LatencyMs int64       `json:"latency_ms"`
	Error     string      `json:"error,omitempty"`
	Detail    interface{} `json:"detail,omitempty"`
}

func executeProbe(ctx context.Context, target, probeType string, port, timeout int) *probeResult {
	if timeout <= 0 {
		timeout = 5
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	var r *probe.Result
	switch probe.ProbeType(probeType) {
	case probe.HTTP:
		r = probe.ProbeHTTP(ctx, target, timeout)
	case probe.Ping:
		r = probe.ProbePing(ctx, target, timeout)
	case probe.TCP:
		if port == 0 {
			port = 80
		}
		r = probe.ProbeTCP(ctx, target, port, timeout)
	case probe.DNS:
		r = probe.ProbeDNS(ctx, target, timeout)
	case probe.Traceroute:
		r = probe.ProbeTraceroute(ctx, target, timeout)
	case probe.IP:
		r = probe.ProbeIP(ctx, target, timeout)
	default:
		return &probeResult{Error: "unsupported type: " + probeType}
	}

	if r == nil {
		return &probeResult{Error: "no result"}
	}

	return &probeResult{
		Success:   r.Success,
		LatencyMs: r.LatencyMs,
		Error:     r.Error,
		Detail:    r.Detail,
	}
}
