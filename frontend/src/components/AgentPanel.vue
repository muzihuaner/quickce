<script setup lang="ts">
import type { Agent } from '../types/probe'

defineProps<{ agents: Agent[] }>()
</script>

<template>
  <div class="agent-panel">
    <div class="agent-header">
      <span class="agent-title">探测节点</span>
      <span v-if="agents.length" class="badge">{{ agents.length }}</span>
      <span v-else class="badge offline">0</span>
    </div>
    <div class="agent-body">
      <div v-if="agents.length === 0" class="empty">暂无在线节点，启动 Agent 客户端即可自动注册</div>
      <div v-for="a in agents" :key="a.id" class="agent-item">
        <div class="agent-info">
          <span class="agent-name">{{ a.name }}</span>
          <span class="agent-location">{{ a.location }}</span>
        </div>
        <div class="agent-meta">
          <span class="agent-ip">{{ a.ip }}</span>
          <span class="dot"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.agent-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.agent-header {
  padding: 12px 24px;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid #f1f5f9;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 10px;
  background: #3b82f6;
  color: #fff;
  font-size: 11px;
  font-weight: 700;
}

.badge.offline {
  background: #e2e8f0;
  color: #94a3b8;
}

.agent-body {
  padding: 12px 24px;
}

.empty {
  color: #94a3b8;
  font-size: 13px;
  text-align: center;
  padding: 8px 0;
}

.agent-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8fafc;
  border-radius: 6px;
  margin-bottom: 6px;
}

.agent-item:last-child {
  margin-bottom: 0;
}

.agent-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.agent-name {
  font-weight: 600;
  font-size: 13px;
  color: #0f172a;
}

.agent-location {
  font-size: 12px;
  color: #64748b;
  background: #e2e8f0;
  padding: 1px 6px;
  border-radius: 3px;
}

.agent-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.agent-ip {
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
  color: #94a3b8;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #22c55e;
}
</style>
