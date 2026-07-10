# Manager UI Contract

**Platform**: Electron desktop (Windows/macOS/Linux) | **UI**: Shadcn + React + Tailwind

## Navigation Structure

| Route | Component | Spec reference |
|-------|-----------|----------------|
| `/monitoring` | MonitoringDashboard | FR-007–FR-009 |
| `/servers` | ServerList | FR-010 |
| `/servers/:id/docker` | DockerManager | FR-011–FR-012 |
| `/servers/:id/network` | NetworkManager | FR-013–FR-015 |
| `/servers/:id/services` | ServiceManager | FR-016–FR-017 |
| `/servers/:id/apps` | AppManager | FR-018–FR-019 |
| `/servers/:id/terminal` | TerminalView | FR-020 |
| `/servers/:id/events` | EventFeed | FR-054–FR-055 |
| `/servers/:id/logs` | LogConfigPanel | FR-041–FR-043 |
| `/servers/:id/actions` | ServerActions | FR-044–FR-046 |
| `/hubs` | HubList | FR-021 |
| `/agents` | AgentStatus | FR-022 |
| `/graph` | GraphNetwork | FR-023–FR-026 |
| `/commands` | CustomCommands | FR-044 |
| `/triggers` | AlertTriggers | FR-047–FR-050 |
| `/deploy` | RemoteDeploy | FR-051–FR-053 |
| `/settings` | Settings | FR-027, SMTP |
| `/about` | About | FR-028 |

---

## Monitoring Dashboard

### Modes

| Mode | View | Data source |
|------|------|-------------|
| `all` | Fleet aggregate cards + summary table | Latest `MetricBatch` per server |
| `per-server` | Server selector + Recharts line charts (24h) | `metric_history` table |

### Metrics displayed

CPU %, GPU % (or "N/A"), RAM %, Storage %, Network in/out, Uptime.

---

## Graph Network

- **Library**: React Flow
- **Node types**: `server`, `external_ip`
- **Edge label**: ping latency (ms), refreshed every 5s (SC-005)
- **Interactions**: drag reposition (persisted to `graph_nodes`), add external IP dialog
- **Visual assets**: icon/badge overlay from `visual_assets`

---

## Alert Trigger Editor

**Input format**:

```text
IF [CPU|GPU|RAM|STORAGE] [>|<] <threshold> THEN [Email|Reboot] <target>
```

**UI**: Structured form mapping to fields in `alert_triggers` table; raw text preview shown.

**Examples**:

```text
IF CPU > 90 THEN Email ops@example.com
IF RAM > 95 THEN Reboot server-web-01
```

---

## Custom Commands

| Field | Type | Required |
|-------|------|----------|
| `command_name` | text | yes |
| `command_script` | textarea (shell) | yes |
| `scope` | select: global / server | yes |
| `server_id` | select | when scope=server |

Execution opens output panel with streaming stdout/stderr.

---

## Server Actions (built-in)

| Action | Confirmation required | Progress indicator |
|--------|----------------------|-------------------|
| Reboot | yes | spinner + result |
| Set datetime | yes (shows picker) | spinner + result |
| Package update | no | streaming apt output |
| Package upgrade | yes | streaming apt output |

---

## Log Configuration Panel

Per server/hub target:

| Field | Default |
|-------|---------|
| Log root directory | `/var/log/nvx` |
| service.txt limit (MB) | 50 |
| errors.txt limit (MB) | 50 |
| connectivity.txt limit (MB) | 20 |

Displays current file sizes fetched via `log.config` command response.

---

## Remote Deploy Wizard

**Steps**:

1. Enter target host, SSH port, username, auth method (key/password)
2. Select component: Hub or Agent
3. Set initial name and description
4. Execute install with progress bar
5. Verify connection handshake

**Failure display**: error code + message; no partial registration (FR-053).

---

## Assets Folder

**Path**: `~/.config/nevarix/assets/` (platform-specific — see `data-model.md`)

**Supported formats**: PNG, SVG, WebP (max 512 KB per file)

**Assignment**: icon and/or badge per server, hub, or agent entity.

---

## Settings

| Section | Fields |
|---------|--------|
| SMTP | host, port, TLS, username, password, from address |
| General | metric retention days (default 30), theme |
| Connections | trusted certificate management |

---

## Real-Time Updates

| Channel | Mechanism | Latency target |
|---------|-----------|----------------|
| Metrics | gRPC stream → Electron main → IPC → renderer | ≤ 5s (SC-001) |
| Events | gRPC stream → event feed | ≤ 5s (SC-014) |
| Graph ping | polling interval 5s | SC-005 |
| Terminal | bidirectional IPC stream | interactive |

---

## Shadcn Components (required)

All screens MUST use Shadcn primitives: `Button`, `Table`, `Dialog`, `Form`, `Input`, `Select`, `Tabs`, `Card`, `Badge`, `Toast`, `Sheet`, `DropdownMenu`, `Chart` (Recharts wrapper).

No raw unstyled HTML form controls in production views.
