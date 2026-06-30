<script setup lang="ts">
import type { PingDetail } from '../types/probe'

defineProps<{ detail: PingDetail }>()
</script>

<template>
  <div class="detail">
    <span class="label">往返延迟</span>
    <span class="value">{{ detail.rtt_ms }} ms</span>
  </div>
</template>

<style scoped>
.detail { display: flex; gap: 8px; }
.label { color: #64748b; font-size: 13px; }
.value { font-weight: 600; font-size: 14px; color: #0284c7; }
</style>
