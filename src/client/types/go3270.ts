import { CLUT } from '$client/state/state';

export type Go3270 = {
  close: () => void;
  focus: (focus: boolean) => void;
  keystroke: (
    code: string,
    key: string,
    alt: boolean,
    ctrl: boolean,
    shift: boolean
  ) => void;
  outbound: (bytes: Uint8ClampedArray) => void;
};

declare global {
  export interface Window {
    Go: any;
    NewGo3270: (
      canvas: HTMLCanvasElement,
      bgColor: string,
      monochrome: boolean,
      clut: CLUT,
      fontSize: number,
      cols: number,
      rows: number,
      dpi: number
    ) => Go3270;
  }
}
