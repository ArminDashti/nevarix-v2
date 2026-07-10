# Research: Server Monitoring Platform

**Feature**: `001-server-monitoring` | **Date**: 2026-07-10

## 1. Inter-Component Communication Protocol

**Decision**: gRPC over mutual TLS (mTLS) for Hub↔Agent and Manager↔Hub data plane; protobuf message schemas in a shared `packages/proto` module.

**Rationale**:
- Bidirectional streaming supports real-time metrics, events, and terminal I/O relay through Hub.
- Protobuf gives versioned, typed contracts across Go (Hub/Agent) and TypeScript (Manager via `@grpc/grpc-js`).
- mTLS satisfies the request-and-accept trust model with certificate pinning per connection pair.

**Alternatives considered**:
- *REST polling*: Simple but poor fit for terminal streaming, event push, and sub-5s metric latency (SC-001, SC-014).
- *Raw WebSocket + JSON*: Flexible but lacks schema enforcement and increases cross-language drift risk.
- *MQTT*: Good for IoT pub/sub but adds broker dependency not required by the three-tier topology.

---

## 2. Manager Desktop Stack

**Decision**: Electron 33+ with React 19, TypeScript 5.x, Vite, Shadcn UI (Radix + Tailwind), TanStack Query for server state, Recharts for monitoring charts, React Flow for graph network.

**Rationale**:
- Electron provides cross-platform desktop shell with native SSH/process capabilities in the main process.
- Shadcn is an explicit spec requirement (FR-028a) and integrates cleanly with React + Tailwind.
- React Flow supports draggable nodes with persisted layout (FR-025).

**Alternatives considered**:
- *Tauri*: Lighter footprint but weaker ecosystem for xterm.js + gRPC native addons on all platforms.
- *WPF/.NET*: Conflicts with clarified TypeScript/Electron direction from original specify input.

---

## 3. Hub and Agent Runtime (Go)

**Decision**: Go 1.23+ services using Cobra for CLI (`nvx-hub`, `nvx-agent`), `google.golang.org/grpc` for RPC, `modernc.org/sqlite` for Hub persistence, `shirou/gopsutil/v4` for metrics collection on Agent, systemd units for service lifecycle on Ubuntu.

**Rationale**:
- Single static binary per component simplifies remote installation (FR-051, FR-052).
- gopsutil is the de facto cross-platform metrics library; Ubuntu-only scope allows direct `/proc` and `systemd` integration.
- modernc.org/sqlite is pure Go (no CGO) — simpler cross-compilation for amd64/arm64 Ubuntu targets.

**Alternatives considered**:
- *CGO sqlite3*: Better performance but complicates ARM builds.
- *PostgreSQL on Hub*: Over-engineered for per-server Hub instance with local operational state.

---

## 4. Hub SQLite Schema Strategy

**Decision**: Hub stores operational data in `/var/lib/nvx/hub.db` (SQLite) with tables: `managers`, `agents`, `connections`, `pending_requests`, `route_cache`, `sync_metadata`, `outbound_queue`.

**Rationale**:
- FR-059 requires local persistence for connections, routing, and sync state.
- SQLite fits single-process Hub service with moderate write volume.
- Outbound queue metadata complements file-based payload cache (see §5).

**Alternatives considered**:
- *BoltDB*: Key-value only; relational queries for connection lists and sync status are simpler in SQL.
- *JSON files*: No transactional guarantees for concurrent connection state updates.

---

## 5. Outbound Cache (Hub & Agent)

**Decision**: File-spool cache at `/var/cache/nvx/outbound/` with JSON envelope files named `<timestamp>_<uuid>.json`; default 24-hour capacity (~10,000 entries or 512 MB, whichever first); FIFO eviction on overflow.

**Rationale**:
- Spec requires chronological flush on reconnection (FR-060, FR-064) and 24-hour disconnected tolerance.
- File spool survives process restarts without DB corruption risk for large payloads (terminal chunks, metric batches).
- SQLite `outbound_queue` table tracks index/order; payload body in files.

**Alternatives considered**:
- *In-memory only*: Lost on restart — unacceptable for extended outages.
- *Redis*: External dependency violates per-server self-contained Agent/Hub design.

---

## 6. Logging Implementation

**Decision**: Structured logging via Go `slog` (Hub/Agent) writing to the spec-defined paths under `/var/log/nvx/<YYYY>/<MM>/<DD>/`; connectivity ticker goroutine every 60s; 5-day retention via daily cleanup job; size limits enforced by logrotate-style truncation per file type.

**Rationale**:
- Matches clarified log structure exactly (FR-056–FR-063).
- `slog` is stdlib in Go 1.21+ with minimal dependencies.
- Separate files (`service.txt`, `errors.txt`, `connectivity.txt`) simplify Manager-side log limit configuration.

**Alternatives considered**:
- *journald only*: Does not match required file layout and Manager-configurable paths.
- *Single combined log*: Violates per-file size limit requirement (FR-042).

---

## 7. Connection Discovery & Request-and-Accept

**Decision**: Manual registration via CLI `search`/`list`/`connection request` using configured host:port endpoints; optional LAN discovery via mDNS service types `_nvx-hub._tcp` and `_nvx-agent._tcp` (opt-in via config flag).

**Rationale**:
- Spec defines explicit CLI connection workflow (FR-032, FR-033, FR-037).
- mDNS as optional accelerator for `search` commands in local networks; production fleets typically use static addresses.

