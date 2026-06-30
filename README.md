# QuickCE - 网络拨测工具

QuickCE 是一个轻量级的网络拨测服务，支持 HTTP、Ping、TCP 端口、DNS 解析、路由追踪和 IP 查询六种探测类型。提供 RESTful API 和 Web 界面。

## 技术栈

- **后端**: Go 1.26 + Gin
- **前端**: Vue 3 + TypeScript + Vite
- **并发控制**: `golang.org/x/sync/semaphore`
- **超时控制**: `context.WithTimeout`

## 项目结构

```
quickce/
├── cmd/
│   ├── server/main.go          # 后端入口（含分布式调度）
│   └── agent/main.go           # Agent 客户端入口
├── pkg/
│   ├── api/
│   │   ├── handler.go          # API 路由和请求处理（本地+远程调度）
│   │   └── agent.go            # Agent 注册、任务轮询、结果上报 API
│   ├── config/config.go        # 环境变量配置
│   ├── store/store.go          # Agent 和任务存储（内存，可替换为数据库）
│   ├── agent/
│   │   ├── client.go           # Agent 客户端（注册、轮询、执行、上报）
│   │   └── probe.go            # Agent 端探测执行封装
│   └── probe/
│       ├── probe.go            # 公共类型定义
│       ├── http.go             # HTTP/HTTPS 探测
│       ├── ping.go             # ICMP Ping 探测
│       ├── tcp.go              # TCP 端口探测
│       ├── dns.go              # DNS 解析探测
│       ├── ip.go               # IP 地理位置查询
│       └── traceroute.go       # 路由追踪探测
└── frontend/
    ├── src/
    │   ├── api/probe.ts        # API 调用
    │   ├── types/probe.ts      # TypeScript 类型
    │   ├── views/HomeView.vue  # 首页
    │   └── components/         # 探测表单和各类型结果组件
    ├── vite.config.ts
    └── package.json
```

## 探测类型

| 类型 | API type | 说明 | 关键实现 |
|------|----------|------|----------|
| HTTP/HTTPS | `http` | GET 请求，返回状态码和响应时间 | 复用 `http.Transport` 连接池 |
| Ping | `ping` | ICMP Echo，返回可达性和 RTT | 调用系统 `ping` 命令，解析输出 |
| TCP 端口 | `tcp` | 检测端口是否开放 | `net.DialTimeout` |
| DNS 解析 | `dns` | 查询 A 记录 | `net.Resolver.LookupHost` |
| 路由追踪 | `traceroute` | 返回前 5 跳 IP 和延迟 | 调用系统 `tracert` / `traceroute` |
| IP 查询 | `ip` | 查询 IP 地理位置和 ISP 信息 | `api.ip.sb/geoip` API |

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+（仅前端开发时需要）

### 运行后端

```bash
# 编译并启动
go build -o quickce-server ./cmd/server
./quickce-server

# 或直接运行
go run ./cmd/server
```

服务默认监听 `:8080`，访问 `http://localhost:8080`。

### 运行前端（开发模式）

```bash
cd frontend
npm install
npm run dev
```

前端 dev server 默认监听 `:5173`，已配置代理将 `/api` 请求转发到后端 `:8080`。

### 单二进制部署（前后端整合）

```bash
# 构建前端
cd frontend
npm install && npm run build

# 构建后端（自动嵌入前端静态文件）
cd ..
go build -o quickce-server ./cmd/server

# 运行
./quickce-server
```

### Docker 部署

```bash
# 使用 docker-compose（推荐）
docker compose up -d

# 或手动构建运行
docker build -t quickce .
docker run -d --name quickce -p 8080:8080 \
  -e MAX_CONCURRENT=100 \
  quickce
```

**Dockerfile** 使用多阶段构建：
1. Node 镜像编译前端
2. Go 镜像编译后端
3. Alpine 镜像作为运行环境（仅含二进制和前端静态文件，镜像约 25MB）

## 配置

通过环境变量配置：

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `PORT` | `8080` | 监听端口 |
| `MAX_CONCURRENT` | `50` | 最大并发探测数 |
| `DEFAULT_TIMEOUT` | `5` | 默认超时秒数 |
| `STATIC_DIR` | `frontend/dist` | 前端静态文件目录 |
| `GIN_MODE` | `release` | Gin 运行模式 |

## API 文档

### 发起探测

```
POST /api/probe
Content-Type: application/json
```

**请求体：**

