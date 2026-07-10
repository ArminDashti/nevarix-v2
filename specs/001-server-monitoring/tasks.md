---
description: "Task list for Server Monitoring Platform implementation"
---

# Tasks: Server Monitoring Platform

**Input**: Design documents from `/specs/001-server-monitoring/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Not requested in spec — test tasks omitted.

**Organization**: Tasks grouped by user story for independent implementation and validation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story label (US1–US16)

## Path Conventions

Monorepo layout per plan.md: `apps/manager/`, `apps/hub/`, `apps/agent/`, `packages/proto/`, `deploy/`, `configs/`, `scripts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize monorepo structure, tooling, and deployment scaffolding

- [x] T001 Create monorepo directory structure per plan.md (`apps/manager/`, `apps/hub/`, `apps/agent/`, `packages/proto/`, `deploy/`, `configs/`, `scripts/dev/`, `scripts/build/`)
- [x] T002 Initialize Go module for Hub in `apps/hub/go.mod` with module path and Go 1.23 toolchain
- [x] T003 [P] Initialize Go module for Agent in `apps/agent/go.mod` with module path and Go 1.23 toolchain
- [x] T004 [P] Initialize Manager Electron project in `apps/manager/package.json` with Electron 33, React 19, Vite, TypeScript 5
- [x] T005 [P] Configure Shadcn UI in `apps/manager/components.json` and install base components in `apps/manager/src/renderer/components/ui/`
- [x] T006 [P] Initialize protobuf workspace in `packages/proto/buf.yaml` and `packages/proto/nvx/v1/`
- [x] T007 Create root `Makefile` with targets: `proto`, `build-hub`, `build-agent`, `build-manager`, `test`, `compose-up`
- [x] T008 [P] Add Hub config template in `configs/hub.yaml` with `hub_name`, `grpc_listen`, `log_root_dir`, `db_path`, `cache_dir`, TLS paths
- [x] T009 [P] Add Agent config template in `configs/agent.yaml` with `agent_name`, `server_alias`, `metrics_interval_sec`, hub endpoints, TLS paths
- [x] T010 [P] Add systemd unit files in `deploy/nvx-hub.service` and `deploy/nvx-agent.service`
- [x] T011 [P] Add dev TLS cert generation script in `scripts/dev/generate-certs.sh`
- [x] T012 [P] Add docker-compose dev fixture in `deploy/docker-compose.dev.yml`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: gRPC protocol, shared infrastructure, and minimal connectivity — **BLOCKS all user stories**

**⚠️ CRITICAL**: No user story work begins until this phase is complete

