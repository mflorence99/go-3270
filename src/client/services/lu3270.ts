import { dumpBytes } from '$lib/dump';

// 🟧 3270 data stream protocol

// 👁️ https://bitsavers.org/pdf/ibm/3270/GA23-0059-07_3270_Data_Stream_Programmers_Reference_199206.pdf
// 👁️ http://www.prycroft6.com.au/misc/3270.html
// 👁️ http://www.tommysprinkle.com/mvs/P3270/start.htm

export class Lu3270 {
  constructor(
    private ctx: CanvasRenderingContext2D,
    private color: string,
    private fontSize: number,
    private width: number,
    private height: number,
    private fontWidth: number,
    private fontHeight: number,
    private paddingLeft: number /* 👈 as a fraction of fontWidth */,
    private paddingTop: number /* 👈 as a fraction of fontHeight */,
    private responder: (bytes: Uint8Array) => void
  ) {}

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
    window.testGo('xxx');
    this.inbound(window.renderGo('xxx'));
    this.refresh();
  }

  refresh(): void {
    // 🔥 TEMPORARY what we really need to do is to refresh the display with "current" data, but with new font size, color etc
    // 👇 establish terminal font and color
    this.ctx.font = `${this.fontSize}px Terminal`;
    this.ctx.textAlign = 'left';
    this.ctx.textBaseline = 'top';
    this.ctx.fillStyle = this.color;
    this.ctx.lineWidth = 2;
    this.ctx.strokeStyle = this.color;
    // 👇 fill every cell with a random character
    const chars =
      'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[{]};:,<.>/? ';
    for (
      let ix = 0, x = 0;
      ix < this.width;
      ix++, x += this.fontWidth
    ) {
      for (
        let iy = 0, y = 0;
        iy < this.height;
        iy++, y += this.fontHeight
      ) {
        if (ix === 0 && iy === 0) {
          this.ctx.strokeRect(x, y, this.fontWidth, this.fontHeight);
        }
        this.ctx.fillText(
          chars.charAt(Math.floor(Math.random() * chars.length)),
          x + this.fontWidth * this.paddingLeft,
          y + this.fontHeight * this.paddingTop
        );
      }
    }
  }
}
