<script setup lang="ts">
import type { IPDetail } from '../types/probe'

defineProps<{ detail: IPDetail }>()
</script>

<template>
  <table class="ip-table">
    <tbody>
      <tr>
        <td class="label">IP 地址</td>
        <td class="value">{{ detail.ip }}</td>
      </tr>
      <tr v-if="detail.isp">
        <td class="label">ISP</td>
        <td class="value">{{ detail.isp }}</td>
      </tr>
      <tr v-if="detail.organization">
        <td class="label">所属组织</td>
        <td class="value">{{ detail.organization }}</td>
      </tr>
      <tr v-if="detail.asn">
        <td class="label">ASN</td>
        <td class="value">AS{{ detail.asn }}</td>
      </tr>
      <tr v-if="detail.asn_organization">
        <td class="label">ASN 组织</td>
        <td class="value">{{ detail.asn_organization }}</td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
.ip-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.ip-table td {
  padding: 6px 0;
}
.label {
  color: #64748b;
  font-weight: 600;
  width: 80px;
  white-space: nowrap;
  vertical-align: top;
}
.value {
  color: #334155;
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
}
</style>