- [ ] T013 Define protobuf messages (`MetricBatch`, `ServerEvent`, `CommandRequest`, `CommandResponse`, `ConnectionRequest`, `InventorySnapshot`, `TerminalFrame`, `LogConfigUpdate`) in `packages/proto/nvx/v1/nvx.proto`
- [ ] T014 Define gRPC services (`NvxHubService`, `NvxAgentService`) in `packages/proto/nvx/v1/nvx.proto` per `contracts/wire-protocol.md`
- [ ] T015 Configure buf code generation for Go stubs in `packages/proto/buf.gen.yaml` and add `make proto` target output to `apps/hub/internal/pb/` and `apps/agent/internal/pb/`
- [ ] T016 [P] Configure TypeScript proto/gRPC stub generation to `apps/manager/src/shared/pb/`
- [ ] T017 [P] Implement mTLS loader and certificate validation in `apps/hub/internal/tls/config.go`
- [ ] T018 [P] Implement mTLS loader and certificate validation in `apps/agent/internal/tls/config.go`
- [ ] T019 [P] Implement mTLS client configuration in `apps/manager/src/main/tls/config.ts`
- [ ] T020 Implement daily log writer (`service.txt`, `errors.txt`, `connectivity.txt`) in `apps/hub/internal/logging/writer.go`
- [ ] T021 [P] Implement daily log writer in `apps/agent/internal/logging/writer.go`
- [ ] T022 Implement connectivity ticker (60s interval, pipe-delimited format) in `apps/hub/internal/logging/connectivity.go`
- [ ] T023 [P] Implement connectivity ticker in `apps/agent/internal/logging/connectivity.go`
- [ ] T024 Implement outbound file-spool cache in `apps/hub/internal/cache/spool.go`
- [ ] T025 [P] Implement outbound file-spool cache in `apps/agent/internal/cache/spool.go`
- [ ] T026 Create Hub SQLite schema migrations in `apps/hub/internal/store/migrations/001_initial.sql` (managers, agents, connection_requests, route_cache, sync_metadata, outbound_queue)
- [ ] T027 Create Manager SQLite schema migrations in `apps/manager/src/main/db/migrations/001_initial.sql` (servers, hubs, agents, metric_history, server_events)
- [ ] T028 Implement Hub SQLite repository in `apps/hub/internal/store/repository.go`
- [ ] T029 Implement Manager SQLite repository in `apps/manager/src/main/db/repository.ts`
- [ ] T030 Implement Hub gRPC server skeleton in `apps/hub/internal/grpc/server.go` with mTLS and stream handlers
- [ ] T031 Implement Agent gRPC client to Hub in `apps/agent/internal/grpc/hub_client.go`
- [ ] T032 Implement Hub-to-Agent gRPC relay in `apps/hub/internal/grpc/relay.go`
- [ ] T033 Implement Manager gRPC client to Hub in `apps/manager/src/main/grpc/hub_client.ts`
- [ ] T034 Implement connection request-and-accept flow in `apps/hub/internal/grpc/connections.go`
- [ ] T035 [P] Create Cobra CLI root and `status` command skeleton in `apps/hub/cmd/nvx-hub/main.go`
- [ ] T036 [P] Create Cobra CLI root and `status` command skeleton in `apps/agent/cmd/nvx-agent/main.go`
- [ ] T037 Setup Electron main process entry in `apps/manager/src/main/index.ts` with IPC bridge
- [ ] T038 Setup preload context bridge in `apps/manager/src/preload/index.ts`
- [ ] T039 Setup React renderer shell with Shadcn layout and router in `apps/manager/src/renderer/App.tsx`
- [ ] T040 Setup Manager navigation routes per `contracts/manager-ui.md` in `apps/manager/src/renderer/routes.tsx`

**Checkpoint**: gRPC connectivity Manager→Hub→Agent established; foundation ready for user stories

---

## Phase 3: User Story 3 — Establish Manager–Hub–Agent Topology (Priority: P1)

**Goal**: Connect Hubs and Agents via request-and-accept; Manager recognizes fleet topology

**Independent Test**: Deploy Hub on Server A and Agent on Server B, establish connection via CLI, verify Manager sees Hub in Hub List and Agent in Agent Status

- [ ] T041 [US3] Implement Hub `connection_requests` store operations in `apps/hub/internal/store/connections.go`
- [ ] T042 [US3] Implement Agent connection request submission in `apps/agent/internal/grpc/connections.go`
- [ ] T043 [US3] Implement Hub agent connection management (list, request, disconnect, status, last-sync) in `apps/hub/internal/grpc/agent_handlers.go`
- [ ] T044 [US3] Implement Hub manager connection management (list, request, disconnect, test, last-sync) in `apps/hub/internal/grpc/manager_handlers.go`
- [ ] T045 [P] [US3] Implement `nvx-hub agent search|list|connection*` commands in `apps/hub/cmd/nvx-hub/agent.go`
- [ ] T046 [P] [US3] Implement `nvx-hub manager search|list|connection*` commands in `apps/hub/cmd/nvx-hub/manager.go`
- [ ] T047 [P] [US3] Implement `nvx-agent hub search|list|connection*` commands in `apps/agent/cmd/nvx-agent/hub.go`
- [ ] T048 [US3] Implement optional mDNS discovery in `apps/hub/internal/discovery/mdns.go`
- [ ] T049 [P] [US3] Implement optional mDNS discovery in `apps/agent/internal/discovery/mdns.go`
- [ ] T050 [US3] Implement Manager Hub List page in `apps/manager/src/renderer/pages/HubList.tsx`
- [ ] T051 [US3] Implement Manager Agent Status page in `apps/manager/src/renderer/pages/AgentStatus.tsx`
- [ ] T052 [US3] Implement connection accept/reject UI and IPC handlers in `apps/manager/src/main/grpc/connection_manager.ts`

**Checkpoint**: Full topology connectivity with CLI and Manager UI

---

## Phase 4: User Story 1 — Monitor Server Health Across the Fleet (Priority: P1) 🎯 MVP

**Goal**: Display fleet-wide and per-server CPU, GPU, RAM, storage, network, uptime metrics with charts

