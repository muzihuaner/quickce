import type { ProbeRequest, ProbeResult, Agent, DispatchResult, TaskDetail } from '../types/probe'

const API_BASE = '/api'

export async function sendProbe(req: ProbeRequest): Promise<ProbeResult | DispatchResult> {
  const res = await fetch(`${API_BASE}/probe`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || 'request failed')
  }
  return res.json()
}

export async function listAgents(): Promise<Agent[]> {
  const res = await fetch(`${API_BASE}/agents`)
  if (!res.ok) return []
  return res.json()
}

export async function getTask(id: string): Promise<TaskDetail> {
  const res = await fetch(`${API_BASE}/tasks/${id}`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || 'task not found')
  }
  return res.json()
}

export async function waitForTask(id: string, timeoutMs = 30000): Promise<TaskDetail> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    const task = await getTask(id)
    if (task.status === 'done') return task
    await new Promise(r => setTimeout(r, 1500))
  }
  throw new Error('task polling timeout')
}
