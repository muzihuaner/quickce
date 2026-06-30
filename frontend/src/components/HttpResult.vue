<script setup lang="ts">
import type { HTTPDetail } from '../types/probe'

defineProps<{ detail: HTTPDetail }>()
</script>

<template>
  <div class="detail">
    <span class="label">状态码</span>
    <span class="value" :class="{ ok: detail.status_code < 400, err: detail.status_code >= 400 }">
      {{ detail.status_code }}
    </span>
  </div>
</template>

<style scoped>
.detail { display: flex; gap: 8px; }
.label { color: #64748b; font-size: 13px; }
.value { font-weight: 600; font-size: 14px; }
.ok { color: #16a34a; }
.err { color: #dc2626; }
</style>
