# Feature Specification: Server Monitoring Platform

**Feature Branch**: `001-server-monitoring`

**Created**: 2026-07-10

**Status**: Draft

**Input**: User description: "A project for monitoring servers consisting of three components (Manager, Hub, Agent) with hierarchical routing (Manager → Hub → Agent), multi-hub agent membership, per-server deployment constraints, and a desktop management application with monitoring, server administration, Docker/network/service/app management, terminal access, topology visualization, and CLI-operable Hub/Agent services."

## Clarifications

### Session 2026-07-10

- Q: What naming conventions should apply across components? → A: Use standard, content-related names for entities, variables, and configuration keys (e.g., `server_alias`, `hub_name`, `agent_name`, `connectivity_status`).
- Q: What UI component approach should the Manager use? → A: Shadcn UI component library for all Manager screens and interactions.
- Q: How should Hub and Agent logs be structured and where should they be stored by default? → A: Default log root `/var/log/nvx`; daily subfolders `<DIR>/<YYYY>/<MM>/<DD>/` containing `service.txt`, `errors.txt`, and `connectivity.txt`; connectivity entries written every 1 minute and retained for 5 days.
- Q: Can operators control log file size limits from the Manager? → A: Yes — per server or hub, operators can set size limits independently for `service.txt`, `errors.txt`, and `connectivity.txt`.
- Q: Can operators change log directories from the Manager? → A: Yes — operators can configure the log root directory per server or hub from the Manager.
- Q: Should operators be able to define reusable commands? → A: Yes — operators can define, store, and execute custom commands from the Manager.
- Q: What built-in server administration actions are required beyond service management? → A: Reboot, datetime configuration, package update, package upgrade, and other common administrative actions.
- Q: How should automated alert triggers be defined? → A: Format `IF [CPU|GPU|RAM|STORAGE] [>|<] <threshold> THEN [Email|Reboot] <target>` where Email target is an email address and Reboot target is a server.
- Q: Which operating systems must Hub and Agent support? → A: Ubuntu only (initial release).
- Q: Where should custom visual assets (icons, badges) be stored? → A: A dedicated Manager assets folder for icons and badges applied to servers, hubs, and agents.
- Q: How should Hub and Agent handle data when upstream connections fail? → A: Both MUST persist unsent outbound data in a local temp/cache store and flush it automatically when the connection is restored.
- Q: How should the Hub persist operational data locally? → A: Each Hub maintains a local SQLite database for connections, routing state, and operational records.
- Q: Can Hub and Agent be installed remotely? → A: Yes — the Manager can remotely install Hub and Agent on target Ubuntu servers.
- Q: Does the Agent report events in addition to metrics? → A: Yes — the Agent sends server events (state changes, action results, alerts) to the Manager via Hub in addition to periodic metrics.
- Q: What protocol must be used for connections between Manager, Hub, and Agent? → A: gRPC — all inter-component communication (Manager↔Hub and Hub↔Agent) MUST use gRPC for metrics, events, commands, and terminal relay.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Monitor Server Health Across the Fleet (Priority: P1)

As an infrastructure operator, I want to view real-time and historical health metrics (CPU, GPU, RAM, storage, network, uptime) for all managed servers in aggregate and per-server views so that I can quickly identify performance issues and capacity constraints.

**Why this priority**: Monitoring is the core value proposition. Without visibility into server health, no other management capability delivers meaningful operational value.

**Independent Test**: Can be fully tested by connecting at least one server through the Manager → Hub → Agent chain and verifying that metrics appear in both "all servers" and "per-server" views, including chart display for a selected server.

**Acceptance Scenarios**:

1. **Given** at least one server is connected and reporting, **When** the operator opens the Monitoring view in "all" mode, **Then** aggregated CPU, GPU, RAM, storage, network, and uptime metrics are displayed for the entire fleet.
2. **Given** at least one server is connected and reporting, **When** the operator switches to "per-server" mode and selects a server, **Then** that server's metrics are displayed with charts showing trends over time.
3. **Given** a server stops reporting metrics, **When** the operator views the Monitoring screen, **Then** the server is clearly marked as unavailable and the last known metrics timestamp is shown.

---

### User Story 2 - Discover and Manage Server Inventory (Priority: P1)

As an infrastructure operator, I want a server list showing hardware specifications, alias name, CPU, RAM, motherboard, storage, network interfaces, virtualization status, OS information, and uptime so that I have a single inventory reference for all managed machines.

