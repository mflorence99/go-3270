import { dumpBytes } from '$lib/dump';

export type EmulatorContext = {
  color: string;
  ctx: CanvasRenderingContext2D;
  dims: [number, number];
  fontHeight: number;
  fontSpec: string;
  fontWidth: number;
  paddingLeft: number /* 👈 as a fraction of fontWidth */;
  paddingTop: number /* 👈 as a fraction of fontHeight */;
  responder: (bytes: Uint8Array) => void;
};

// 🟧 3270 data stream protocol

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

export class Lu3270 {
  private constructor(public ectx: EmulatorContext) {}

  // 🔥 keep factory pattern congruent with Tn3270,
  //    even if not strictly needed
  static lu3270(ectx: EmulatorContext): Lu3270 {
    return new Lu3270(ectx);
  }

  // 🔥 similarly, code close for symmetry and in case
  //    we have resources to free
  close(): void {}

  // 🔥 think differently! this code emulates the device and
  //    inbound data streams are sent FROM the device
  //    TO application code
  inbound(bytes: Uint8Array): void {
    const { responder } = this.ectx;
    dumpBytes(bytes, 'Inbound 3270 -> Application', true, 'palegreen');
    responder(bytes);
  }

  // 🔥 think differently! this code emulates the device and
  //    outbound data streams flow FROM application code
  //    TO the device
  outbound(bytes: Uint8Array): void {
    dumpBytes(bytes, 'Outbound Application -> 3270', true, 'yellow');
    // 🔥 TEMPORARY
    this.inbound(new Uint8Array([193, 194, 195] /* 👈 EBCDIC "ABC" */));
    this.refresh();
  }

  refresh(): void {
    // 🔥 TEMPORARY
    const {
      color,
      ctx,
      dims,
      fontHeight,
      fontSpec,
      fontWidth,
      paddingLeft,
      paddingTop
    } = this.ectx;
    // 👇 establish terminal font and color
    ctx.font = fontSpec;
    ctx.textAlign = 'left';
    ctx.textBaseline = 'top';
    ctx.fillStyle = color;
    ctx.lineWidth = 2;
    ctx.strokeStyle = color;
    // 👇 fill every cell with a random character
    const chars =
      'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/? ';
    for (let ix = 0, x = 0; ix < dims[0]; ix++, x += fontWidth) {
      for (let iy = 0, y = 0; iy < dims[1]; iy++, y += fontHeight) {
        if (ix === 0 && iy === 0) {
          ctx.strokeRect(x, y, fontWidth, fontHeight);
        }
        ctx.fillText(
          chars.charAt(Math.floor(Math.random() * chars.length)),
          x + fontWidth * paddingLeft,
          y + fontHeight * paddingTop
        );
      }
    }
  }
}
