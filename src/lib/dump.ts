// 🔥 we use this funky conversion just for the purpose of producing a readable dump -- the real conversion is done in the Go code -- we skip the first 64 entries and start on line 64 so it's easy to read the line number as the EDCDIC character and the value as the ASCII equivalent

const ebcdic: any[] = [
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  // start on line 64 to make reconciliation easier
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '¢',
  '.',
  '<',
  '(',
  '+',
  '|',
  '&',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '!',
  '$',
  '*',
  ')',
  ';',
  '¬',
  '-',
  '/',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '|',
  ',',
  '%',
  '_',
  '>',
  '?',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '`',
  ':',
  '#',
  '@',
  "'",
  '=',
  '\"',
  '\u00a0',
  'a',
  'b',
  'c',
  'd',
  'e',
  'f',
  'g',
  'h',
  'i',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  'j',
  'k',
  'l',
  'm',
  'n',
  'o',
  'p',
  'q',
  'r',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  's',
  't',
  'u',
  'v',
  'w',
  'x',
  'y',
  'z',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '`',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '{',
  'A',
  'B',
  'C',
  'D',
  'E',
  'F',
  'G',
  'H',
  'I',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '}',
  'J',
  'K',
  'L',
  'M',
  'N',
  'O',
  'P',
  'Q',
  'R',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\\',
  '\u00a0',
  'S',
  'T',
  'U',
  'V',
  'W',
  'X',
  'Y',
  'Z',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '0',
  '1',
  '2',
  '3',
  '4',
  '5',
  '6',
  '7',
  '8',
  '9',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0',
  '\u00a0'
];

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

export function e2a(e: number[]): string {
  let a = '';
  for (let i = 0; i < e.length; i++) {
    // @ts-ignore 🔥 we know e[i] is always valid
    if (e[i] >= 64) a += ebcdic[e[i] - 64];
    else a += '\u2022';
  }
  return a;
}

function toHex(num: number, pad: number): string {
  const padding = '0000000000000000'.substring(0, pad);
  const hex = num.toString(16);
  return padding.substring(0, padding.length - hex.length) + hex;
}