**Why this priority**: Server inventory is foundational for every other management action and provides context for monitoring alerts and operational decisions.

**Independent Test**: Can be fully tested by verifying that each connected server appears in the list with accurate hardware and OS details retrieved from its agent.

**Acceptance Scenarios**:

1. **Given** a newly connected server, **When** the operator opens the Server List, **Then** the server appears with alias name, CPU, RAM, motherboard, storage, NICs, VM indicator, OS info, and uptime.
2. **Given** multiple servers are connected, **When** the operator views the Server List, **Then** each server is listed as a distinct entry regardless of whether it also hosts a Hub on the same machine.
3. **Given** a server goes offline, **When** the operator views the Server List, **Then** the server entry remains visible with an offline status indicator and last-seen timestamp.

---

### User Story 3 - Establish Manager–Hub–Agent Topology (Priority: P1)

As an infrastructure operator, I want to connect Hubs and Agents through a request-and-accept connection model so that the Manager can reach agents via their designated Hub while agents can participate in multiple Hub memberships.

**Why this priority**: The three-tier topology is the architectural backbone. All monitoring and remote management flows depend on reliable Hub and Agent connectivity.

**Independent Test**: Can be fully tested by deploying Hub and Agent on separate servers, initiating connection requests from both sides, and verifying the Manager can reach the Agent through the Hub path.

**Acceptance Scenarios**:

1. **Given** a Hub is running on Server A and an Agent on Server B, **When** the Agent sends a connection request to the Hub and the Hub accepts, **Then** the Agent appears in the Hub's agent list and reports as connected.
2. **Given** a Manager is configured, **When** the Hub sends a connection request to the Manager and the Manager accepts, **Then** the Hub appears in the Manager's Hub list with name, server, and description.
3. **Given** an Agent is connected to Hub-2 and Hub-5, **When** the Manager views agent status, **Then** the Agent is associated with one recognized server identity while its multi-hub memberships are visible at the Hub level.
4. **Given** a server hosts both a Hub and an Agent, **When** both services are running, **Then** each is managed independently while the Manager recognizes the server as a single entity.

---

### User Story 4 - Manage Docker Containers and Images (Priority: P2)

As an infrastructure operator, I want to view Docker containers and images on managed servers so that I can inspect containerized workloads without logging in manually.

**Why this priority**: Container management is a high-frequency operational task for teams running containerized workloads, but it depends on established server connectivity (P1).

**Independent Test**: Can be fully tested by selecting a Docker-enabled server and verifying container and image lists are displayed with accurate status information.

**Acceptance Scenarios**:

1. **Given** a server with Docker installed and running, **When** the operator opens Docker Manager for that server, **Then** a list of running and stopped containers is displayed.
2. **Given** a server with Docker images present, **When** the operator opens Docker Manager for that server, **Then** a list of available images is displayed.
3. **Given** a server without Docker installed, **When** the operator opens Docker Manager, **Then** a clear message indicates Docker is not available on that server.

---

### User Story 5 - Manage System Services (Priority: P2)

As an infrastructure operator, I want to list, enable, disable, and restart system services on managed servers so that I can perform routine service administration remotely.

**Why this priority**: Service management is a common operational need that reduces the need for direct server login.

**Independent Test**: Can be fully tested by listing services on a connected server and performing enable, disable, and restart actions with verified state changes.

**Acceptance Scenarios**:

1. **Given** a connected server, **When** the operator opens Service Manager, **Then** a list of system services with their current status (enabled/disabled, running/stopped) is displayed.
2. **Given** a stopped service, **When** the operator enables and starts it, **Then** the service status updates to enabled and running.
3. **Given** a running service, **When** the operator disables it, **Then** the service stops and its status updates to disabled.

---

### User Story 6 - Perform Network Diagnostics (Priority: P2)

As an infrastructure operator, I want to view network interfaces and run diagnostic tools (ping, traceroute, DNS lookup) on managed servers so that I can troubleshoot connectivity issues remotely.

**Why this priority**: Network diagnostics are essential for incident response and complement the monitoring and topology features.

**Independent Test**: Can be fully tested by selecting a server, viewing its NIC list, and running ping, traceroute, and DNS lookup against a target with results displayed in the Manager.

**Acceptance Scenarios**:

1. **Given** a connected server, **When** the operator opens Network Manager, **Then** all network interfaces with their configuration details are listed.
2. **Given** a connected server, **When** the operator runs ping against a target address, **Then** ping results (latency, packet loss) are displayed in real time.
3. **Given** a connected server, **When** the operator runs traceroute or DNS lookup, **Then** results are displayed within the Network Manager view.

