<script setup lang="ts">
import type { TracerouteDetail } from '../types/probe'

defineProps<{ detail: TracerouteDetail }>()
</script>

<template>
  <div>
    <span class="label">路由跳数</span>
    <table class="hop-table">
      <thead>
        <tr>
          <th>跳数</th>
          <th>地址</th>
          <th>延迟</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="hop in detail.hops" :key="hop.ttl">
          <td>{{ hop.ttl }}</td>
          <td>{{ hop.address }}</td>
          <td>{{ hop.rtt }} ms</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.label { color: #64748b; font-size: 13px; display: block; margin-bottom: 8px; }
.hop-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.hop-table th {
  text-align: left;
  color: #64748b;
  font-weight: 600;
  padding: 6px 12px;
  border-bottom: 1px solid #e2e8f0;
}
.hop-table td {
  padding: 6px 12px;
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
}
</style>
