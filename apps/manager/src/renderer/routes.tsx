import { MonitoringDashboard } from "@/pages/MonitoringDashboard";
import { ServerList } from "@/pages/ServerList";
import { HubList } from "@/pages/HubList";
import { AgentStatus } from "@/pages/AgentStatus";
import { Settings } from "@/pages/Settings";
import { About } from "@/pages/About";

export const routes = [
  { path: "/", element: <MonitoringDashboard />, label: "Monitoring" },
  { path: "/servers", element: <ServerList />, label: "Servers" },
  { path: "/hubs", element: <HubList />, label: "Hubs" },
  { path: "/agents", element: <AgentStatus />, label: "Agents" },
  { path: "/settings", element: <Settings />, label: "Settings" },
  { path: "/about", element: <About />, label: "About" },
];
