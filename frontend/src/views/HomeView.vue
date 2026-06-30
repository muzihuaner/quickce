<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { ProbeResult, Agent, DispatchResult, DispatchTask, TaskDetail, HTTPDetail, PingDetail, TCPDetail, DNSDetail, TracerouteDetail, IPDetail } from '../types/probe'
import { sendProbe, listAgents, waitForTask } from '../api/probe'
import ProbeForm from '../components/ProbeForm.vue'
import ResultCard from '../components/ResultCard.vue'
import AgentPanel from '../components/AgentPanel.vue'
import HttpResult from '../components/HttpResult.vue'
import PingResult from '../components/PingResult.vue'
import TcpResult from '../components/TcpResult.vue'
import DnsResult from '../components/DnsResult.vue'
import TracerouteResult from '../components/TracerouteResult.vue'
import IpResult from '../components/IpResult.vue'
import DispatchModal from '../components/DispatchModal.vue'

function isHttp(d: unknown): d is HTTPDetail { return d !== undefined && 'status_code' in (d as HTTPDetail) }
function isPing(d: unknown): d is PingDetail { return d !== undefined && 'rtt_ms' in (d as PingDetail) }
function isTcp(d: unknown): d is TCPDetail { return d !== undefined && 'port_open' in (d as TCPDetail) }
function isDns(d: unknown): d is DNSDetail { return d !== undefined && 'ips' in (d as DNSDetail) }
function isTraceroute(d: unknown): d is TracerouteDetail { return d !== undefined && 'hops' in (d as TracerouteDetail) }
function isIp(d: unknown): d is IPDetail { return d !== undefined && 'ip' in (d as IPDetail) }

const result = ref<ProbeResult | DispatchResult | null>(null)
const loading = ref(false)
const errorBanner = ref('')
const agents = ref<Agent[]>([])
const taskResults = ref<TaskDetail[]>([])
const showDispatchModal = ref(false)
const dispatchTasks = ref<DispatchTask[]>([])
const expanded = ref<Set<string>>(new Set())

function agentName(t: TaskDetail): string {
  return dispatchTasks.value.find(d => d.task_id === t.id)?.agent || t.agent_id.slice(0, 8)
}

function toggleExpand(id: string) {
  const s = new Set(expanded.value)
  if (s.has(id)) s.delete(id); else s.add(id)
  expanded.value = s
}

async function fetchAgents() {
  try {
    agents.value = await listAgents()
  } catch {
    agents.value = []
  }
}

onMounted(fetchAgents)

