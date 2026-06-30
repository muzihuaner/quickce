<script setup lang="ts">
import type { TCPDetail } from '../types/probe'

defineProps<{ detail: TCPDetail }>()
</script>

<template>
  <div class="detail">
    <span class="label">端口状态</span>
    <span class="value" :class="{ open: detail.port_open, closed: !detail.port_open }">
      {{ detail.port_open ? '开放' : '关闭' }}
    </span>
  </div>
  <div class="detail">
    <span class="label">连接耗时</span>
    <span class="value">{{ detail.rtt_ms }} ms</span>
  </div>
</template>

<style scoped>
.detail { display: flex; gap: 8px; margin-bottom: 4px; }
.label { color: #64748b; font-size: 13px; }
.value { font-weight: 600; font-size: 14px; }
.open { color: #16a34a; }
.closed { color: #dc2626; }
</style>