**Independent Test**: Connect one agent through hub; verify metrics appear in "all" and "per-server" modes with charts within 5 seconds

- [ ] T053 [US1] Implement metrics collector (gopsutil + nvidia-smi GPU fallback) in `apps/agent/internal/collector/metrics.go`
- [ ] T054 [US1] Implement 15s metric collection loop and MetricBatch emission in `apps/agent/internal/collector/scheduler.go`
- [ ] T055 [US1] Implement Hub metric relay stream (Agent→Hub→Manager) in `apps/hub/internal/grpc/metrics_relay.go`
- [ ] T056 [US1] Implement Manager metric ingestion and `metric_history` persistence in `apps/manager/src/main/metrics/ingest.ts`
- [ ] T057 [P] [US1] Implement Monitoring "all servers" aggregate view in `apps/manager/src/renderer/pages/MonitoringDashboard.tsx`
- [ ] T058 [P] [US1] Implement Monitoring "per-server" chart view with Recharts in `apps/manager/src/renderer/pages/MonitoringServerDetail.tsx`
- [ ] T059 [US1] Add stale/offline indicator logic for unavailable servers in `apps/manager/src/renderer/components/ServerStatusBadge.tsx`
- [ ] T060 [US1] Wire TanStack Query metric subscriptions in `apps/manager/src/renderer/hooks/useMetrics.ts`

**Checkpoint**: MVP monitoring functional — fleet and per-server views

---

## Phase 5: User Story 2 — Discover and Manage Server Inventory (Priority: P1)

**Goal**: Server list with hardware specs, alias, OS info, uptime, online status

**Independent Test**: New agent connection populates Server List with CPU, RAM, motherboard, storage, NICs, VM, OS within 30 seconds

- [ ] T061 [US2] Implement hardware/OS inventory collector in `apps/agent/internal/collector/inventory.go`
- [ ] T062 [US2] Implement InventorySnapshot gRPC message emission on connect and periodic refresh in `apps/agent/internal/grpc/inventory.go`
- [ ] T063 [US2] Implement Hub inventory relay to Manager in `apps/hub/internal/grpc/inventory_relay.go`
- [ ] T064 [US2] Implement Manager server upsert and `hardware_json` persistence in `apps/manager/src/main/inventory/service.ts`
- [ ] T065 [P] [US2] Implement Server List page in `apps/manager/src/renderer/pages/ServerList.tsx`
- [ ] T066 [US2] Implement server detail panel with full hardware breakdown in `apps/manager/src/renderer/components/ServerDetailPanel.tsx`
- [ ] T067 [US2] Add offline/last-seen indicator for unavailable servers in `apps/manager/src/renderer/pages/ServerList.tsx`

**Checkpoint**: Server inventory auto-populated from agents

---

## Phase 6: User Story 10 — Operate Hub and Agent via Command Line (Priority: P2)

**Goal**: Full `nvx-hub` and `nvx-agent` CLI for status, service control, log dir, and connections

**Independent Test**: Run all CLI commands from `contracts/hub-cli.md` and `contracts/agent-cli.md`; verify expected output

- [ ] T068 [P] [US10] Implement `nvx-hub service disable|enable|restart` in `apps/hub/cmd/nvx-hub/service.go`
- [ ] T069 [P] [US10] Implement `nvx-hub log set dir` in `apps/hub/cmd/nvx-hub/log.go`
- [ ] T070 [P] [US10] Implement `nvx-agent service disable|enable|restart` in `apps/agent/cmd/nvx-agent/service.go`
- [ ] T071 [P] [US10] Implement `nvx-agent log set dir` in `apps/agent/cmd/nvx-agent/log.go`
- [ ] T072 [US10] Enhance `nvx-hub status` with version, connection counts, last sync in `apps/hub/cmd/nvx-hub/status.go`
- [ ] T073 [US10] Enhance `nvx-agent status` with version, hub count, cache pending in `apps/agent/cmd/nvx-agent/status.go`
- [ ] T074 [US10] Add `--json` output flag to all Hub CLI commands in `apps/hub/cmd/nvx-hub/root.go`
- [ ] T075 [US10] Add `--json` output flag to all Agent CLI commands in `apps/agent/cmd/nvx-agent/root.go`

**Checkpoint**: Complete CLI surface for Hub and Agent

---

## Phase 7: User Story 15 — Receive Server Events in Real Time (Priority: P2)

