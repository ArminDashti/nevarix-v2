# Implementation Plan: Server Monitoring Platform

**Branch**: `001-server-monitoring` | **Date**: 2026-07-10 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-server-monitoring/spec.md`

## Summary

Build a three-tier server monitoring and management platform (**Manager → Hub → Agent**) targeting Ubuntu servers. The **Manager** is an Electron + React + TypeScript desktop app using Shadcn UI for fleet monitoring, remote administration, alert triggers, and topology visualization. **Hub** and **Agent** are Go services with Cobra CLIs (`nvx-hub`, `nvx-agent`), communicating via gRPC/mTLS, with SQLite persistence on Hub, file-spool outbound caching, and structured daily log files under `/var/log/nvx`.

## Technical Context

**Language/Version**: TypeScript 5.x (Manager), Go 1.23+ (Hub/Agent), Protocol Buffers 3

**Primary Dependencies**:
- Manager: Electron 33, React 19, Vite, Shadcn UI, TanStack Query, Recharts, React Flow, xterm.js, `@grpc/grpc-js`, better-sqlite3, Nodemailer
- Hub/Agent: Cobra, gRPC, modernc.org/sqlite (Hub), gopsutil, golang.org/x/crypto/ssh, slog

**Storage**:
- Manager: SQLite (`manager.db`) + assets folder
- Hub: SQLite (`/var/lib/nvx/hub.db`)
- Hub/Agent: file-spool cache (`/var/cache/nvx/outbound/`), daily log files (`/var/log/nvx/`)

**Testing**: Go `testing`+testify, Vitest (Manager unit), Playwright (Manager E2E), docker-compose integration fixture

**Target Platform**: Manager — Windows/macOS/Linux desktop; Hub/Agent — Ubuntu 22.04+ (systemd)

**Project Type**: Multi-component monorepo (desktop app + 2 Go services + shared proto)

**Performance Goals**: 100-server fleet per Manager (SC-008); metrics every 15s; graph refresh 5s; CLI response <2s

**Constraints**: mTLS required; Manager reaches Agents only via Hub; one Agent + one Hub max per server; outbound cache 512 MB / 24h

**Scale/Scope**: 3 applications, ~15 Manager screens, 65 functional requirements, 16 user stories

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

No project constitution ratified (`.specify/memory/constitution.md` absent). Applied default gates from constitution template:

| Gate | Status | Notes |
|------|--------|-------|
| Components independently testable | ✅ PASS | Hub, Agent, Manager each have CLI/status contracts |
| Contract tests for inter-service comms | ✅ PASS | Protobuf schemas + `contracts/wire-protocol.md` |
| Structured logging required | ✅ PASS | slog + spec-defined log file layout |
| Complexity justified | ✅ PASS | Three-tier topology mandated by spec (FR-001) |
| Security baseline | ✅ PASS | mTLS + request-and-accept (research §8) |

**Post-design re-check**: All gates pass. No unjustified violations.

## Project Structure

### Documentation (this feature)

```text
specs/001-server-monitoring/
├── plan.md              # This file
├── research.md          # Phase 0 — technology decisions
├── data-model.md        # Phase 1 — entities and schemas
├── quickstart.md        # Phase 1 — validation guide
├── contracts/           # Phase 1 — interface contracts
│   ├── wire-protocol.md
│   ├── hub-cli.md
│   ├── agent-cli.md
│   └── manager-ui.md
└── tasks.md             # Phase 2 (/speckit-tasks — not yet created)
```

### Source Code (repository root)

```text
apps/
├── manager/                    # Electron + React + TypeScript + Shadcn
│   ├── src/
│   │   ├── main/               # Electron main: gRPC, SSH, SQLite, SMTP
│   │   ├── renderer/           # React UI pages and Shadcn components
│   │   ├── preload/            # Context bridge IPC
│   │   └── shared/             # Types shared main/renderer
│   └── tests/
│       ├── unit/
│       └── e2e/
├── hub/                        # Go Hub service
│   ├── cmd/nvx-hub/            # Cobra CLI entrypoint
│   ├── internal/
│   │   ├── grpc/               # gRPC server + relay
│   │   ├── store/              # SQLite repository
│   │   ├── cache/              # Outbound file spool
│   │   ├── logging/            # Daily log writer + connectivity ticker
│   │   └── discovery/          # mDNS search
│   └── tests/
└── agent/                      # Go Agent service
    ├── cmd/nvx-agent/          # Cobra CLI entrypoint
    ├── internal/
    │   ├── grpc/               # gRPC client/server
    │   ├── collector/          # Metrics (gopsutil + nvidia-smi)
    │   ├── executor/           # Command execution (systemd, docker, apt)
    │   ├── cache/              # Outbound file spool
    │   ├── logging/
    │   └── events/             # Event emitter
    └── tests/

