declare global {
  export interface Window {
    Go: any;
    renderGo: (...args: any) => Uint8Array;
    testGo: (name: string) => string;
  }
}
export {};