---

### User Story 7 - Visualize Network Topology with Real-Time Latency (Priority: P2)

As an infrastructure operator, I want an interactive graph showing all servers as nodes with real-time ping times between them so that I can visually assess network health and latency across the fleet.

**Why this priority**: The topology graph provides intuitive situational awareness that complements tabular monitoring data.

**Independent Test**: Can be fully tested by displaying the graph with at least two connected servers and verifying that ping times update in real time and node positions can be rearranged.

**Acceptance Scenarios**:

1. **Given** multiple servers are connected, **When** the operator opens the Graph Network view, **Then** each server appears as a node with connections showing current ping latency.
2. **Given** the Graph Network view is open, **When** ping times change, **Then** latency values on the graph update in real time without requiring a page refresh.
3. **Given** the Graph Network view is displayed, **When** the operator drags a node to a new position, **Then** the node remains at the chosen position on subsequent visits.
4. **Given** the operator wants to monitor an external endpoint, **When** they add an external IP address, **Then** it appears as a node in the graph with ping latency displayed.

---

### User Story 8 - Manage Applications on Servers (Priority: P3)

As an infrastructure operator, I want to list, install, and uninstall applications on managed servers so that I can perform software management without direct server access.

**Why this priority**: Application management extends operational control but is less critical than monitoring, services, and network diagnostics for initial release.

**Independent Test**: Can be fully tested by listing installed applications on a server and performing an install and uninstall operation with verified results.

**Acceptance Scenarios**:

1. **Given** a connected server, **When** the operator opens App Manager, **Then** a list of installed applications is displayed.
2. **Given** a connected server, **When** the operator initiates an application install, **Then** installation progress is shown and the application appears in the list upon completion.
3. **Given** an installed application, **When** the operator initiates uninstall, **Then** the application is removed and no longer appears in the list.

---

### User Story 9 - Access Remote Terminal (Priority: P3)

As an infrastructure operator, I want to open an interactive terminal session to a managed server from the Manager so that I can run ad-hoc commands without a separate SSH client.

**Why this priority**: Terminal access is a powerful but advanced capability that operators expect in comprehensive management tools.

**Independent Test**: Can be fully tested by opening a terminal session to a connected server and executing a command with output displayed in the Manager.

**Acceptance Scenarios**:

1. **Given** a connected server, **When** the operator opens Terminal for that server, **Then** an interactive shell session is established and command output is displayed.
2. **Given** an active terminal session, **When** the server becomes unreachable, **Then** the terminal displays a connection-lost message and offers reconnection.

---

### User Story 10 - Operate Hub and Agent via Command Line (Priority: P2)

As a system administrator, I want to manage Hub and Agent services through command-line tools so that I can automate deployment, inspect status, and manage connections without the desktop Manager.

**Why this priority**: CLI tools are essential for headless server deployment and automation, and they define the operational contract between components.

**Independent Test**: Can be fully tested by running status, service control, log configuration, and connection management commands on both Hub and Agent installations.

**Acceptance Scenarios**:

1. **Given** a Hub is installed, **When** the administrator runs the status command, **Then** the Hub reports its running state, version, and active connections.
2. **Given** a Hub is running, **When** the administrator searches for and lists agents, **Then** discovered and connected agents are displayed with their names and connection status.
3. **Given** an Agent is installed, **When** the administrator sends a connection request to a Hub, **Then** the request is registered and connection status can be verified via connection test and last-sync commands.
4. **Given** a Hub or Agent service, **When** the administrator sets the log directory, **Then** subsequent log output is written to the specified directory.
5. **Given** a Hub is connected to a Manager, **When** the administrator runs manager connection test, **Then** connectivity and last sync time are reported.

---

### User Story 11 - Configure Logging and Retention from Manager (Priority: P2)

As an infrastructure operator, I want to configure log directories and per-file size limits for each server and hub from the Manager so that I can control log storage without logging into individual machines.

**Why this priority**: Centralized log configuration reduces operational overhead and ensures consistent retention policies across the fleet.

**Independent Test**: Can be fully tested by changing log directory and size limits for a connected server/hub from the Manager and verifying new log entries follow the updated configuration.

**Acceptance Scenarios**:

1. **Given** a connected server or hub, **When** the operator sets the log root directory from the Manager, **Then** subsequent log entries are written under `<DIR>/<YYYY>/<MM>/<DD>/` on that target.
2. **Given** a connected server or hub, **When** the operator sets size limits for `service.txt`, `errors.txt`, or `connectivity.txt`, **Then** each log file is rotated or truncated according to its configured limit.
3. **Given** a hub or agent with default configuration, **When** no custom directory is set, **Then** logs are written to the default root `/var/log/nvx`.

---

### User Story 12 - Run Custom and Built-In Server Actions (Priority: P2)

As an infrastructure operator, I want to execute built-in administrative actions (reboot, datetime, package update/upgrade) and save reusable custom commands so that I can perform routine and repeatable operations efficiently.

**Why this priority**: Administrative actions and custom commands are high-frequency operational tasks that complement service and app management.

**Independent Test**: Can be fully tested by running a built-in action (e.g., package update) and a saved custom command on a connected server, verifying results in the Manager.

**Acceptance Scenarios**:

1. **Given** a connected server, **When** the operator triggers reboot, datetime change, package update, or package upgrade, **Then** the action executes remotely and the result (success/failure) is displayed.
2. **Given** the operator defines a custom command with a name and script, **When** they save it, **Then** the command appears in the reusable command list for that server or globally.
3. **Given** a saved custom command, **When** the operator runs it against a selected server, **Then** command output is displayed in the Manager.

---

### User Story 13 - Define and Respond to Alert Triggers (Priority: P2)

As an infrastructure operator, I want to define alert triggers based on resource thresholds so that the system automatically sends email notifications or reboots a server when conditions are met.

**Why this priority**: Automated triggers enable proactive incident response without continuous manual monitoring.

**Independent Test**: Can be fully tested by creating a trigger (e.g., `IF CPU > 90 THEN Email ops@example.com`), simulating the threshold condition, and verifying the configured action executes.

**Acceptance Scenarios**:

1. **Given** the operator defines a trigger in the format `IF [CPU|GPU|RAM|STORAGE] [>|<] <threshold> THEN [Email|Reboot] <target>`, **When** the trigger is saved, **Then** it is evaluated continuously against incoming metrics.
2. **Given** a trigger with action Email, **When** the condition is met, **Then** an email is sent to the specified address within 1 minute.
3. **Given** a trigger with action Reboot, **When** the condition is met, **Then** the specified server reboots and the Manager reflects the state change.

---

### User Story 14 - Deploy Hub and Agent Remotely (Priority: P2)

As an infrastructure operator, I want to install Hub and Agent on remote Ubuntu servers from the Manager so that I can onboard new machines without manual CLI installation on each host.

**Why this priority**: Remote deployment accelerates fleet expansion and reduces onboarding friction.

**Independent Test**: Can be fully tested by selecting an Ubuntu server and initiating remote Hub or Agent installation from the Manager, then verifying the component appears in status lists.

**Acceptance Scenarios**:

1. **Given** a reachable Ubuntu server with valid credentials, **When** the operator initiates remote Agent installation from the Manager, **Then** the Agent is installed, started, and appears in Agent Status.
2. **Given** a reachable Ubuntu server with valid credentials, **When** the operator initiates remote Hub installation from the Manager, **Then** the Hub is installed, started, and appears in Hub List.
3. **Given** remote installation fails (network error, permission denied), **When** the operation completes, **Then** the Manager displays a descriptive error and no partial registration occurs.

---

### User Story 15 - Receive Server Events in Real Time (Priority: P2)

As an infrastructure operator, I want to receive server events (state changes, action results, connectivity changes) from agents in addition to periodic metrics so that I can respond to incidents as they happen.

**Why this priority**: Event-driven notifications complement polling-based metrics and enable faster incident awareness.

**Independent Test**: Can be fully tested by triggering a server state change (e.g., service restart) and verifying the event appears in the Manager event stream within seconds.

**Acceptance Scenarios**:

1. **Given** a connected agent, **When** a significant server event occurs (service state change, reboot, connectivity loss), **Then** the event is sent to the Manager and displayed in the server event feed.
2. **Given** the Manager is viewing a server's detail page, **When** new events arrive, **Then** they appear in real time without manual refresh.
3. **Given** a hub or agent connection is temporarily lost, **When** connectivity is restored, **Then** cached events are delivered to the Manager in chronological order.

---

### User Story 16 - Customize Server Visual Identity (Priority: P3)

As an infrastructure operator, I want to assign custom icons and badges to servers, hubs, and agents so that I can visually distinguish roles and status in lists and the graph network.

**Why this priority**: Visual customization improves fleet recognition in large deployments but is not required for core functionality.

