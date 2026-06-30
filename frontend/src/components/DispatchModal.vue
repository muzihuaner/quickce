<script setup lang="ts">
import type { DispatchTask } from '../types/probe'

defineProps<{
  show: boolean
  tasks: DispatchTask[]
}>()

const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="overlay" @click.self="emit('close')">
      <div class="modal">
        <div class="modal-header">
          <div class="modal-icon">✓</div>
          <div>
            <div class="modal-title">任务已分发</div>
            <div class="modal-msg">探针任务已下发到以下节点，结果将实时展示</div>
          </div>
        </div>
        <div class="modal-body">
          <div v-for="t in tasks" :key="t.task_id" class="modal-task">
            <span class="modal-agent">{{ t.agent }}</span>
            <span class="modal-location">{{ t.location }}</span>
            <span class="modal-id">{{ t.task_id.slice(0, 8) }}...</span>
          </div>
        </div>
        <div class="modal-footer">
          <button class="modal-btn" @click="emit('close')">知道了</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  border-radius: 12px;
  width: 420px;
  max-width: 90vw;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 24px 24px 0;
}

.modal-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #dcfce7;
  color: #16a34a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  flex-shrink: 0;
}

.modal-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
}

.modal-msg {
  font-size: 13px;
  color: #64748b;
  margin-top: 2px;
}

.modal-body {
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.modal-task {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: #f8fafc;
  border-radius: 6px;
}

.modal-agent {
  font-weight: 600;
  font-size: 13px;
  color: #0f172a;
}

.modal-location {
  font-size: 12px;
  color: #64748b;
  background: #e2e8f0;
  padding: 1px 6px;
  border-radius: 3px;
}

.modal-id {
  margin-left: auto;
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
  color: #94a3b8;
}

.modal-footer {
  padding: 0 24px 20px;
}

.modal-btn {
  width: 100%;
  padding: 10px;
  border: none;
  border-radius: 8px;
  background: #3b82f6;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.modal-btn:hover {
  background: #2563eb;
}
</style>
