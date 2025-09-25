export type Go3270 = {
  close: () => Uint8Array;
  datastream: (bytes: Uint8Array) => Uint8Array;
  restore: (bytes: Uint8Array) => void;
  testPattern: () => void;
};

declare global {
  export interface Window {
    Go: any;
    NewGo3270: (
      canvas: HTMLCanvasElement,
      color: string,
      fontSize: number,
      cols: number,
      rows: number,
      dpi: number
    ) => Go3270;
  }
}