**Independent Test**: Can be fully tested by uploading an icon to the assets folder, assigning it to a server, and verifying it appears in Server List and Graph Network views.

**Acceptance Scenarios**:

1. **Given** the operator adds an icon or badge to the assets folder, **When** they assign it to a server, hub, or agent, **Then** the custom image appears in all relevant Manager views.
2. **Given** a server has a custom badge indicating status (e.g., production, staging), **When** the operator views the Graph Network, **Then** the badge is displayed on the corresponding node.

---

### Edge Cases

- What happens when a Hub goes offline while Agents are connected through it? Agents remain running but become unreachable from the Manager until the Hub reconnects or an alternate Hub path exists; gRPC streams between Hub and Manager drop and cached outbound data is spooled locally.
- What happens when an Agent belongs to multiple Hubs but one Hub fails? The Agent remains reachable through other connected Hubs; the Manager continues to recognize the Agent under its single server identity.
- What happens when a server hosts both Hub and Agent and one component crashes? The surviving component continues operating independently; the Manager reflects partial availability for that server.
- What happens when a connection request is sent but never accepted? The request remains in a pending state with a configurable timeout after which it expires.
- What happens when two servers attempt to register with the same identity? The second registration is rejected and the administrator is notified of the conflict.
- What happens when metrics collection fails temporarily on an Agent? The Manager displays the last known values with a staleness indicator and retries automatically.
- What happens when an external IP added to the graph becomes unreachable? The node is displayed with an unreachable status and ping failure indicator.
- What happens when a Docker, service, or app management action fails mid-operation? The operator receives a descriptive error message and the system state reflects the actual outcome (no partial state reported as success).
- What happens when a Hub or Agent cannot reach its upstream peer for an extended period? Unsent data accumulates in the local temp/cache store; when the connection restores, cached data is flushed in chronological order. If cache storage reaches its limit, oldest entries are discarded with a warning logged to `errors.txt`.
- What happens when a trigger condition is met while the target server is offline? Email triggers still send; Reboot triggers are queued and executed when the server reconnects, or reported as failed if the server remains unreachable beyond a configurable timeout.
- What happens when remote Hub or Agent installation is interrupted? The Manager reports failure status; the operator must retry or complete installation manually via CLI on the target Ubuntu server.
- What happens when a log file reaches its configured size limit? The file is rotated according to the configured policy (truncate oldest entries or rotate to archive) without stopping the service.
- What happens when a server has no GPU? GPU metrics are reported as unavailable and GPU-based triggers are skipped for that server.
- What happens when two triggers conflict (e.g., one reboots while another sends email for the same condition)? All matching triggers fire independently; the operator is responsible for avoiding conflicting trigger definitions.

## Requirements *(mandatory)*

### Functional Requirements

#### Topology and Identity

- **FR-001**: The system MUST implement a three-tier routing model where the Manager reaches Agents exclusively through Hubs (Manager → Hub → Agent).
- **FR-002**: Each physical or virtual server MUST host at most one Agent instance and at most one Hub instance.
- **FR-003**: A single server MAY host both an Agent and a Hub simultaneously, with each component managed independently.
- **FR-004**: An Agent MUST be able to connect to one or more Hubs simultaneously.
- **FR-005**: The Manager MUST recognize each server as a single entity regardless of how many Hubs an Agent on that server is connected to.
- **FR-006**: The system MUST support connection establishment via explicit request-and-accept between Agent↔Hub and Hub↔Manager pairs.
- **FR-006a**: Hub and Agent MUST support Ubuntu operating systems only in the initial release.
- **FR-006b**: All entity names, configuration keys, and CLI parameters MUST use standard, content-related naming (e.g., `server_alias`, `hub_name`, `agent_name`, `connectivity_status`).
- **FR-006c**: All inter-component communication between Manager↔Hub and Hub↔Agent MUST use gRPC, including metrics streaming, event delivery, remote commands, and terminal session relay.

#### Manager Application — Monitoring

- **FR-007**: The Manager MUST display server health metrics including CPU, GPU, RAM, storage, network utilization, and uptime.
- **FR-008**: The Manager MUST provide an "all servers" monitoring mode showing fleet-wide aggregated metrics.
- **FR-009**: The Manager MUST provide a "per-server" monitoring mode showing individual server metrics with charts displaying trends over time.

#### Manager Application — Server List

- **FR-010**: The Manager MUST display a Server List showing for each server: alias name, CPU, RAM, motherboard, storage, network interfaces, virtualization status, OS information, and uptime.

