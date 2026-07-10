# Agent CLI Contract (`nvx-agent`)

**Binary**: `nvx-agent` | **Platform**: Ubuntu 22.04+ | **Framework**: Cobra

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Config file path (default `/etc/nvx/agent.yaml`) |
| `--json` | JSON output format |

---

## Commands

### `nvx-agent status`

**Output**:

```text
Agent: agent-1
Version: 1.0.0
Status: running
Hubs connected: 2
Last metric sent: 2026-07-10T08:30:00Z
Cache pending: 0 entries
```

---

### `nvx-agent service <disable|enable|restart>`

Controls systemd unit `nvx-agent.service`.

| Subcommand | Action |
|------------|--------|
| `disable` | `systemctl disable --now nvx-agent` |
| `enable` | `systemctl enable --now nvx-agent` |
| `restart` | `systemctl restart nvx-agent` |

---

### `nvx-agent log set dir --dir=<path>`

Sets `log_root_dir` in `/etc/nvx/agent.yaml`.

**Example**:

```bash
nvx-agent log set dir --dir=/var/log/nvx
```

---

### Hub Connection Commands

| Command | Description |
|---------|-------------|
| `nvx-agent hub search` | Discover hubs (mDNS + configured endpoints) |
| `nvx-agent hub list` | List known/connected hubs |
| `nvx-agent hub connection request --hub=<HUB-NAME>` | Send connection request |
| `nvx-agent hub connection disconnect --hub=<HUB-NAME>` | Disconnect from hub |
| `nvx-agent hub connection test --hub=<HUB-NAME>` | Test gRPC connectivity + latency |
| `nvx-agent hub last-sync --hub=<HUB-NAME>` | Last sync timestamp |

**`hub list` output columns**: `hub_name`, `connectivity_status`, `last_sync_at`

---

## Configuration File (`/etc/nvx/agent.yaml`)

```yaml
agent_name: "agent-1"
server_alias: "web-server-01"
grpc_listen: ":9444"
log_root_dir: "/var/log/nvx"
cache_dir: "/var/cache/nvx/outbound"
cache_max_mb: 512
metrics_interval_sec: 15
hub_endpoints:
  - "hub-2.internal:9443"
mdns_enabled: true
tls:
  cert_file: "/etc/nvx/certs/agent.crt"
  key_file: "/etc/nvx/certs/agent.key"
  ca_file: "/etc/nvx/certs/ca.crt"
```

---

## Log File Layout

Identical to Hub — see `hub-cli.md`.

Default `log_root_dir`: `/var/log/nvx`

---

## Metrics Collected

| Metric | Source |
|--------|--------|
| CPU % | gopsutil |
| GPU % | nvidia-smi (optional) |
| RAM % | gopsutil |
| Storage % | gopsutil disk usage |
| Network in/out | gopsutil net IO counters |
| Uptime | gopsutil / proc uptime |

Collection interval: 15 seconds (configurable via `metrics_interval_sec`).

---

## Events Emitted

| Event type | Trigger |
|------------|---------|
| `connectivity.lost` | Hub connection dropped |
| `connectivity.restored` | Hub connection re-established |
| `service.state_changed` | systemd unit state change |
| `action.completed` | Remote command finished |
| `action.failed` | Remote command error |
| `system.reboot` | Server reboot initiated |

See `wire-protocol.md` → `ServerEvent` for wire format.
