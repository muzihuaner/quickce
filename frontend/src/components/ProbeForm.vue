<script setup lang="ts">
import { reactive } from 'vue'
import type { ProbeType, Agent } from '../types/probe'

const props = defineProps<{ agents: Agent[] }>()

const emit = defineEmits<{
  submit: [payload: {
    target: string
    type: ProbeType
    port: number
    timeout: number
    nodes: string[]
  }]
}>()

const PROBE_TYPES: { value: ProbeType; label: string }[] = [
  { value: 'http', label: 'HTTP/HTTPS' },
  { value: 'ping', label: 'Ping (ICMP)' },
  { value: 'tcp', label: 'TCP 端口' },
  { value: 'dns', label: 'DNS 解析' },
  { value: 'traceroute', label: '路由追踪' },
  { value: 'ip', label: 'IP 查询' },
]

const form = reactive({
  target: '',
  type: 'http' as ProbeType,
  port: 80,
  timeout: 5,
  nodes: [] as string[],
})

function toggleNode(name: string) {
  const idx = form.nodes.indexOf(name)
  if (idx >= 0) {
    form.nodes.splice(idx, 1)
  } else {
    form.nodes.push(name)
  }
}

function handleSubmit() {
  if (!form.target.trim()) return
  emit('submit', { ...form, target: form.target.trim() })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="probe-form">
    <div class="form-row">
      <div class="field target-field">
        <label for="target">目标地址</label>
        <input
          id="target"
          v-model="form.target"
          type="text"
          placeholder="域名或 IP 地址，如 example.com"
          autocomplete="off"
        />
      </div>
      <div class="field type-field">
        <label for="type">探测类型</label>
        <select id="type" v-model="form.type">
          <option v-for="t in PROBE_TYPES" :key="t.value" :value="t.value">
            {{ t.label }}
          </option>
        </select>
      </div>
      <div v-if="form.type === 'tcp'" class="field port-field">
        <label for="port">端口号</label>
        <input id="port" v-model.number="form.port" type="number" min="1" max="65535" />
      </div>
      <div class="field timeout-field">
        <label for="timeout">超时 (秒)</label>
        <input id="timeout" v-model.number="form.timeout" type="number" min="1" max="60" />
      </div>
    </div>

    <div v-if="agents.length" class="nodes-section">
      <span class="nodes-label">分发节点（可选，勾选后将由远端节点执行）</span>
      <div class="nodes-list">
        <label
          v-for="a in agents"
          :key="a.id"
          class="node-chip"
          :class="{ selected: form.nodes.includes(a.name) }"
        >
          <input
            type="checkbox"
            :checked="form.nodes.includes(a.name)"
            @change="toggleNode(a.name)"
          />
          {{ a.name }}
          <span class="node-loc">{{ a.location }}</span>
        </label>
      </div>
    </div>

    <button type="submit" class="submit-btn">
      {{ form.nodes.length ? '分发探测' : '开始探测' }}
    </button>
  </form>
</template>

<style scoped>
.probe-form {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 24px;
}

.form-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.target-field {
  flex: 2;
  min-width: 240px;
}

.type-field { flex: 1; min-width: 140px; }
.port-field { flex: 0 0 100px; }
.timeout-field { flex: 0 0 100px; }

.field label {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.field input,
.field select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  background: #fff;
}

.field input:focus,
.field select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.nodes-section {
  margin-bottom: 16px;
}

.nodes-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 8px;
}

.nodes-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.node-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  user-select: none;
  color: #475569;
  transition: all 0.15s;
}

.node-chip:hover {
  border-color: #3b82f6;
  background: #f0f5ff;
}

.node-chip.selected {
  border-color: #3b82f6;
  background: #eff6ff;
  color: #1d4ed8;
}

.node-chip input {
  display: none;
}

.node-loc {
  font-size: 11px;
  color: #94a3b8;
  background: #f1f5f9;
  padding: 0 4px;
  border-radius: 2px;
}

.submit-btn {
  background: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}

.submit-btn:hover {
  background: #2563eb;
}

.submit-btn:active {
  background: #1d4ed8;
}
</style>