**Goal**: Agent sends server events; Manager displays real-time event feed per server

**Independent Test**: Restart a service on agent; event appears in Manager feed within 5 seconds

- [ ] T076 [US15] Implement event emitter for state changes in `apps/agent/internal/events/emitter.go`
- [ ] T077 [US15] Implement ServerEvent gRPC stream publishing in `apps/agent/internal/grpc/events.go`
- [ ] T078 [US15] Implement Hub event relay stream in `apps/hub/internal/grpc/events_relay.go`
- [ ] T079 [US15] Implement Manager event ingestion and `server_events` persistence in `apps/manager/src/main/events/ingest.ts`
- [ ] T080 [P] [US15] Implement per-server Event Feed page in `apps/manager/src/renderer/pages/EventFeed.tsx`
- [ ] T081 [US15] Wire real-time event IPC push to renderer in `apps/manager/src/main/events/broadcast.ts`
- [ ] T082 [US15] Implement cached event flush on reconnection in `apps/hub/internal/cache/flush.go` and `apps/agent/internal/cache/flush.go`

**Checkpoint**: Real-time events with offline cache recovery

---

## Phase 8: User Story 11 — Configure Logging and Retention from Manager (Priority: P2)

**Goal**: Configure log directories and per-file size limits per server/hub from Manager

**Independent Test**: Change log dir and size limits from Manager; verify new entries follow updated config

- [ ] T083 [US11] Add `log_configurations` table migration in `apps/manager/src/main/db/migrations/002_log_config.sql`
- [ ] T084 [US11] Implement LogConfigUpdate gRPC command handler in `apps/hub/internal/grpc/log_config.go`
- [ ] T085 [US11] Implement LogConfigUpdate gRPC command handler in `apps/agent/internal/grpc/log_config.go`
- [ ] T086 [US11] Implement log file size limit enforcement and rotation in `apps/hub/internal/logging/rotation.go`
- [ ] T087 [P] [US11] Implement log file size limit enforcement in `apps/agent/internal/logging/rotation.go`
- [ ] T088 [P] [US11] Implement Log Configuration panel in `apps/manager/src/renderer/pages/LogConfigPanel.tsx`
- [ ] T089 [US11] Implement 5-day connectivity log purge job in `apps/hub/internal/logging/retention.go`
- [ ] T090 [P] [US11] Implement 5-day connectivity log purge job in `apps/agent/internal/logging/retention.go`

**Checkpoint**: Centralized log configuration from Manager

---

## Phase 9: User Story 5 — Manage System Services (Priority: P2)

**Goal**: List, enable, disable, and restart systemd services remotely

**Independent Test**: List services on connected server; disable and re-enable a service; verify status updates

- [ ] T091 [US5] Implement systemd service executor in `apps/agent/internal/executor/systemd.go`
- [ ] T092 [US5] Register `service.list|enable|disable|restart` command handlers in `apps/agent/internal/grpc/commands.go`
- [ ] T093 [US5] Implement Hub command relay for service operations in `apps/hub/internal/grpc/command_relay.go`
- [ ] T094 [P] [US5] Implement Service Manager page in `apps/manager/src/renderer/pages/ServiceManager.tsx`
- [ ] T095 [US5] Implement Manager command dispatch IPC in `apps/manager/src/main/commands/dispatch.ts`

**Checkpoint**: Remote systemd service management functional

---

## Phase 10: User Story 6 — Perform Network Diagnostics (Priority: P2)

**Goal**: View NICs and run ping, traceroute, nslookup remotely

**Independent Test**: View NIC list; run ping and traceroute; results display in Network Manager

- [ ] T096 [US6] Implement NIC enumeration in `apps/agent/internal/executor/network_info.go`
- [ ] T097 [US6] Implement ping, traceroute, nslookup executors in `apps/agent/internal/executor/network_diag.go`
- [ ] T098 [US6] Register `network.*` command handlers in `apps/agent/internal/grpc/commands.go`
- [ ] T099 [P] [US6] Implement Network Manager page in `apps/manager/src/renderer/pages/NetworkManager.tsx`
- [ ] T100 [US6] Implement streaming diagnostic output relay in `apps/hub/internal/grpc/command_relay.go`

**Checkpoint**: Network diagnostics executable from Manager

---

## Phase 11: User Story 4 — Manage Docker Containers and Images (Priority: P2)