packages/
└── proto/                      # .proto files + generated Go/TS stubs
    ├── nvx/v1/
    └── buf.yaml

deploy/
├── nvx-hub.service
├── nvx-agent.service
├── docker-compose.dev.yml
└── deb/                        # .deb packaging scripts

scripts/
├── dev/
│   ├── generate-certs.sh
│   └── smoke-test.sh
└── build/

configs/
├── hub.yaml
└── agent.yaml

Makefile                        # build, proto, test, compose-up
```

**Structure Decision**: Monorepo with three `apps/` subprojects plus shared `packages/proto`. Matches three spec-defined components while keeping wire contracts in one protobuf module. Manager renderer follows Shadcn page-per-feature routing aligned with `contracts/manager-ui.md`.

## Complexity Tracking

> No constitution violations requiring justification. Three-component architecture is spec-mandated (FR-001).

| Decision | Why Needed | Simpler Alternative Rejected |
|----------|------------|------------------------------|
| gRPC streaming | Terminal I/O, events, metrics push | REST polling — latency + streaming poor fit |
| Hub SQLite | Connection state, routing, queue index | JSON files — no transactional integrity |
| File-spool cache | Large payload survival across restarts | In-memory queue — lost on crash |
| Electron main process gRPC | Native TLS + SSH from desktop | Renderer-only — blocked by browser sandbox |

## Phase 0: Research — Complete

All NEEDS CLARIFICATION items resolved. See [research.md](./research.md).

Key decisions:
- gRPC + protobuf + mTLS inter-component protocol
- Go 1.23 + Cobra for Hub/Agent; Electron + Shadcn for Manager
- Hub SQLite at `/var/lib/nvx/hub.db`; outbound cache at `/var/cache/nvx/outbound/`
- Default logs at `/var/log/nvx/<YYYY>/<MM>/<DD>/`
- Manager-side trigger evaluation + Nodemailer SMTP
- 15s metric collection interval; nvidia-smi for GPU

## Phase 1: Design — Complete

### Data Model

See [data-model.md](./data-model.md) — Manager SQLite (11 tables), Hub SQLite (6 tables), protobuf wire entities, file-spool envelope format.

### Contracts

| Contract | Path | Covers |
|----------|------|--------|
| Wire protocol | [contracts/wire-protocol.md](./contracts/wire-protocol.md) | gRPC services, messages, error codes |
| Hub CLI | [contracts/hub-cli.md](./contracts/hub-cli.md) | All `nvx-hub` commands |
| Agent CLI | [contracts/agent-cli.md](./contracts/agent-cli.md) | All `nvx-agent` commands |
| Manager UI | [contracts/manager-ui.md](./contracts/manager-ui.md) | Routes, screens, Shadcn requirements |

### Quickstart

See [quickstart.md](./quickstart.md) — 14-step validation guide covering P1 stories and SC-001 through SC-016.

### Agent Context Update

No agent-context script present in `.specify/scripts/`. Skipped — plan artifacts serve as implementation context.

## Phase 2: Task Breakdown — Next

Run `/speckit-tasks` to generate `tasks.md` with dependency-ordered implementation tasks.

Recommended implementation order:
1. `packages/proto` — protobuf definitions and code generation
2. `apps/agent` — metrics, events, CLI, logging, cache
3. `apps/hub` — relay, SQLite, connection management, logging, cache
4. `apps/manager` — gRPC client, SQLite, core monitoring + server list
5. Manager feature modules — Docker, network, services, terminal, graph, triggers, deploy
6. Integration tests via docker-compose + quickstart validation

## Risk Register

| Risk | Mitigation |
|------|------------|
| Electron gRPC native addon complexity | Use `@grpc/grpc-js` pure JS in main process |
| 100-server metric volume | Batch metrics; Hub relays aggregated streams |
| GPU detection variance | Graceful `gpu_status: unavailable`; skip GPU triggers |
| Remote install SSH failures | Structured error codes (FR-053); idempotent .deb install |
| mTLS cert rotation | Settings UI for cert import; CLI documented in quickstart |