```json
{
  "target": "example.com",
  "type": "http",
  "port": 443,
  "timeout": 5
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `target` | string | 是 | 目标域名或 IP |
| `type` | string | 是 | 探测类型：`http`、`ping`、`tcp`、`dns`、`traceroute`、`ip` |
| `port` | number | 仅 tcp | 目标端口 |
| `timeout` | number | 否 | 超时秒数，默认 5 |
| `nodes` | string[] | 否 | 指定 Agent 名称/ID 列表，启用分布式拨测 |

**分布式探测请求体：**
```json
{
  "target": "example.com",
  "type": "http",
  "port": 443,
  "timeout": 5,
  "nodes": ["beijing-01", "shanghai-01"]
}
```

- 指定 `nodes` 后，Server 将任务分发给对应 Agent
- 响应为 `202 Accepted`，返回各 Agent 的 `task_id`

### 响应示例

**HTTP 探测：**
```json
{
  "success": true,
  "latency_ms": 234,
  "detail": { "status_code": 200 }
}
```

**Ping 探测：**
```json
{
  "success": true,
  "latency_ms": 42,
  "detail": { "rtt_ms": 42 }
}
```

**TCP 探测：**
```json
{
  "success": true,
  "latency_ms": 85,
  "detail": { "port_open": true, "rtt_ms": 85 }
}
```

**DNS 探测：**
```json
{
  "success": true,
  "latency_ms": 23,
  "detail": { "ips": ["93.184.216.34"] }
}
```

**IP 查询：**
```json
{
  "success": true,
  "latency_ms": 428,
  "detail": {
    "ip": "1.1.1.1",
    "isp": "Cloudflare",
    "organization": "Cloudflare",
    "asn": 13335,
    "asn_organization": "Cloudflare, Inc."
  }
}
```

**路由追踪：**
```json
{
  "success": true,
  "latency_ms": 1203,
  "detail": {
    "hops": [
      { "ttl": 1, "address": "192.168.1.1", "rtt": 1200000 },
      { "ttl": 2, "address": "104.20.23.154", "rtt": 15000000 }
    ]
  }
}
```

### curl 示例

```bash
# HTTP 探测
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"http","timeout":5}'

# Ping 探测
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"ping","timeout":5}'

# TCP 端口探测
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"tcp","port":443,"timeout":5}'

# DNS 解析探测
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"dns","timeout":5}'

# 路由追踪探测
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"traceroute","timeout":30}'

# IP 查询
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"1.1.1.1","type":"ip","timeout":10}'

# 分布式拨测（指定 agent 节点）
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"http","nodes":["beijing-01","shanghai-01"]}'
```

## 权限说明

- Ping 探测在 Linux 上调用系统 `ping` 命令，需要 `CAP_NET_RAW` 能力或 root 权限。
- 路由追踪在 Linux 上调用系统 `traceroute` 命令，可能需要 root 权限。
- Windows 下无额外权限要求。

## 分布式架构

QuickCE 支持单机运行和分布式多节点拨测两种模式。通过 Agent 客户端，可在多个地域部署探测节点，由中央服务器统一调度。

### 架构图

```
                        ┌──────────────┐
  POST /api/probe       │              │       GET /api/agents
  (含 nodes 参数) ───────▶    Server     ◀─────── 注册/轮询
                        │  (调度器)     │
                        └──────┬───────┘
                               │ 分配任务
                    ┌──────────┼──────────┐
                    │          │          │
                    ▼          ▼          ▼
              ┌──────────┐ ┌──────────┐ ┌──────────┐
              │ Agent A  │ │ Agent B  │ │ Agent C  │
              │ Beijing  │ │ Shanghai │ │ Shenzhen │
              └──────────┘ └──────────┘ └──────────┘
```

### 工作流程

1. **Agent 注册**：每个 Agent 启动时向 Server 注册（名称 + 地域）
2. **任务分发**：`POST /api/probe` 携带 `nodes` 参数指定目标 Agent
3. **轮询获取**：Agent 定期 `GET /api/agents/{id}/tasks` 拉取任务
4. **本地执行**：Agent 收到任务后调用本地探测函数执行
5. **结果上报**：Agent `POST /api/agents/{id}/results` 上报结果

### 启动 Agent

```bash
# 编译
go build -o quickce-agent ./cmd/agent

# 运行（注册到本地 Server）
./quickce-agent -name=beijing-01 -location=Beijing

# 指定 Server 地址和 Agent ID（重启后复用同一 ID）
./quickce-agent -server=http://server-ip:8080 -name=beijing-01 \
  -location=Beijing -agent-id=uuid-string
```

### 分布式拨测示例

```bash
# 查看已注册的 Agent
curl -s http://localhost:8080/api/agents

# 向指定 Agent 分发探测任务
curl -s -X POST http://localhost:8080/api/probe \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com","type":"http","timeout":5,"nodes":["beijing-01"]}'

# 响应：202 Accepted，返回 task_id
# Agent 执行完成后结果自动上报到 Server
```

### 关键 API

| 方法 | 路由 | 说明 |
|------|------|------|
| POST | `/api/agents/register` | Agent 注册 |
| GET | `/api/agents` | 列出所有 Agent |
| GET | `/api/agents/:id/tasks` | Agent 轮询任务（无任务返回 204） |
| POST | `/api/agents/:id/results` | Agent 上报结果 |

### 扩展设计

- `pkg/store/store.go` 定义了存储接口，当前使用内存实现，可替换为 MySQL / PostgreSQL
- Agent 支持 `-agent-id` 参数和 `AGENT_ID` 环境变量，重启后可保持同一身份
- Server 端 `nodes` 参数支持按 Agent 名称或 ID 筛选

## 设计要点

- **并发控制**：使用 `semaphore.Weighted` 限制全局最大并发探测数，防止资源耗尽。
- **超时控制**：每个探测任务通过 `context.WithTimeout` 独立超时，互不影响。
- **HTTP 连接池**：全局复用 `http.Transport`，避免每次探测新建连接，减少开销。
- **可扩展性**：`pkg/probe` 下的探测函数均遵循统一签名，新增探测类型只需添加新文件并在 `handler.go` 注册即可。
- **前后端一体部署**：后端自动检测 `frontend/dist` 目录并提供静态文件服务，支持 SPA 路由回退。