**Goal**: View Docker containers and images on managed servers

**Independent Test**: Select Docker-enabled server; verify container and image lists; verify graceful message when Docker absent

- [ ] T101 [US4] Implement Docker container/image list executor in `apps/agent/internal/executor/docker.go`
- [ ] T102 [US4] Register `docker.containers|docker.images` command handlers with NVX-008 error in `apps/agent/internal/grpc/commands.go`
- [ ] T103 [P] [US4] Implement Docker Manager page in `apps/manager/src/renderer/pages/DockerManager.tsx`
- [ ] T104 [US4] Add Docker-not-available empty state in `apps/manager/src/renderer/components/DockerUnavailable.tsx`

**Checkpoint**: Docker listing functional with unavailable fallback

---

## Phase 12: User Story 12 — Run Custom and Built-In Server Actions (Priority: P2)

**Goal**: Reboot, datetime, package update/upgrade, and reusable custom commands

**Independent Test**: Run package update and a saved custom command; verify output displayed

- [ ] T105 [US12] Add `custom_commands` table migration in `apps/manager/src/main/db/migrations/003_custom_commands.sql`
- [ ] T106 [US12] Implement built-in action executors (reboot, datetime, apt update/upgrade) in `apps/agent/internal/executor/actions.go`
- [ ] T107 [US12] Implement custom command script executor in `apps/agent/internal/executor/custom.go`
- [ ] T108 [P] [US12] Implement Server Actions page in `apps/manager/src/renderer/pages/ServerActions.tsx`
- [ ] T109 [P] [US12] Implement Custom Commands page in `apps/manager/src/renderer/pages/CustomCommands.tsx`
- [ ] T110 [US12] Implement custom command CRUD service in `apps/manager/src/main/commands/custom_commands.ts`

**Checkpoint**: Built-in actions and custom commands executable from Manager

---

## Phase 13: User Story 13 — Define and Respond to Alert Triggers (Priority: P2)

**Goal**: Define `IF [metric] [op] [threshold] THEN [Email|Reboot] [target]` triggers with automatic evaluation

**Independent Test**: Create CPU > 80 Email trigger; simulate threshold breach; verify email within 60 seconds

- [ ] T111 [US13] Add `alert_triggers` and `smtp_settings` table migrations in `apps/manager/src/main/db/migrations/004_triggers_smtp.sql`
- [ ] T112 [US13] Implement trigger parser and validator in `apps/manager/src/main/triggers/parser.ts`
- [ ] T113 [US13] Implement continuous trigger evaluation engine in `apps/manager/src/main/triggers/evaluator.ts`
- [ ] T114 [US13] Implement SMTP email sender via Nodemailer in `apps/manager/src/main/triggers/email.ts`
- [ ] T115 [US13] Implement Reboot action dispatch via gRPC command in `apps/manager/src/main/triggers/reboot.ts`
- [ ] T116 [P] [US13] Implement Alert Triggers page in `apps/manager/src/renderer/pages/AlertTriggers.tsx`
- [ ] T117 [P] [US13] Implement SMTP Settings section in `apps/manager/src/renderer/pages/Settings.tsx`

**Checkpoint**: Alert triggers evaluate metrics and fire Email/Reboot actions

---

## Phase 14: User Story 7 — Visualize Network Topology with Real-Time Latency (Priority: P2)

**Goal**: Interactive graph with draggable nodes, external IPs, real-time ping latency

**Independent Test**: Display graph with 2+ servers; drag node; add external IP; latency refreshes every 5 seconds

- [ ] T118 [US7] Add `graph_nodes` table migration in `apps/manager/src/main/db/migrations/005_graph_nodes.sql`
- [ ] T119 [US7] Implement inter-server ping latency collector dispatch in `apps/agent/internal/executor/network_diag.go`
- [ ] T120 [US7] Implement graph latency aggregation service in `apps/manager/src/main/graph/latency.ts`
- [ ] T121 [US7] Implement Graph Network page with React Flow in `apps/manager/src/renderer/pages/GraphNetwork.tsx`
- [ ] T122 [US7] Implement node drag persistence and external IP node creation in `apps/manager/src/renderer/components/GraphCanvas.tsx`
- [ ] T123 [US7] Implement 5-second latency refresh polling in `apps/manager/src/renderer/hooks/useGraphLatency.ts`

**Checkpoint**: Interactive topology graph with real-time latency

---

