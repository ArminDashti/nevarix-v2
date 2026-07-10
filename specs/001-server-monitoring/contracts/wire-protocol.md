# Wire Protocol Contract

**Version**: 1.0.0 | **Transport**: gRPC over mTLS | **Schema**: Protocol Buffers v3

## Services

### `NvxHubService` (Manager → Hub, Agent → Hub)

```protobuf
service NvxHubService {
  rpc ConnectManager(stream ManagerFrame) returns (stream HubFrame);
  rpc ConnectAgent(stream AgentFrame) returns (stream HubAgentFrame);
  rpc SubmitConnectionRequest(ConnectionRequest) returns (ConnectionRequestAck);
  rpc AcceptConnectionRequest(AcceptRequest) returns (AcceptResponse);
}
```

### `NvxAgentService` (Hub → Agent)

```protobuf
service NvxAgentService {
  rpc Connect(stream HubAgentCommand) returns (stream AgentResponse);
  rpc SubmitConnectionRequest(ConnectionRequest) returns (ConnectionRequestAck);
}
```

---

## Core Messages

### MetricBatch

```protobuf
message MetricBatch {
  string agent_name = 1;
  string server_identity = 2;
  google.protobuf.Timestamp recorded_at = 3;
  double cpu_percent = 4;
  optional double gpu_percent = 5;  // absent = GPU unavailable
  double ram_percent = 6;
  double storage_percent = 7;
  uint64 network_bytes_in = 8;
  uint64 network_bytes_out = 9;
  uint64 uptime_seconds = 10;
}
```

**Frequency**: Every 15 seconds from Agent.

---

### ServerEvent

```protobuf
message ServerEvent {
  string event_id = 1;
  string agent_name = 2;
  string server_identity = 3;
  string event_type = 4;       // e.g. "service.restarted", "connectivity.lost"
  google.protobuf.Timestamp occurred_at = 5;
  string payload_json = 6;
}
```

**Delivery**: Real-time stream; cached in outbound spool on failure.

---

### CommandRequest / CommandResponse

```protobuf
message CommandRequest {
  string request_id = 1;
  string command_type = 2;   // "service.restart", "docker.list", "custom", "reboot", etc.
  string target_agent_name = 3;
  string payload_json = 4;
}

message CommandResponse {
  string request_id = 1;
  bool success = 2;
  string output = 3;
  string error_message = 4;
  google.protobuf.Timestamp completed_at = 5;
}
```

**Command types** (non-exhaustive):

| `command_type` | Payload | Agent action |
|----------------|---------|--------------|
| `service.list` | `{}` | `systemctl list-units` |
| `service.enable` | `{"name":"nginx"}` | systemctl enable + start |
| `service.disable` | `{"name":"nginx"}` | systemctl disable + stop |
| `service.restart` | `{"name":"nginx"}` | systemctl restart |
| `docker.containers` | `{}` | docker ps -a |
| `docker.images` | `{}` | docker images |
| `network.ping` | `{"host":"8.8.8.8"}` | ping execution |
| `network.traceroute` | `{"host":"8.8.8.8"}` | traceroute |
| `network.nslookup` | `{"host":"example.com"}` | nslookup/dig |
| `action.reboot` | `{}` | reboot |
| `action.datetime` | `{"datetime":"ISO8601"}` | timedatectl |
| `action.package_update` | `{}` | apt update |
| `action.package_upgrade` | `{}` | apt upgrade -y |
| `custom` | `{"script":"..."}` | execute saved script |
| `log.config` | `LogConfigUpdate` | update log settings |

---

### LogConfigUpdate

```protobuf
message LogConfigUpdate {
  string log_root_dir = 1;           // default: /var/log/nvx
  int32 service_txt_limit_mb = 2;
  int32 errors_txt_limit_mb = 3;
  int32 connectivity_txt_limit_mb = 4;
}
```

---

### ConnectionRequest (request-and-accept)

```protobuf
message ConnectionRequest {
  string request_id = 1;
  string requester_name = 2;
  string requester_type = 3;   // "manager" | "agent" | "hub"
  string target_name = 4;
  string target_type = 5;
  bytes certificate_pem = 6;
  google.protobuf.Timestamp requested_at = 7;
}

message ConnectionRequestAck {
  string request_id = 1;
  string status = 2;           // "pending" | "accepted" | "rejected"
}
```

---

### TerminalFrame (bidirectional)

```protobuf
message TerminalFrame {
  string session_id = 1;
  string agent_name = 2;
  bytes data = 3;              // stdin/stdout binary chunk
  bool close = 4;
}
```

---

### InventorySnapshot

```protobuf
message InventorySnapshot {
  string server_alias = 1;
  string cpu_model = 2;
  int32 cpu_cores = 3;
  uint64 ram_bytes = 4;
  string motherboard = 5;
  repeated StorageDevice storage = 6;
  repeated NetworkInterface nics = 7;
  bool is_vm = 8;
  string os_name = 9;
  string os_version = 10;
  uint64 uptime_seconds = 11;
}
```

---

## Connectivity Log Format (file contract)

Written to `connectivity.txt` every 60 seconds by Hub and Agent:

```text
<ISO8601-timestamp> | direction=<inbound|outbound> | peer=<name> | status=<connected|disconnected|failed> | latency_ms=<integer> | detail=<message>
```

**Example**:

```text
2026-07-10T08:30:00Z | direction=outbound | peer=hub-2 | status=connected | latency_ms=12 | detail=grpc stream active
```

**Retention**: 5 days; purged by daily cleanup job.

---

## Error Codes

| Code | Name | Description |
|------|------|-------------|
| `NVX-001` | `PEER_NOT_FOUND` | Target agent/hub/manager unknown |
| `NVX-002` | `CONNECTION_PENDING` | Request awaiting accept |
| `NVX-003` | `CONNECTION_REJECTED` | Request denied |
| `NVX-004` | `ROUTE_UNAVAILABLE` | Hub cannot reach agent |
| `NVX-005` | `COMMAND_TIMEOUT` | Agent did not respond in 30s |
| `NVX-006` | `CACHE_OVERFLOW` | Outbound spool evicting oldest |
| `NVX-007` | `GPU_UNAVAILABLE` | GPU metric/trigger skipped |
| `NVX-008` | `DOCKER_NOT_INSTALLED` | Docker commands rejected |

---

## Versioning

- Protocol version sent in gRPC metadata header `nvx-protocol-version: 1.0.0`.
- Breaking changes increment MAJOR; backward-compatible additions increment MINOR.