#### Manager Application — Docker Manager

- **FR-011**: The Manager MUST display a list of Docker containers (running and stopped) for a selected server.
- **FR-012**: The Manager MUST display a list of Docker images for a selected server.

#### Manager Application — Network Manager

- **FR-013**: The Manager MUST display network interfaces for a selected server.
- **FR-014**: The Manager MUST provide remote execution of ping, traceroute, and DNS lookup tools from a selected server.
- **FR-015**: The Manager MUST provide additional common network diagnostic utilities beyond ping, traceroute, and DNS lookup.

#### Manager Application — Service Manager

- **FR-016**: The Manager MUST display a list of system services on a selected server with their current status.
- **FR-017**: The Manager MUST allow operators to enable, disable, and restart services on a selected server.

#### Manager Application — App Manager

- **FR-018**: The Manager MUST display a list of installed applications on a selected server.
- **FR-019**: The Manager MUST allow operators to install and uninstall applications on a selected server.

#### Manager Application — Terminal

- **FR-020**: The Manager MUST provide an interactive remote terminal session to a selected connected server.

#### Manager Application — Hub and Agent Status

- **FR-021**: The Manager MUST display a Hub List showing each Hub's name, associated server, and description.
- **FR-022**: The Manager MUST display an Agent Status view showing each Agent's name, associated server, and description.

#### Manager Application — Graph Network

- **FR-023**: The Manager MUST display an interactive graph with servers represented as nodes and connections showing real-time ping latency.
- **FR-024**: The Manager MUST allow operators to add external IP addresses as nodes in the graph network.
- **FR-025**: The Manager MUST allow operators to reposition nodes in the graph, persisting layout preferences across sessions.
- **FR-026**: The Manager MUST update ping latency values on the graph in real time.

#### Manager Application — General

- **FR-027**: The Manager MUST provide a Settings area for configuring application preferences.
- **FR-028**: The Manager MUST provide an About section displaying application version and relevant information.
- **FR-028a**: The Manager MUST use the Shadcn UI component library for all screens, forms, tables, charts, and navigation elements.
- **FR-028b**: The Manager MUST provide a dedicated assets folder for storing custom icons and badges assignable to servers, hubs, and agents.

#### Manager Application — Log Management

- **FR-041**: The Manager MUST allow operators to configure the log root directory per connected server or hub.
- **FR-042**: The Manager MUST allow operators to set independent size limits for `service.txt`, `errors.txt`, and `connectivity.txt` per server or hub.
- **FR-043**: The Manager MUST display current log configuration (directory, limits, current file sizes) for each server and hub.

#### Manager Application — Custom Commands and Server Actions

- **FR-044**: The Manager MUST allow operators to define, name, store, and execute custom commands against selected servers.
- **FR-045**: The Manager MUST provide built-in server administration actions including reboot, datetime configuration, package update, and package upgrade.
- **FR-046**: The Manager MUST display action results (success, failure, output) for both built-in actions and custom commands.

#### Manager Application — Alert Triggers

- **FR-047**: The Manager MUST allow operators to define alert triggers using the format: `IF [CPU|GPU|RAM|STORAGE] [>|<] <threshold> THEN [Email|Reboot] <target>`.
- **FR-048**: The system MUST evaluate active triggers continuously against incoming metrics from connected agents.
- **FR-049**: When a trigger with action Email fires, the system MUST send a notification to the specified email address.
- **FR-050**: When a trigger with action Reboot fires, the system MUST initiate a reboot on the specified target server.

#### Manager Application — Remote Deployment

- **FR-051**: The Manager MUST allow operators to remotely install the Agent on reachable Ubuntu servers.
- **FR-052**: The Manager MUST allow operators to remotely install the Hub on reachable Ubuntu servers.
- **FR-053**: Remote installation MUST report progress and final status (success or failure with reason) in the Manager.

#### Manager Application — Events

- **FR-054**: The Manager MUST display a real-time event feed per server showing state changes, action results, and connectivity events reported by agents.
- **FR-055**: Events MUST be delivered to the Manager via the Hub relay path alongside periodic metrics.

#### Hub — Logging

