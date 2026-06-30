package probe

import "time"

type ProbeType string

const (
	HTTP       ProbeType = "http"
	Ping       ProbeType = "ping"
	TCP        ProbeType = "tcp"
	DNS        ProbeType = "dns"
	Traceroute ProbeType = "traceroute"
	IP         ProbeType = "ip"
)

type Result struct {
	Success   bool        `json:"success"`
	LatencyMs int64       `json:"latency_ms"`
	Error     string      `json:"error,omitempty"`
	Detail    interface{} `json:"detail,omitempty"`
}

type HTTPDetail struct {
	StatusCode int    `json:"status_code"`
}

type PingDetail struct {
	RTT int64 `json:"rtt_ms"`
}

type TCPDetail struct {
	PortOpen bool  `json:"port_open"`
	RTT      int64 `json:"rtt_ms"`
}

type DNSDetail struct {
	Ips []string `json:"ips"`
}

type TracerouteDetail struct {
	Hops []Hop `json:"hops"`
}

type IPDetail struct {
	IP         string `json:"ip"`
	ISP        string `json:"isp,omitempty"`
	Org        string `json:"organization,omitempty"`
	ASN        int    `json:"asn,omitempty"`
	ASNOrg     string `json:"asn_organization,omitempty"`
}

type Hop struct {
	TTL     int           `json:"ttl"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
}