**Alternatives considered**:
- *Central registry service*: Adds fourth component — rejected per three-tier architecture.
- *Automatic trust on discovery*: Violates request-and-accept security model.

---

## 8. Authentication & Security

**Decision**: mTLS with component-specific X.509 certificates; connection requests carry signed tokens exchanged during accept; Manager stores trusted Hub certificates; Hub stores trusted Agent/Manager certificates; SSH key-based auth for remote install and terminal (OpenSSH via `golang.org/x/crypto/ssh`).

**Rationale**:
- Resolves spec assumption "authentication mechanisms defined during planning."
- Request-and-accept maps to pending request records until operator approves via CLI or Manager UI.
- SSH aligns with Ubuntu remote install requirement (FR-051, FR-052).

**Alternatives considered**:
- *Shared API keys only*: No identity binding between components; insufficient for multi-hub topology.
- *OAuth2*: Overkill for server-to-server component mesh without human SSO requirement.

---

## 9. Manager Local Persistence

**Decision**: Manager uses SQLite at `%APPDATA%/Nevarix/manager.db` (Windows) / `~/.config/nevarix/manager.db` (Linux/macOS) for servers, hubs, agents, triggers, custom commands, graph layout, visual asset references, and SMTP settings; assets folder at `%APPDATA%/Nevarix/assets/`.

**Rationale**:
- Manager must persist triggers, custom commands, graph positions, and asset assignments across sessions.
- SQLite avoids external DB dependency for desktop app.
- Assets folder satisfies FR-028b with filesystem storage for icons/badges.

**Alternatives considered**:
- *JSON config files*: Poor query support for 100-server fleet (SC-008).
- *Electron localStorage*: Size and structure limits for metrics history cache.

---

## 10. Metrics & GPU Collection (Agent)

**Decision**: Agent collects metrics every 15s via gopsutil; GPU via `nvidia-smi` JSON output when available, otherwise reports `gpu_status: unavailable`; batches sent to Hub via gRPC stream.

**Rationale**:
- 15s interval balances SC-001 (5s view load uses cached/latest) with network overhead for 100 servers.
- nvidia-smi is standard on Ubuntu GPU servers; avoids CGO NVML binding complexity in v1.

**Alternatives considered**:
- *5s collection*: Higher load on Hub relay for 100-agent fleet.
- *DCGM exporter sidecar*: Extra deployment complexity per server.

---

## 11. Alert Trigger Evaluation

**Decision**: Manager evaluates triggers against incoming metric stream in Electron main process; Email via Nodemailer + configurable SMTP (Settings); Reboot dispatched as gRPC command Manager→Hub→Agent.

**Rationale**:
- Triggers defined in Manager UI (FR-047); evaluation co-located with SMTP config.
- Reboot requires Agent execution path — consistent with other server actions.

**Alternatives considered**:
- *Agent-side evaluation*: Duplicates trigger config across agents; harder to manage centrally.
- *Hub-side evaluation*: Hub lacks SMTP access and operator-defined trigger store.

---

## 12. Remote Terminal

**Decision**: xterm.js in Manager renderer; SSH session initiated from Electron main process via `ssh2` (Node) or Go ssh helper spawned by main; stdin/stdout bridged through gRPC bidirectional stream Agent→Hub→Manager.

**Rationale**:
- FR-020 requires interactive terminal; SSH is standard on Ubuntu servers.
- Bridging through existing gRPC path avoids opening additional firewall ports beyond Hub.

**Alternatives considered**:
- *Direct Manager→Agent SSH*: Bypasses Hub topology — violates FR-001.
- *WebSocket-only shell*: Would require custom shell daemon on Agent — redundant with SSH.

---

## 13. Remote Hub/Agent Installation

**Decision**: Manager main process uses SSH + SFTP to upload pre-built `.deb` packages (or tarball fallback), runs `dpkg -i` / `apt install -f`, registers systemd units (`nvx-hub.service`, `nvx-agent.service`), and verifies health via gRPC handshake.

**Rationale**:
- Ubuntu-only scope (FR-006a) makes `.deb` the natural packaging format.
- systemd matches `service enable/disable/restart` CLI commands (FR-030, FR-035).

**Alternatives considered**:
- *curl | bash install script*: Less idempotent; harder to report structured progress (FR-053).
- *Ansible dependency*: External tool requirement on operator workstation.

---

## 14. Testing Strategy

**Decision**:
- **Go (Hub/Agent)**: `testing` + `testify`; integration tests with testcontainers or docker-compose fixture; contract tests against protobuf schemas.
- **Manager**: Vitest for unit tests; Playwright for Electron E2E smoke tests.
- **Cross-component**: docker-compose stack (1 Manager dev build + 2 Hubs + 3 Agents) for quickstart validation.

**Rationale**:
- Matches constitution-template emphasis on contract and integration testing (applied as default gates — no project constitution ratified yet).

**Alternatives considered**:
- *Manual-only QA*: Insufficient for 65 functional requirements and streaming protocol.

---

## Resolved NEEDS CLARIFICATION Items

| Item | Resolution |
|------|------------|
| Authentication mechanism | mTLS + signed connection tokens |
| Hub SQLite schema scope | 7 tables — see `data-model.md` |
| SMTP configuration | Manager Settings → Nodemailer; fields: host, port, TLS, user, password, from_address |
| Cache capacity | 512 MB / 10,000 entries FIFO |
| Metric collection interval | 15 seconds |
| Protocol between components | gRPC + protobuf + mTLS |
