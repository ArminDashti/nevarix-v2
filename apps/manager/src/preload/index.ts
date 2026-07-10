import { contextBridge } from "electron";

contextBridge.exposeInMainWorld("nvx", {
  version: "0.1.0",
});
