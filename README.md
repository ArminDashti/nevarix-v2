# nevarix-v2

Server monitoring and management platform (Manager → Hub → Agent).

## Prerequisites

- Go 1.23+
- Node.js 20+
- buf CLI (`brew install bufbuild/buf/buf` or see https://buf.build/docs/installation)
- Docker (optional, for dev compose)
- OpenSSL (for dev cert generation)

## Quick Start

```bash
# Generate dev TLS certificates
make certs

# Build Hub and Agent
make build

# Generate protobuf stubs (requires buf)
make proto

# Manager UI (dev)
cd apps/manager
npm install
npm run dev
```

## Project Structure

```text
apps/manager/     Electron + React + TypeScript + Shadcn UI
apps/hub/         Go Hub service (nvx-hub)
apps/agent/       Go Agent service (nvx-agent)
packages/proto/   gRPC/protobuf contracts
configs/          Default YAML configs
deploy/           systemd units, Docker, compose
specs/            Feature specifications and tasks
```

## Documentation

- [Implementation plan](specs/001-server-monitoring/plan.md)
- [Tasks](specs/001-server-monitoring/tasks.md)
- [Quickstart validation](specs/001-server-monitoring/quickstart.md)

## CLI

```bash
nvx-hub status
nvx-agent status
```

See `specs/001-server-monitoring/contracts/` for full CLI and wire protocol contracts.