## Phase 15: User Story 14 — Deploy Hub and Agent Remotely (Priority: P2)

**Goal**: Install Hub and Agent on remote Ubuntu servers from Manager via SSH

**Independent Test**: Install agent on fresh Ubuntu VM from Manager; agent appears in Agent Status within 5 minutes

- [ ] T124 [US14] Create `.deb` packaging scripts in `deploy/deb/build-hub-deb.sh` and `deploy/deb/build-agent-deb.sh`
- [ ] T125 [US14] Implement SSH connection helper in `apps/manager/src/main/deploy/ssh_client.ts`
- [ ] T126 [US14] Implement remote Hub install workflow (upload deb, dpkg, systemd enable) in `apps/manager/src/main/deploy/install_hub.ts`
- [ ] T127 [US14] Implement remote Agent install workflow in `apps/manager/src/main/deploy/install_agent.ts`
- [ ] T128 [US14] Implement post-install gRPC handshake verification in `apps/manager/src/main/deploy/verify.ts`
- [ ] T129 [P] [US14] Implement Remote Deploy wizard page in `apps/manager/src/renderer/pages/RemoteDeploy.tsx`

**Checkpoint**: Remote deployment of Hub and Agent from Manager

---

## Phase 16: User Story 8 — Manage Applications on Servers (Priority: P3)

**Goal**: List, install, and uninstall applications via apt

**Independent Test**: List installed apps; install and uninstall a package; verify list updates

- [ ] T130 [US8] Implement apt list/install/uninstall executor in `apps/agent/internal/executor/apt.go`
- [ ] T131 [US8] Register app management command handlers in `apps/agent/internal/grpc/commands.go`
- [ ] T132 [P] [US8] Implement App Manager page in `apps/manager/src/renderer/pages/AppManager.tsx`

**Checkpoint**: Application management via apt functional

---

## Phase 17: User Story 9 — Access Remote Terminal (Priority: P3)

**Goal**: Interactive remote terminal session from Manager via gRPC relay

**Independent Test**: Open terminal for connected server; execute command; output displayed; reconnect on disconnect

- [ ] T133 [US9] Implement SSH shell session handler in `apps/agent/internal/executor/terminal.go`
- [ ] T134 [US9] Implement TerminalFrame bidirectional gRPC stream in `apps/agent/internal/grpc/terminal.go`
- [ ] T135 [US9] Implement Hub terminal stream relay in `apps/hub/internal/grpc/terminal_relay.go`
- [ ] T136 [US9] Implement Manager terminal IPC bridge in `apps/manager/src/main/terminal/session.ts`
- [ ] T137 [P] [US9] Implement Terminal page with xterm.js in `apps/manager/src/renderer/pages/TerminalView.tsx`

**Checkpoint**: Interactive terminal over gRPC relay path

---

## Phase 18: User Story 16 — Customize Server Visual Identity (Priority: P3)

**Goal**: Custom icons and badges for servers, hubs, and agents in lists and graph

**Independent Test**: Upload icon to assets folder; assign to server; verify in Server List and Graph Network

- [ ] T138 [US16] Add `visual_assets` table migration in `apps/manager/src/main/db/migrations/006_visual_assets.sql`
- [ ] T139 [US16] Implement assets folder manager in `apps/manager/src/main/assets/store.ts`
- [ ] T140 [P] [US16] Implement asset upload and assignment UI in `apps/manager/src/renderer/components/AssetManager.tsx`
- [ ] T141 [US16] Integrate icon/badge rendering in `apps/manager/src/renderer/pages/ServerList.tsx` and `apps/manager/src/renderer/components/GraphCanvas.tsx`

**Checkpoint**: Visual customization applied across Manager views

---

## Phase 19: Polish & Cross-Cutting Concerns

**Purpose**: Settings, About, integration validation, and hardening

- [ ] T142 [P] Implement Settings page (general preferences, metric retention) in `apps/manager/src/renderer/pages/Settings.tsx`
- [ ] T143 [P] Implement About page with version info in `apps/manager/src/renderer/pages/About.tsx`
- [ ] T144 Implement Manager IPC type definitions for all channels in `apps/manager/src/shared/ipc.ts`
- [ ] T145 [P] Add error code mapping (NVX-001–NVX-008) in `apps/manager/src/shared/errors.ts`
- [ ] T146 Create smoke test script validating quickstart steps in `scripts/dev/smoke-test.sh`
- [ ] T147 Run full quickstart.md validation against docker-compose dev fixture
- [ ] T148 [P] Update README.md with build, deploy, and development instructions

