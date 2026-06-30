<script setup lang="ts">
import type { ProbeResult, HTTPDetail, PingDetail, TCPDetail, DNSDetail, TracerouteDetail, IPDetail, DispatchResult } from '../types/probe'
import HttpResult from './HttpResult.vue'
import PingResult from './PingResult.vue'
import TcpResult from './TcpResult.vue'
import DnsResult from './DnsResult.vue'
import TracerouteResult from './TracerouteResult.vue'
import IpResult from './IpResult.vue'

defineProps<{
  result: ProbeResult | DispatchResult | null
  loading: boolean
}>()

function isHttp(d: unknown): d is HTTPDetail { return d !== undefined && 'status_code' in (d as HTTPDetail) }
function isPing(d: unknown): d is PingDetail { return d !== undefined && 'rtt_ms' in (d as PingDetail) }
function isTcp(d: unknown): d is TCPDetail { return d !== undefined && 'port_open' in (d as TCPDetail) }
function isDns(d: unknown): d is DNSDetail { return d !== undefined && 'ips' in (d as DNSDetail) }
function isTraceroute(d: unknown): d is TracerouteDetail { return d !== undefined && 'hops' in (d as TracerouteDetail) }
function isIp(d: unknown): d is IPDetail { return d !== undefined && 'ip' in (d as IPDetail) }

function isDispatch(r: unknown): r is DispatchResult {
  return r !== null && typeof r === 'object' && 'dispatched' in r
}

function isProbe(r: unknown): r is ProbeResult {
  return r !== null && typeof r === 'object' && 'success' in r
}
</script>

<template>
  <div class="result-card">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>探测中...</span>
    </div>
    <div v-else-if="!result" class="placeholder">
      输入目标地址并点击「开始探测」发起网络拨测
    </div>

    <template v-else-if="isDispatch(result)">
      <div class="dispatch-header">
        <div class="dispatch-icon">↗</div>
        <div>
          <div class="dispatch-title">任务已分发</div>
          <div class="dispatch-msg">{{ result.message }}</div>
        </div>
      </div>
      <div class="dispatch-tasks">
        <div v-for="t in result.tasks" :key="t.task_id" class="dispatch-task">
          <div class="task-left">
            <span class="task-agent">{{ t.agent }}</span>
            <span class="task-location">{{ t.location }}</span>
          </div>
          <div class="task-id">{{ t.task_id.slice(0, 8) }}...</div>
        </div>
      </div>
    </template>

    <template v-else-if="isProbe(result)">
      <div class="status" :class="{ success: result.success, fail: !result.success }">
        {{ result.success ? '成功' : '失败' }}
      </div>
      <div class="meta">
        <div class="meta-item">
          <span class="label">延迟</span>
          <span class="value">{{ result.latency_ms }} ms</span>
        </div>
      </div>
      <div v-if="result.error" class="error-msg">{{ result.error }}</div>
      <div v-if="result.detail" class="detail-section">
        <HttpResult v-if="isHttp(result.detail)" :detail="result.detail" />
        <PingResult v-if="isPing(result.detail)" :detail="result.detail" />
        <TcpResult v-if="isTcp(result.detail)" :detail="result.detail" />
        <DnsResult v-if="isDns(result.detail)" :detail="result.detail" />
        <TracerouteResult v-if="isTraceroute(result.detail)" :detail="result.detail" />
        <IpResult v-if="isIp(result.detail)" :detail="result.detail" />
      </div>
    </template>
  </div>
</template>

<style scoped>
.result-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 24px;
  min-height: 120px;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #64748b;
  font-size: 14px;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.placeholder { color: #94a3b8; font-size: 14px; text-align: center; padding: 20px 0; }

.status {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 16px;
}

.success { background: #dcfce7; color: #16a34a; }
.fail { background: #fee2e2; color: #dc2626; }

.meta { display: flex; gap: 24px; margin-bottom: 12px; }
.meta-item { display: flex; gap: 8px; }
.meta .label { color: #64748b; font-size: 13px; }
.meta .value { font-weight: 600; font-size: 14px; color: #334155; }

.error-msg {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 13px;
  margin-bottom: 12px;
}

.detail-section {
  border-top: 1px solid #f1f5f9;
  padding-top: 12px;
}

.dispatch-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.dispatch-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #eff6ff;
  color: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
}

.dispatch-title {
  font-size: 15px;
  font-weight: 600;
  color: #0f172a;
}

.dispatch-msg {
  font-size: 13px;
  color: #64748b;
}

.dispatch-tasks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dispatch-task {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: #f8fafc;
  border-radius: 6px;
}

.task-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-agent {
  font-weight: 600;
  font-size: 13px;
  color: #0f172a;
}

.task-location {
  font-size: 12px;
  color: #64748b;
  background: #e2e8f0;
  padding: 1px 6px;
  border-radius: 3px;
}

.task-id {
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
  color: #94a3b8;
}
</style>
