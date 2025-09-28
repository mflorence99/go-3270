// 🔥 Bun bundler doesn't like $lib here ???
import { e2a } from './convert';

// 🟧 Log a Uint8ClampedArray like an old-fashioned dump

export function dumpBytes(
  data: Uint8ClampedArray,
  title: string,
  ebcdic = false,
  color = 'blue'
): void {
  const sliceSize = 32;
  let offset = 0;
  const total = data.length;
  console.groupCollapsed(
    `%c${title} ${ebcdic ? '(EBCDIC-encoded)' : ''}`,
    `color: ${color}`
  );
  console.log(
    '%c       00       04       08       0c       10       14       18       1c        00  04  08  0c  10  14  18  1c  ',
    'color: skyblue; font-weight: bold'
  );
  while (true) {
    const slice = new Uint8ClampedArray(
      data.slice(offset, Math.min(offset + sliceSize, total))
    );
    const { hex, str } = dumpSlice(slice, sliceSize, ebcdic);
    console.log(
      `%c${toHex(offset, 6)} %c${hex} %c${str}`,
      'color: skyblue; font-weight: bold',
      'color: white',
      'color: wheat'
    );
    // setup for next time
    if (slice.length < sliceSize) break;
    offset += sliceSize;
  }
  console.groupEnd();
}

// 🟦 Helpers

function dumpSlice(
  bytes: Uint8ClampedArray,
  sliceSize: number,
  ebcdic: boolean
): { hex: string; str: string } {
  let hex = '';
  let str = '';
  let ix = 0;
  // 👇 decode to hex and string equiv
  for (; ix < bytes.length; ix++) {
    const byte = bytes[ix];
    if (byte == null) break;
    hex += toHex(byte, 2);
    const char = ebcdic ? e2a([byte]) : String.fromCharCode(byte);
    // 👇 use special character in string as a visual aid to counting
    str += char === '\u00a0' || char === ' ' ? '\u2022' : char;
    if (ix > 0 && ix % 4 === 3) hex += ' ';
  }
  // 👇 pad remainder of slice
  for (; ix < sliceSize; ix++) {
    hex += '  ';
    str += ' ';
    if (ix > 0 && ix % 4 === 3) hex += ' ';
  }
  return { hex, str };
}

function toHex(num: number, pad: number): string {
  const padding = '0000000000000000'.substring(0, pad);
  const hex = num.toString(16);
  return padding.substring(0, padding.length - hex.length) + hex;
}
