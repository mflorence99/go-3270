export type Go3270 = {
  close: () => void;
  keystroke: (
    code: string,
    key: string,
    alt: boolean,
    ctrl: boolean,
    shift: boolean
  ) => void;
  receiveFromApp: (bytes: Uint8ClampedArray) => void;
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
