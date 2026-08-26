# 服务注册发现中心（service-registry）

一个基于纯 Go 标准库（`net/http`，零第三方依赖）的服务注册与发现后端，提供服务/实例注册、健康检查、心跳、负载均衡选择、服务依赖拓扑、事件审计与多维统计。内置鉴权、限流、CORS、请求 ID、日志与 panic 恢复等中间件，并提供原生前端控制台。

## 运行

```bash
cd origin
go run ./cmd/server
# 默认监听 :8080，可用 PORT / ADDR 覆盖
# 可选环境变量：AUTH_TOKEN（启用 Bearer 鉴权）、RATE_LIMIT（每 IP 每分钟限流，默认 200）
# HEARTBEAT_TIMEOUT_SEC（默认 30）、HEALTH_CHECK_INTERVAL_SEC（默认 10）
```

访问前端控制台：`http://localhost:8080/`（需在 `origin/` 目录下启动）。

## 分层结构

```
origin/
├── cmd/server/main.go        # 入口 + 前端挂载 + 优雅关闭
├── frontend/index.html       # 原生前端控制台（零构建）
├── internal/
│   ├── app/ config/ model/ store/ service/ handler/
└── pkg/ httpx/ idgen/ logger/ semver/ strutil/ backoff/ retry/
```

## API 一览（核心）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET | /api/services | 注册 / 列表（keyword/status/owner） |
| GET | /api/services/{id} | 服务详情 |
| PUT | /api/services/{id} | 更新服务 |
| PATCH | /api/services/{id}/status | 上下线 |
| DELETE | /api/services/{id} | 注销服务（级联） |
| GET | /api/services/{id}/versions | 版本列表（semver 降序） |
| GET | /api/services/{id}/latest | 最新版本 |
| GET | /api/services/{id}/stats | 服务统计 |
| GET | /api/services/{id}/select | 负载均衡选一个实例 |
| POST/GET | /api/instances | 注册 / 列表 |
| GET | /api/instances/{id} | 实例详情 |
| PATCH | /api/instances/{id}/status | 状态机流转 |
| DELETE | /api/instances/{id} | 注销实例 |
| GET | /api/instances/{id}/health | 实例健康视图 |
| POST/GET | /api/health-checks | 创建 / 列表 |
| POST | /api/health-checks/{id}/report | 上报检查结果 |
| POST/GET | /api/heartbeats | 心跳 / 查询 |
| POST | /api/maintenance/run | 巡检超时实例 |
| POST/GET | /api/load-balances | 负载均衡策略 |
| POST/GET | /api/dependencies | 依赖声明 / 列表 |
| POST | /api/dependencies/{id}/call | 调用计数 +1 |
| GET | /api/events | 事件流（分页） |
| GET | /api/stats/registry | 整体统计 |
| GET | /api/stats/top-dependencies | 依赖 Top N |
| GET | /api/export · POST /api/import | 快照导出 / 导入 |
| POST/GET | /api/labels | 服务标签 |
| GET/PUT/DELETE | /api/configs[/{key}] | 注册中心配置 |
| GET | /api/search?q= | 全局搜索 |
| GET | /api/topology · /api/topology/cycles | 依赖拓扑 / 环检测 |
| GET | /api/reports/{regions,algorithms,events} | 分布报告 |
| GET | /api/metrics/instances | 实例延迟指标 |
| GET | /api/health-matrix | 服务健康矩阵 |
| GET | /api/probe · /api/probe/{id} | 实例健康探测 |
| POST | /api/reconcile | 一致性巡检 |
| GET | /api/closure/{id} | 依赖传递闭包 |
| GET | /api/versions/{name} · /api/versions/outdated | 版本汇总 / 过期服务 |
| GET | /api/health-score | 整体健康分 |
| GET | /healthz · /readyz | 存活 / 就绪探针 |

## 关键设计

- **状态机**：实例 `starting → up → down → up`；健康检查 `healthy ↔ unhealthy`，联动实例上下线并记录事件。
- **负载均衡**：`round_robin` / `random` / `weighted` / `least_conn` 四种算法，仅选择 `up` 状态实例。
- **心跳超时**：`LastHeartbeat` 超过 `HEARTBEAT_TIMEOUT_SEC` 的 up 实例会被巡检置为 down。
- **依赖拓扑**：支持上游/下游查询、最短路径（BFS）、传递闭包与环检测。
- **版本管理**：`pkg/semver` 做语义化版本比较，支持最新版本与过期服务识别。

## 统一响应

```json
{ "code": 0, "message": "ok", "data": { } }
```
