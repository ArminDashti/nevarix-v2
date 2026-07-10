export interface NvxApi {
  version: string;
}

declare global {
  interface Window {
    nvx: NvxApi;
  }
}

export {};