- **FR-056**: The Hub MUST write logs to `<DIR>/<YYYY>/<MM>/<DD>/service.txt`, `errors.txt`, and `connectivity.txt` where `<DIR>` defaults to `/var/log/nvx`.
- **FR-057**: The Hub MUST append a connectivity log entry to `connectivity.txt` every 1 minute in the format: `<ISO8601-timestamp> | direction=<inbound|outbound> | peer=<name> | status=<connected|disconnected|failed> | latency_ms=<integer> | detail=<message>`.
- **FR-058**: The Hub MUST retain connectivity log entries for 5 days; entries older than 5 days MUST be purged automatically.

#### Hub — Data Persistence and Resilience

- **FR-059**: Each Hub MUST maintain a local SQLite database for operational data including connection records, routing state, and sync metadata.
- **FR-060**: The Hub MUST persist unsent outbound data in a local temp/cache store when upstream connections fail and flush cached data in chronological order upon reconnection.

#### Agent — Logging

- **FR-061**: The Agent MUST write logs to `<DIR>/<YYYY>/<MM>/<DD>/service.txt`, `errors.txt`, and `connectivity.txt` where `<DIR>` defaults to `/var/log/nvx`.
- **FR-062**: The Agent MUST append a connectivity log entry to `connectivity.txt` every 1 minute using the same standard format as the Hub.
- **FR-063**: The Agent MUST retain connectivity log entries for 5 days; entries older than 5 days MUST be purged automatically.

#### Agent — Data Persistence and Resilience

- **FR-064**: The Agent MUST persist unsent outbound data (metrics, events) in a local temp/cache store when Hub connections fail and flush cached data in chronological order upon reconnection.
- **FR-065**: The Agent MUST collect and report server events including state changes, action results, and connectivity changes to connected Hub(s) in addition to periodic metrics.

#### Hub Command-Line Interface

- **FR-029**: The Hub MUST provide a status command reporting running state and connection summary.
- **FR-030**: The Hub MUST provide service control commands to disable, enable, and restart the Hub service.
- **FR-031**: The Hub MUST provide a command to configure the log output directory.
- **FR-032**: The Hub MUST provide commands to search for, list, request connection to, disconnect from, inspect connection status of, and check last sync time with Agents.
- **FR-033**: The Hub MUST provide commands to search for, list, request connection to, disconnect from, test connection with, and check last sync time with Managers.

#### Agent Command-Line Interface

- **FR-034**: The Agent MUST provide a status command reporting running state and connection summary.
- **FR-035**: The Agent MUST provide service control commands to disable, enable, and restart the Agent service.
- **FR-036**: The Agent MUST provide a command to configure the log output directory.
- **FR-037**: The Agent MUST provide commands to search for, list, request connection to, disconnect from, test connection with, and check last sync time with Hubs.

#### Data Collection (Agent Responsibilities)

- **FR-038**: The Agent MUST collect and report hardware specifications (CPU, GPU, RAM, motherboard, storage, NICs, VM status, OS info) to its connected Hub(s).
- **FR-039**: The Agent MUST collect and report real-time and historical metrics (CPU, GPU, RAM, storage, network, uptime) to its connected Hub(s).
- **FR-040**: The Agent MUST execute remote management commands relayed from the Manager via Hub (Docker queries, service control, app management, network diagnostics, terminal sessions).

### Key Entities

