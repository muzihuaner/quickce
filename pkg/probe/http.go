package probe

import (
	"context"
	"net/http"
	"time"
)

var sharedClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	},
	Timeout: 0,
}

func ProbeHTTP(ctx context.Context, target string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	url := target
	if !(len(url) > 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "http://" + url
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &Result{Error: err.Error()}
	}

	resp, err := sharedClient.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &Result{Error: err.Error(), LatencyMs: latency}
	}
	defer resp.Body.Close()

	return &Result{
		Success:   true,
		LatencyMs: latency,
		Detail: HTTPDetail{
			StatusCode: resp.StatusCode,
		},
	}
}
