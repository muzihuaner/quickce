package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ipSBResponse struct {
	IP             string `json:"ip"`
	ISP            string `json:"isp"`
	Organization   string `json:"organization"`
	ASN            int    `json:"asn"`
	ASNOrg         string `json:"asn_organization"`
}

func ProbeIP(ctx context.Context, target string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	ip := target
	if net.ParseIP(target) == nil {
		ips, err := net.LookupHost(target)
		if err != nil || len(ips) == 0 {
			return &Result{
				Error: "cannot resolve domain to IP: " + target,
			}
		}
		ip = ips[0]
	}

	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &Result{Error: err.Error()}
	}

	req.Header.Set("User-Agent", "QuickCE/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &Result{
			Error:     err.Error(),
			LatencyMs: latency,
		}
	}
	defer resp.Body.Close()

	var data ipSBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &Result{
			Error:     "decode response failed: " + err.Error(),
			LatencyMs: latency,
		}
	}

	return &Result{
		Success:   true,
		LatencyMs: latency,
		Detail: IPDetail{
			IP:     data.IP,
			ISP:    data.ISP,
			Org:    data.Organization,
			ASN:    data.ASN,
			ASNOrg: data.ASNOrg,
		},
	}
}
