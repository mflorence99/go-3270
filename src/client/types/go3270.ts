export type Go3270 = {
  close: () => void;
  focussed: (focussed: boolean) => void;
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
      bgColor: string,
      color: string,
      fontSize: number,
      cols: number,
      rows: number,
      dpi: number
    ) => Go3270;
  }
}
