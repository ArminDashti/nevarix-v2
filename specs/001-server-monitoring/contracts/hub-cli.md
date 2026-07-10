# Hub CLI Contract (`nvx-hub`)

**Binary**: `nvx-hub` | **Platform**: Ubuntu 22.04+ | **Framework**: Cobra

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Config file path (default `/etc/nvx/hub.yaml`) |
| `--json` | JSON output format |

---

## Commands

### `nvx-hub status`

**Output** (human-readable):

```text
Hub: hub-2
Version: 1.0.0
Status: running
Managers connected: 1
Agents connected: 3
Last sync: 2026-07-10T08:30:00Z
```

**Exit codes**: `0` running, `1` stopped/error.

---

### `nvx-hub service <disable|enable|restart>`

Controls systemd unit `nvx-hub.service`.

| Subcommand | Action |
|------------|--------|
| `disable` | `systemctl disable --now nvx-hub` |
| `enable` | `systemctl enable --now nvx-hub` |
| `restart` | `systemctl restart nvx-hub` |

---

### `nvx-hub log set dir --dir=<path>`

Sets `log_root_dir` in `/etc/nvx/hub.yaml`. Creates directory if missing.

**Example**:

```bash
nvx-hub log set dir --dir=/var/log/nvx
```

---

### Agent Connection Commands

| Command | Description |
|---------|-------------|
| `nvx-hub agent search` | Discover agents (mDNS + configured endpoints) |
| `nvx-hub agent list` | List known/connected agents |
| `nvx-hub agent connection request --agent=<AGENT-NAME>` | Send connection request |
| `nvx-hub agent connection disconnect --agent=<AGENT-NAME>` | Disconnect agent |
| `nvx-hub agent connection --agent=<AGENT-NAME>` | Show connection status |
| `nvx-hub agent last-sync --agent=<AGENT-NAME>` | Last sync timestamp |

**`agent list` output columns**: `agent_name`, `connectivity_status`, `last_sync_at`

---

### Manager Connection Commands

| Command | Description |
|---------|-------------|
| `nvx-hub manager search` | Discover managers on network |
| `nvx-hub manager list` | List known/connected managers |
| `nvx-hub manager connection request --manager=<MANAGER-NAME>` | Send connection request |
| `nvx-hub manager connection disconnect --manager=<MANAGER-NAME>` | Disconnect manager |
| `nvx-hub manager connection test --manager=<MANAGER-NAME>` | Test gRPC connectivity + latency |
| `nvx-hub manager last-sync --manager=<MANAGER-NAME>` | Last sync timestamp |

---

## Configuration File (`/etc/nvx/hub.yaml`)

```yaml
hub_name: "hub-2"
grpc_listen: ":9443"
log_root_dir: "/var/log/nvx"
db_path: "/var/lib/nvx/hub.db"
cache_dir: "/var/cache/nvx/outbound"
mdns_enabled: true
tls:
  cert_file: "/etc/nvx/certs/hub.crt"
  key_file: "/etc/nvx/certs/hub.key"
  ca_file: "/etc/nvx/certs/ca.crt"
```

---

## Log File Layout

```text
<log_root_dir>/<YYYY>/<MM>/<DD>/service.txt
<log_root_dir>/<YYYY>/<MM>/<DD>/errors.txt
<log_root_dir>/<YYYY>/<MM>/<DD>/connectivity.txt
```

Default `log_root_dir`: `/var/log/nvx`