- **Server**: A physical or virtual machine identified uniquely by the Manager; hosts at most one Agent and at most one Hub; attributes include alias name, hardware specs, OS info, uptime, and online status.
- **Manager**: The central desktop application through which operators monitor and manage the fleet; maintains connections to Hubs and presents unified views.
- **Hub**: An intermediary service on a server that relays communication between the Manager and Agents; maintains its own connections to Managers and Agents.
- **Agent**: A lightweight service on a server that collects metrics, reports hardware info, and executes management commands; may connect to multiple Hubs.
- **Connection**: A trust relationship between two components (Agent↔Hub or Hub↔Manager) established via request-and-accept; has states including pending, connected, and disconnected.
- **Metric Snapshot**: A point-in-time or aggregated record of CPU, GPU, RAM, storage, network, and uptime data for a server.
- **Hub Record**: A Manager-side representation of a Hub with name, associated server, and description.
- **Agent Record**: A Manager-side representation of an Agent with name, associated server, and description.
- **Graph Node**: A visual representation of a server or external IP in the topology graph, with position coordinates, optional custom icon/badge, and real-time ping latency to connected nodes.
- **Docker Container / Image**: A containerized workload or its image as reported from a server's Docker environment.
- **System Service**: An OS-level service with name, status (enabled/disabled, running/stopped), manageable remotely.
- **Application**: An installed software package on a server, supporting list, install, and uninstall operations.
- **Log Configuration**: Per server/hub settings defining log root directory (default `/var/log/nvx`), per-file size limits, and current usage for `service.txt`, `errors.txt`, and `connectivity.txt`.
- **Custom Command**: A user-defined, named command script stored in the Manager and executable against one or more servers.
- **Alert Trigger**: A rule in the format `IF [CPU|GPU|RAM|STORAGE] [>|<] <threshold> THEN [Email|Reboot] <target>` evaluated continuously against incoming metrics.
- **Server Event**: A discrete notification from an Agent describing a state change, action result, or connectivity event, delivered to the Manager in real time.
- **Visual Asset**: A custom icon or badge stored in the Manager assets folder and assignable to servers, hubs, or agents.
- **Outbound Cache**: A local temp/cache store on Hub or Agent that holds unsent metrics, events, or commands during connection failures.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can view health metrics for all connected servers within 5 seconds of opening the Monitoring view.
- **SC-002**: Per-server monitoring charts display metric history covering at least the past 24 hours.
- **SC-003**: Server inventory details (hardware specs, OS, uptime) are populated automatically within 30 seconds of a new Agent connection being established.
- **SC-004**: Hub-to-Agent and Hub-to-Manager connection establishment completes within 10 seconds under normal network conditions.
- **SC-005**: Graph Network ping latency values refresh at least every 5 seconds for all visible nodes.
- **SC-006**: Remote service enable/disable/restart actions complete and reflect updated status within 15 seconds.
- **SC-007**: Network diagnostic tools (ping, traceroute, DNS lookup) return initial results within 3 seconds of execution.
- **SC-008**: Operators can manage a fleet of at least 100 servers from a single Manager instance without degraded responsiveness in list and monitoring views.
- **SC-009**: 95% of remote management actions (service control, Docker listing, app install/uninstall) complete successfully on first attempt when the target server is online.
- **SC-010**: CLI status and connection commands produce output within 2 seconds when run locally on the Hub or Agent host.
- **SC-011**: Connectivity log entries are written at 1-minute intervals with no more than 5 seconds drift on Hub and Agent hosts.
- **SC-012**: Cached outbound data is fully flushed within 30 seconds of connection restoration.
- **SC-013**: Alert triggers evaluate conditions and fire configured actions within 60 seconds of threshold breach.
- **SC-014**: Server events appear in the Manager event feed within 5 seconds of occurrence on the agent host.
- **SC-015**: Remote Hub or Agent installation on Ubuntu completes within 5 minutes under normal network conditions.
- **SC-016**: Custom commands saved by an operator are available for execution immediately without application restart.

## Assumptions

- Operators have network connectivity between Manager workstations and Hub servers, and between Hub and Agent servers.
- Hub and Agent components run on Ubuntu servers only in the initial release; other Linux distributions and operating systems are out of scope.
- Agents and Hubs can be installed via CLI on target servers or remotely from the Manager; remote installation requires SSH access with sufficient privileges on the target Ubuntu server.
- Docker management applies only to servers where Docker is installed and accessible to the Agent.
- Application management and package update/upgrade operate through Ubuntu's native package management system (`apt`); custom or proprietary installers are out of scope for the initial release.
- Connection security uses a mutual trust model established through the request-and-accept workflow; all component connections use gRPC as the transport protocol (FR-006c).
- The Manager is a desktop application installed on the operator's workstation using the Shadcn UI component library; web-based access is out of scope for the initial release.
- Metric retention defaults to 30 days of historical data per server; longer retention is configurable via Settings.
- Terminal sessions support a single concurrent session per server per operator; multiple operators may each hold their own session.
- External IP nodes in the graph network are ping targets only and do not participate in the Hub-Agent management topology.
- Log directory changes via Manager or CLI take effect immediately for new log entries; existing logs are not migrated automatically.
- Default log root directory is `/var/log/nvx` on all Hub and Agent hosts; operators may override per server or hub.
- Connectivity logs are retained for 5 days; service and error logs follow configured size limits set per server or hub.
- Hub operational data is persisted in a local SQLite database; database schema details will be defined during planning.
- Outbound cache storage on Hub and Agent has a default capacity sufficient for 24 hours of disconnected operation; overflow behavior discards oldest entries.
- Email trigger delivery requires SMTP configuration defined in Manager Settings; SMTP setup details will be defined during planning.
- GPU metrics are collected when a GPU is present; servers without GPU report the metric as unavailable.
- Custom icons and badges are stored in a Manager-local assets folder and are not synchronized across Manager installations.
