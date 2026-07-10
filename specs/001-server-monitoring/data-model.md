# Data Model: Server Monitoring Platform

**Feature**: `001-server-monitoring` | **Date**: 2026-07-10

## Overview

Three persistence layers:

| Layer | Technology | Location | Owner |
|-------|------------|----------|-------|
| Manager store | SQLite | `~/.config/nevarix/manager.db` | Manager app |
| Hub store | SQLite | `/var/lib/nvx/hub.db` | Hub service |
| Outbound cache | File spool | `/var/cache/nvx/outbound/` | Hub & Agent |
| Logs | Text files | `/var/log/nvx/<YYYY>/<MM>/<DD>/` | Hub & Agent |
| Visual assets | Filesystem | `~/.config/nevarix/assets/` | Manager app |

---

## Manager Database (`manager.db`)

### `servers`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `server_id` | UUID | PK | Unique server identity |
| `server_alias` | TEXT | NOT NULL, UNIQUE | Display name |
| `agent_name` | TEXT | UNIQUE | Linked agent identifier |
| `hub_id` | UUID | FK → `hubs.hub_id` | Primary routing hub |
| `hardware_json` | JSON | | CPU, RAM, MB, storage, NICs, VM, OS snapshot |
| `online_status` | ENUM | `online`, `offline`, `degraded` | Current connectivity |
| `last_seen_at` | TIMESTAMP | | Last metric/event received |
| `icon_asset_id` | UUID | FK → `visual_assets.asset_id`, NULL | Custom icon |
| `badge_asset_id` | UUID | FK → `visual_assets.asset_id`, NULL | Custom badge |
| `created_at` | TIMESTAMP | NOT NULL | |
| `updated_at` | TIMESTAMP | NOT NULL | |

**Validation**: One `server_id` per physical machine (FR-005). `server_alias` max 128 chars.

---

### `hubs`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `hub_id` | UUID | PK | |
| `hub_name` | TEXT | NOT NULL, UNIQUE | |
| `server_id` | UUID | FK → `servers.server_id`, NULL | Host server if co-located |
| `description` | TEXT | | |
| `endpoint` | TEXT | NOT NULL | `host:port` |
| `connectivity_status` | ENUM | `connected`, `disconnected`, `pending` | |
| `last_sync_at` | TIMESTAMP | | |
| `created_at` | TIMESTAMP | NOT NULL | |

---

### `agents`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `agent_id` | UUID | PK | |
| `agent_name` | TEXT | NOT NULL, UNIQUE | |
| `server_id` | UUID | FK → `servers.server_id`, NOT NULL | |
| `description` | TEXT | | |
| `connectivity_status` | ENUM | `connected`, `disconnected`, `pending` | Via primary hub |
| `last_sync_at` | TIMESTAMP | | |
| `created_at` | TIMESTAMP | NOT NULL | |

---

### `log_configurations`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `config_id` | UUID | PK | |
| `target_type` | ENUM | `server`, `hub` | |
| `target_id` | UUID | NOT NULL | FK to server or hub |
| `log_root_dir` | TEXT | DEFAULT `/var/log/nvx` | |
| `service_txt_limit_mb` | INTEGER | DEFAULT 50 | |
| `errors_txt_limit_mb` | INTEGER | DEFAULT 50 | |
| `connectivity_txt_limit_mb` | INTEGER | DEFAULT 20 | |
| `updated_at` | TIMESTAMP | NOT NULL | |

---

### `custom_commands`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `command_id` | UUID | PK | |
| `command_name` | TEXT | NOT NULL | |
| `command_script` | TEXT | NOT NULL | Shell script body |
| `scope` | ENUM | `global`, `server` | |
| `server_id` | UUID | FK, NULL | Required when scope=server |
| `created_at` | TIMESTAMP | NOT NULL | |

**Validation**: `command_name` unique per scope+server.

---

### `alert_triggers`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `trigger_id` | UUID | PK | |
| `metric_type` | ENUM | `CPU`, `GPU`, `RAM`, `STORAGE` | |
| `operator` | ENUM | `>`, `<` | |
| `threshold` | REAL | NOT NULL | Percentage 0–100 |
| `action_type` | ENUM | `Email`, `Reboot` | |
| `action_target` | TEXT | NOT NULL | Email address or server_id |
| `server_id` | UUID | FK, NULL | NULL = fleet-wide |
| `enabled` | BOOLEAN | DEFAULT true | |
| `created_at` | TIMESTAMP | NOT NULL | |

**Validation**: GPU triggers skipped when `gpu_status = unavailable`.

---

### `graph_nodes`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `node_id` | UUID | PK | |
| `node_type` | ENUM | `server`, `external_ip` | |
| `reference_id` | UUID | NULL | server_id when type=server |
| `external_ip` | TEXT | NULL | When type=external_ip |
| `position_x` | REAL | NOT NULL | |
| `position_y` | REAL | NOT NULL | |
| `updated_at` | TIMESTAMP | NOT NULL | |

---

### `visual_assets`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `asset_id` | UUID | PK | |
| `asset_name` | TEXT | NOT NULL | |
| `asset_type` | ENUM | `icon`, `badge` | |
| `file_path` | TEXT | NOT NULL | Relative to assets folder |
| `created_at` | TIMESTAMP | NOT NULL | |

---

### `metric_history`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `metric_id` | BIGINT | PK AUTO | |
| `server_id` | UUID | FK, NOT NULL | |
| `recorded_at` | TIMESTAMP | NOT NULL | |
| `cpu_percent` | REAL | | |
| `gpu_percent` | REAL | NULL | NULL = unavailable |
| `ram_percent` | REAL | | |
| `storage_percent` | REAL | | |
| `network_bytes_in` | BIGINT | | |
| `network_bytes_out` | BIGINT | | |
| `uptime_seconds` | BIGINT | | |

