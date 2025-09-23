import { Go3270 } from '$client/types/go3270';

import { dumpBytes } from '$lib/dump';

// 🟧 3270 data stream protocol

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

export class Lu3270 {
  go3270: Go3270;

  constructor(
    canvas: HTMLCanvasElement,
    color: string,
    fontSize: number,
    cols: number,
    rows: number,
    dpi: number,
    private responder: (bytes: Uint8Array) => void
  ) {
    // 🔥 ctor will be called before WASM is initialized
    this.go3270 = window.NewGo3270?.(
      canvas,
      color,
      fontSize,
      cols,
      rows,
      dpi
    );
  }

  // 🔥 we nay have resources to free etc
  close(): void {}

  // 🔥 this class emulates the device and inbound data streams are sent FROM the device TO application code
  inbound(bytes: Uint8Array): void {
    dumpBytes(bytes, 'Inbound 3270 -> Application', true, 'palegreen');
    this.responder?.(bytes);
  }

  // 🔥 this class emulates the device and outbound data streams flow FROM application code TO the device
  outbound(bytes: Uint8Array): void {
    dumpBytes(bytes, 'Outbound Application -> 3270', true, 'yellow');
    // 🔥 TEMPORARY
    this.inbound(this.go3270.inbound());
  }

  // 🔥 TEMPORARY what we really need to do is to refresh the display with "current" data, but with new font size, color etc
  refresh(): void {
    this.go3270?.testPattern();
  }
}
