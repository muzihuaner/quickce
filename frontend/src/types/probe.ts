export type ProbeType = 'http' | 'ping' | 'tcp' | 'dns' | 'traceroute' | 'ip'

export interface ProbeRequest {
  target: string
  type: ProbeType
  port?: number
  timeout?: number
}

export interface HTTPDetail {
  status_code: number
}

export interface PingDetail {
  rtt_ms: number
}

export interface TCPDetail {
  port_open: boolean
  rtt_ms: number
}

export interface DNSDetail {
  ips: string[]
}

export interface Hop {
  ttl: number
  address: string
  rtt: number
}

export interface IPDetail {
  ip: string
  isp?: string
  organization?: string
  asn?: number
  asn_organization?: string
}

export interface TracerouteDetail {
  hops: Hop[]
}

export interface ProbeResult {
  success: boolean
  latency_ms: number
  error?: string
  detail?: HTTPDetail | PingDetail | TCPDetail | DNSDetail | TracerouteDetail | IPDetail
}

export interface Agent {
  id: string
  name: string
  location: string
  ip: string
  last_seen: string
  created_at: string
}

export interface DispatchTask {
  task_id: string
  agent_id: string
  agent: string
  location: string
}

export interface DispatchResult {
  dispatched: boolean
  tasks: DispatchTask[]
  message: string
}

export interface TaskDetail {
  id: string
  agent_id: string
  status: string
  target: string
  type: string
  port: number
  timeout: number
  created_at: string
  done_at: string | null
  result: ProbeResult | null
}