**Retention**: 30 days default; indexed on `(server_id, recorded_at)`.

---

### `server_events`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `event_id` | UUID | PK | |
| `server_id` | UUID | FK, NOT NULL | |
| `event_type` | TEXT | NOT NULL | e.g. `service.restarted`, `connectivity.lost` |
| `event_payload` | JSON | | |
| `occurred_at` | TIMESTAMP | NOT NULL | |
| `received_at` | TIMESTAMP | NOT NULL | |

---

### `smtp_settings`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `setting_id` | INTEGER | PK (singleton row) | |
| `smtp_host` | TEXT | | |
| `smtp_port` | INTEGER | DEFAULT 587 | |
| `smtp_tls` | BOOLEAN | DEFAULT true | |
| `smtp_user` | TEXT | | |
| `smtp_password_enc` | TEXT | | Encrypted at rest |
| `from_address` | TEXT | | |

---

## Hub Database (`hub.db`)

### `managers`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `manager_id` | UUID | PK | |
| `manager_name` | TEXT | NOT NULL | |
| `certificate_fingerprint` | TEXT | NOT NULL | mTLS identity |
| `connectivity_status` | ENUM | `connected`, `disconnected`, `pending` | |
| `last_sync_at` | TIMESTAMP | | |

---

### `agents`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `agent_id` | UUID | PK | |
| `agent_name` | TEXT | NOT NULL, UNIQUE | |
| `certificate_fingerprint` | TEXT | NOT NULL | |
| `connectivity_status` | ENUM | `connected`, `disconnected`, `pending` | |
| `last_sync_at` | TIMESTAMP | | |
| `server_identity` | TEXT | NOT NULL | Canonical server fingerprint |

---

### `connection_requests`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `request_id` | UUID | PK | |
| `requester_type` | ENUM | `manager`, `agent` | |
| `requester_name` | TEXT | NOT NULL | |
| `target_type` | ENUM | `manager`, `agent` | |
| `target_name` | TEXT | NOT NULL | |
| `status` | ENUM | `pending`, `accepted`, `rejected`, `expired` | |
| `requested_at` | TIMESTAMP | NOT NULL | |
| `expires_at` | TIMESTAMP | NOT NULL | Default +24h |

**State transitions**: `pending` → `accepted` | `rejected` | `expired`

---

### `route_cache`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `agent_id` | UUID | PK, FK | |
| `manager_id` | UUID | FK | Primary manager route |
| `last_routed_at` | TIMESTAMP | | |

---

### `sync_metadata`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `peer_type` | ENUM | `manager`, `agent` | |
| `peer_id` | UUID | PK composite | |
| `last_sync_at` | TIMESTAMP | | |
| `sync_version` | INTEGER | DEFAULT 0 | |

---

### `outbound_queue`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `queue_id` | BIGINT | PK AUTO | |
| `direction` | ENUM | `to_manager`, `to_agent` | |
| `peer_id` | UUID | NOT NULL | |
| `payload_file` | TEXT | NOT NULL | Path under `/var/cache/nvx/outbound/` |
| `enqueued_at` | TIMESTAMP | NOT NULL | |
| `priority` | INTEGER | DEFAULT 0 | Lower = older FIFO |

---

## Agent Local State (minimal SQLite optional; primary cache is file spool)

Agent uses `/var/cache/nvx/outbound/` with identical envelope schema to Hub. Agent config at `/etc/nvx/agent.yaml`:

```yaml
agent_name: "agent-1"
log_root_dir: "/var/log/nvx"
hub_endpoints: []
cache_max_mb: 512
metrics_interval_sec: 15
```

---

## Protobuf Message Types (wire entities)

See `contracts/wire-protocol.md` for full schema. Key messages:

| Message | Direction | Purpose |
|---------|-----------|---------|
| `MetricBatch` | Agent → Hub → Manager | Periodic metrics |
| `ServerEvent` | Agent → Hub → Manager | Real-time events |
| `CommandRequest` | Manager → Hub → Agent | Actions, custom commands |
| `CommandResponse` | Agent → Hub → Manager | Action results |
| `TerminalFrame` | Bidirectional | SSH I/O chunks |
| `ConnectionRequest` | Any peer | Request-and-accept handshake |
| `LogConfigUpdate` | Manager → Hub/Agent | Log dir and limits |
| `InventorySnapshot` | Agent → Hub → Manager | Hardware/OS info |

---

## Entity Relationship Summary

```text
Manager ──(1:N)── Hub ──(N:M)── Agent ──(N:1)── Server

Server 1:1 Agent (max one agent per server)
Server 1:1 Hub  (max one hub per server)
Server may have both Agent and Hub
Agent N:M Hub   (multi-hub membership)
Manager N:1 Hub (primary route per agent via hub_id on server)
```

---

## Validation Rules Summary

| Rule | Source |
|------|--------|
| Unique `server_alias` per fleet | FR-005, FR-010 |
| Unique `agent_name`, `hub_name` globally | FR-006b |
| Log paths must match `<DIR>/<YYYY>/<MM>/<DD>/` pattern | FR-056, FR-061 |
| Connectivity log format pipe-delimited with ISO8601 timestamp | FR-057 |
| Trigger threshold 0–100 for percentage metrics | FR-047 |
| GPU triggers inactive when GPU unavailable | Edge case |
| Connection requests expire after 24h if pending | Edge case |