async function onSubmit(payload: {
  target: string
  type: 'http' | 'ping' | 'tcp' | 'dns' | 'traceroute' | 'ip'
  port: number
  timeout: number
  nodes: string[]
}) {
  loading.value = true
  result.value = null
  taskResults.value = []
  errorBanner.value = ''

  try {
    const resp = await sendProbe({
      target: payload.target,
      type: payload.type,
      port: payload.type === 'tcp' ? payload.port : undefined,
      timeout: payload.timeout,
      ...(payload.nodes.length ? { nodes: payload.nodes } : {}),
    })

    if ('dispatched' in resp) {
      dispatchTasks.value = resp.tasks
      showDispatchModal.value = true
      const tasks = resp.tasks.map(t => waitForTask(t.task_id))
      taskResults.value = await Promise.all(tasks)
    } else {
      result.value = resp
    }
  } catch (e) {
    errorBanner.value = e instanceof Error ? e.message : '未知错误'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="home">
    <h1 class="title">QuickCE</h1>
    <p class="subtitle">网络拨测工具</p>

    <AgentPanel :agents="agents" />

    <ProbeForm :agents="agents" @submit="onSubmit" />

    <details class="guide">
      <summary>使用说明</summary>
      <div class="guide-content">
        <p><strong>支持的探测类型：</strong></p>
        <table>
          <thead>
            <tr><th>类型</th><th>说明</th><th>必填参数</th></tr>
          </thead>
          <tbody>
            <tr><td>HTTP/HTTPS</td><td>发送 GET 请求，检测网站可访问性及响应状态码</td><td>目标地址</td></tr>
            <tr><td>Ping (ICMP)</td><td>检测目标是否可达，返回往返延迟（RTT）</td><td>目标地址</td></tr>
            <tr><td>TCP 端口</td><td>检测指定端口是否开放，测量连接耗时</td><td>目标地址 + 端口号</td></tr>
            <tr><td>DNS 解析</td><td>查询域名的 A 记录，返回解析到的 IP 列表</td><td>目标地址</td></tr>
            <tr><td>路由追踪</td><td>追踪到目标的路由路径，显示前 5 跳信息</td><td>目标地址</td></tr>
            <tr><td>IP 查询</td><td>查询 IP 的地理位置和 ISP 信息（数据来源 api.ip.sb）</td><td>IP 地址或域名</td></tr>
          </tbody>
        </table>
        <p class="guide-note">
          超时时间默认为 5 秒，可根据实际情况调整。<br>
          路由追踪耗时较长，建议设置 30 秒以上超时。<br>
          勾选「分发节点」可将探测任务下发到远端 Agent 执行。
        </p>
      </div>
    </details>

    <div v-if="errorBanner" class="error-banner">{{ errorBanner }}</div>

    <DispatchModal
      :show="showDispatchModal"
      :tasks="dispatchTasks"
      @close="showDispatchModal = false"
    />

    <div v-if="taskResults.length" class="task-results">
      <h3 class="section-title">节点探测结果</h3>
      <table class="result-table">
        <thead>
          <tr>
            <th>节点</th>
            <th>状态</th>
            <th>延迟</th>
            <th>详情</th>
            <th>错误</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in taskResults" :key="t.id">
            <td class="col-node">
              {{ dispatchTasks.find(d => d.task_id === t.id)?.agent || t.agent_id.slice(0, 8) }}
            </td>
            <td>
              <span class="tag" :class="{ ok: t.result?.success, fail: !t.result?.success }">
                {{ t.result ? (t.result.success ? '成功' : '失败') : '等待中' }}
              </span>
            </td>
            <td class="col-latency">{{ t.result ? t.result.latency_ms + ' ms' : '-' }}</td>
            <td class="col-detail">
              <template v-if="t.result?.detail">
                <span v-if="isHttp(t.result.detail)">状态码 {{ t.result.detail.status_code }}</span>
                <span v-if="isPing(t.result.detail)">{{ t.result.detail.rtt_ms }} ms</span>
                <span v-if="isTcp(t.result.detail)">{{ t.result.detail.port_open ? '开放' : '关闭' }} ({{ t.result.detail.rtt_ms }} ms)</span>
                <span v-if="isDns(t.result.detail)">{{ t.result.detail.ips.length }} 条记录</span>
                <span v-if="isIp(t.result.detail)">{{ t.result.detail.isp || t.result.detail.organization || t.result.detail.ip }}</span>
                <span v-if="isTraceroute(t.result.detail)">{{ t.result.detail.hops.length }} 跳</span>
              </template>
              <span v-else class="muted">-</span>
            </td>
            <td class="col-error">
              <span v-if="t.result?.error" class="err-text">{{ t.result.error }}</span>
              <span v-else class="muted">-</span>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-for="t in taskResults" :key="'d'+t.id" class="detail-expand">
        <div
          v-if="t.result?.detail"
          class="expand-toggle"
          @click="toggleExpand(t.id)"
        >
          {{ expanded.has(t.id) ? '收起' : '查看' }} {{ agentName(t) }} 完整结果
        </div>
        <div v-if="expanded.has(t.id) && t.result?.detail" class="expand-body">
          <HttpResult v-if="isHttp(t.result.detail)" :detail="t.result.detail" />
          <PingResult v-if="isPing(t.result.detail)" :detail="t.result.detail" />
          <TcpResult v-if="isTcp(t.result.detail)" :detail="t.result.detail" />
          <DnsResult v-if="isDns(t.result.detail)" :detail="t.result.detail" />
          <TracerouteResult v-if="isTraceroute(t.result.detail)" :detail="t.result.detail" />
          <IpResult v-if="isIp(t.result.detail)" :detail="t.result.detail" />
        </div>
      </div>
    </div>

    <ResultCard :result="result" :loading="loading" />
  </div>
</template>

<style scoped>
.home {
  max-width: 720px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
  text-align: center;
}

.subtitle {
  font-size: 14px;
  color: #64748b;
  text-align: center;
  margin: -12px 0 0;
}

.error-banner {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  padding: 10px 16px;
  border-radius: 6px;
  font-size: 14px;
}

.guide {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.guide summary {
  padding: 12px 24px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  user-select: none;
}

.guide summary:hover {
  background: #f8fafc;
}

.guide-content {
  padding: 0 24px 16px;
  font-size: 14px;
  color: #334155;
}

.guide-content p {
  margin: 8px 0;
}

.guide-content table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.guide-content th {
  text-align: left;
  color: #64748b;
  font-weight: 600;
  padding: 8px 12px;
  border-bottom: 2px solid #e2e8f0;
}

.guide-content td {
  padding: 8px 12px;
  border-bottom: 1px solid #f1f5f9;
}

.guide-note {
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.8;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #0f172a;
  margin: 0 0 12px;
}

.task-results {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.result-table thead {
  background: #f8fafc;
}

.result-table th {
  text-align: left;
  color: #64748b;
  font-weight: 600;
  padding: 10px 14px;
  border-bottom: 1px solid #e2e8f0;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.result-table td {
  padding: 10px 14px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
}

.result-table tbody tr:last-child td {
  border-bottom: none;
}

.result-table tbody tr:hover {
  background: #f8fafc;
}

.col-node {
  font-weight: 600;
  color: #0f172a;
  white-space: nowrap;
}

.col-latency {
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
}

.col-error {
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.tag.ok { background: #dcfce7; color: #16a34a; }
.tag.fail { background: #fee2e2; color: #dc2626; }

.err-text {
  color: #dc2626;
  font-size: 12px;
}

.muted {
  color: #cbd5e1;
}

.detail-expand {
  margin-top: 4px;
}

.expand-toggle {
  font-size: 13px;
  color: #3b82f6;
  cursor: pointer;
  padding: 4px 0;
  user-select: none;
}

.expand-toggle:hover {
  color: #2563eb;
}

.expand-body {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 14px;
  margin-top: 4px;
}

.expand-body :deep(.detail) {
  margin-bottom: 4px;
}
</style>
