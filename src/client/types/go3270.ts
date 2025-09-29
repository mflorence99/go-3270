export type Go3270 = {
  close: () => Uint8ClampedArray;
  receive: (bytes: Uint8ClampedArray) => Uint8ClampedArray;
  restore: (bytes: Uint8ClampedArray) => void;
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