---

## Dependencies & Execution Order

### Phase Dependencies

| Phase | Depends on | Blocks |
|-------|-----------|--------|
| 1 Setup | — | Phase 2 |
| 2 Foundational | Phase 1 | All user stories (3–18) |
| 3 US3 Topology | Phase 2 | US1, US2, all P2 stories |
| 4 US1 Monitoring | Phase 3 | US13 (triggers need metrics) |
| 5 US2 Inventory | Phase 3 | — |
| 6–18 User Stories | Phase 3 (most also benefit from US15 events) | — |
| 19 Polish | Desired stories complete | — |

### User Story Dependencies

| Story | Priority | Depends on | Notes |
|-------|----------|-----------|-------|
| US3 | P1 | Foundational | Topology backbone — implement first |
| US1 | P1 | US3 | MVP monitoring |
| US2 | P1 | US3 | Inventory on connect |
| US10 | P2 | Foundational | Can parallel with US3+ |
| US15 | P2 | US3 | Events enhance all management stories |
| US11 | P2 | US3 | Log config via gRPC |
| US5 | P2 | US3, US15 | Service events |
| US6 | P2 | US3 | Network diag |
| US4 | P2 | US3 | Docker listing |
| US12 | P2 | US3, US15 | Actions emit events |
| US13 | P2 | US1 | Triggers need metrics stream |
| US7 | P2 | US1, US2 | Graph needs servers + latency |
| US14 | P2 | US3 | Remote install + verify |
| US8 | P3 | US3 | apt management |
| US9 | P3 | US3 | Terminal relay |
| US16 | P3 | US2, US7 | Assets in list + graph |

### Parallel Opportunities

- **Phase 1**: T003–T012 all marked [P] can run in parallel after T001–T002
- **Phase 2**: TLS (T017–T019), logging (T020–T023), cache (T024–T025), CLI skeletons (T035–T036) parallelizable
- **After Phase 3**: US10 CLI (Phase 6) can run parallel to US1/US2 (Phases 4–5)
- **P2 stories**: US4, US5, US6, US8 have no cross-dependencies — parallelizable after US3
- **P3 stories**: US8, US9, US16 parallelizable after their respective dependencies

---

## Parallel Example: User Story 3

```bash
# Parallel CLI commands (different files):
T045: apps/hub/cmd/nvx-hub/agent.go
T046: apps/hub/cmd/nvx-hub/manager.go
T047: apps/agent/cmd/nvx-agent/hub.go

# Parallel discovery (different apps):
T048: apps/hub/internal/discovery/mdns.go
T049: apps/agent/internal/discovery/mdns.go
```

---

## Parallel Example: P2 Management Stories

```bash
# After US3 complete, launch in parallel:
Developer A: US5 Service Manager (T091–T095)
Developer B: US6 Network Manager (T096–T100)
Developer C: US4 Docker Manager (T101–T104)
```

---

## Implementation Strategy

### MVP First (P1 stories only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (**critical**)
3. Complete Phase 3: US3 Topology
4. Complete Phase 4: US1 Monitoring
5. Complete Phase 5: US2 Inventory
6. **STOP and VALIDATE** against quickstart.md steps 3–6

### Incremental Delivery

| Increment | Stories | Value delivered |
|-----------|---------|-----------------|
| MVP | US3 + US1 + US2 | Connected fleet with monitoring and inventory |
| Ops | US10 + US5 + US6 + US4 | CLI + service/network/docker management |
| Automation | US12 + US13 + US15 | Actions, triggers, events |
| Visual | US7 + US16 | Topology graph and custom assets |
| Deploy | US14 | Remote onboarding |
| Advanced | US8 + US9 | App management and terminal |

### Suggested MVP Scope

**User Story 1 (Monitoring)** — requires US3 (topology) and US2 (inventory) for full P1 experience:

- Minimum viable path: **Setup → Foundational → US3 → US1**
- Recommended P1 demo: add **US2** for complete server inventory

---

## Notes

- All inter-component communication MUST use gRPC (FR-006c) — no alternate protocols
- Total tasks: **148**
- Commit after each task or logical group; validate at each checkpoint
- Run `make proto` after any `.proto` changes before building Hub, Agent, or Manager
