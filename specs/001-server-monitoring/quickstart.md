# Quickstart: Server Monitoring Platform

**Feature**: `001-server-monitoring` | **Purpose**: End-to-end validation guide

## Prerequisites

- Ubuntu 22.04+ VMs or containers: 1 Hub host, 2 Agent hosts (or docker-compose)
- Go 1.23+, Node.js 20+, Docker (optional, for compose fixture)
- OpenSSL for test CA generation
- Manager dev environment (Electron) on operator workstation

---

## 1. Build Components

```bash
# From repository root (after implementation)
make proto          # Generate protobuf stubs
make build-hub      # → bin/nvx-hub
make build-agent    # → bin/nvx-agent
make build-manager  # → apps/manager release build
```

---

## 2. Generate Test TLS Certificates

```bash
./scripts/dev/generate-certs.sh
# Outputs to ./certs/ — deploy to /etc/nvx/certs/ on each host
```

---

## 3. Install Hub (manual or CLI)

On **hub-host**:

```bash
sudo cp bin/nvx-hub /usr/local/bin/
sudo cp configs/hub.yaml /etc/nvx/hub.yaml
sudo cp -r certs/* /etc/nvx/certs/
sudo cp deploy/nvx-hub.service /etc/systemd/system/
sudo systemctl enable --now nvx-hub
nvx-hub status
```

**Expected**: Status `running`, gRPC listening on `:9443`.

---

## 4. Install Agents

On **agent-1** and **agent-2**:

```bash
sudo cp bin/nvx-agent /usr/local/bin/
sudo cp configs/agent.yaml /etc/nvx/agent.yaml   # unique agent_name each
sudo cp -r certs/* /etc/nvx/certs/
sudo cp deploy/nvx-agent.service /etc/systemd/system/
sudo systemctl enable --now nvx-agent
nvx-agent status
```

---

## 5. Establish Connections (CLI)

On each agent:

```bash
nvx-agent hub connection request --hub=hub-2
```

On hub:

```bash
nvx-hub agent list                              # verify pending
nvx-hub agent connection --agent=agent-1      # verify connected
```

On hub (register with Manager):

```bash
nvx-hub manager connection request --manager=manager-local
```

---

## 6. Launch Manager

```bash
cd apps/manager
npm run dev
```

Accept hub connection request in Manager UI → **Hubs** page.

**Validate** (User Story 1, 2, 3):

- [ ] Server List shows both agents with hardware info within 30s (SC-003)
- [ ] Monitoring "all" mode shows fleet metrics within 5s (SC-001)
- [ ] Monitoring "per-server" mode shows 24h charts (SC-002)

---

## 7. Validate Logging

On hub-host:

```bash
ls /var/log/nvx/$(date +%Y/%m/%d)/
# Expected: service.txt  errors.txt  connectivity.txt

tail -3 /var/log/nvx/$(date +%Y/%m/%d)/connectivity.txt
# Expected: pipe-delimited entries, ~1/minute
```

In Manager → server → **Logs**:

- [ ] Change log root directory; verify new entries use new path
- [ ] Set size limit; verify rotation on next write cycle

---

## 8. Validate Server Actions

In Manager → server → **Actions**:

```bash
# Via UI buttons
```

- [ ] Package update completes with apt output (FR-045)
- [ ] Custom command: save `echo hello` → run → output displayed (FR-044, SC-016)

---

## 9. Validate Alert Trigger

Create trigger in Manager:

```text
IF CPU > 80 THEN Email test@example.com
```

Load CPU on agent (stress tool) → verify email within 60s (SC-013).

---

## 10. Validate Events

Restart a service on agent-1 via Manager Service Manager.

- [ ] Event `service.state_changed` appears in event feed within 5s (SC-014)

---

## 11. Validate Graph Network

- [ ] All servers appear as nodes with ping latency
- [ ] Drag node → reload → position persisted (FR-025)
- [ ] Add external IP `8.8.8.8` → node shows ping (FR-024)
- [ ] Latency refreshes every ≤5s (SC-005)

---

## 12. Validate Connection Resilience

```bash
# On hub-host
sudo systemctl stop nvx-hub
# Wait 2 minutes, generate metrics on agent
sudo systemctl start nvx-hub
```

- [ ] Agent cache pending > 0 while hub down
- [ ] Cached metrics/events flush within 30s of hub restart (SC-012)
- [ ] `errors.txt` logs cache overflow warning if limit exceeded

---

## 13. Validate Remote Deploy (optional)

In Manager → **Deploy**:

- [ ] Install agent on fresh Ubuntu VM via SSH
- [ ] Installation completes within 5 minutes (SC-015)
- [ ] Agent appears in Agent Status without manual CLI steps

---

## 14. Docker Compose Fixture (dev)

```bash
docker compose -f deploy/docker-compose.dev.yml up -d
# Spins up: hub-1, agent-1, agent-2 on isolated network
./scripts/dev/smoke-test.sh
```

**Expected exit code**: 0 with all contract checks passing.

---

## Contract References

| Topic | Document |
|-------|----------|
| gRPC messages | [wire-protocol.md](./contracts/wire-protocol.md) |
| Hub CLI | [hub-cli.md](./contracts/hub-cli.md) |
| Agent CLI | [agent-cli.md](./contracts/agent-cli.md) |
| Manager screens | [manager-ui.md](./contracts/manager-ui.md) |
| Database schema | [data-model.md](./data-model.md) |

---

## Success Criteria Checklist

| ID | Target | Validation step |
|----|--------|-----------------|
| SC-001 | Metrics in ≤5s | Step 6 |
| SC-003 | Inventory in ≤30s | Step 6 |
| SC-005 | Graph refresh ≤5s | Step 11 |
| SC-011 | Connectivity log 1/min | Step 7 |
| SC-012 | Cache flush ≤30s | Step 12 |
| SC-013 | Trigger fire ≤60s | Step 9 |
| SC-014 | Events in ≤5s | Step 10 |
| SC-015 | Remote install ≤5min | Step 13 |
| SC-016 | Custom cmd immediate | Step 8 |
