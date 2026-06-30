<script setup lang="ts">
import type { DNSDetail } from '../types/probe'

defineProps<{ detail: DNSDetail }>()
</script>

<template>
  <div class="detail">
    <span class="label">解析结果</span>
    <ul class="ip-list">
      <li v-for="ip in detail.ips" :key="ip">{{ ip }}</li>
    </ul>
  </div>
</template>

<style scoped>
.detail { display: flex; gap: 8px; }
.label { color: #64748b; font-size: 13px; white-space: nowrap; padding-top: 2px; }
.ip-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 4px; }
.ip-list li {
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 13px;
  color: #0284c7;
  background: #f0f9ff;
  padding: 2px 8px;
  border-radius: 4px;
}
</style>
